package loophole

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/apiclient"
	"github.com/loophole/cli/internal/pkg/closehandler"
	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/loophole/cli/internal/pkg/httpserver"
	"github.com/loophole/cli/internal/pkg/keys"
	"github.com/loophole/cli/internal/pkg/urlmaker"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
)

// remote forwarding port (on remote SSH server network)
var remoteEndpoint = lm.Endpoint{
	Host: "127.0.0.1",
	Port: 80,
}

func handleClient(client net.Conn, local net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	// Start local -> client data transfer
	go func() {
		_, err := io.Copy(client, local)
		if err != nil {
			if el := log.Debug(); el.Enabled() {
				el.Err(err).Msg("Error copying local -> client:")
			}
		}
		chDone <- true
	}()

	// Start client -> local data transfer
	go func() {
		_, err := io.Copy(local, client)
		if err != nil {
			if el := log.Debug(); el.Enabled() {
				el.Err(err).Msg("Error copying client -> local:")
			}
		}
		chDone <- true
	}()

	<-chDone
}

func registerDomain(apiURL string, publicKey *ssh.PublicKey, requestedSiteID string) string {
	communication.StartLoading("Registering your domain...")
	siteID, err := apiclient.RegisterSite(apiURL, *publicKey, requestedSiteID)
	if err != nil {
		communication.LoadingFailure()
		if requestErr, ok := err.(apiclient.RequestError); ok {
			log.Error().Int("status", requestErr.StatusCode).Msg("Request ended")
			log.Error().Msg(requestErr.Message)
			log.Error().Msg(fmt.Sprintf("Details: %s", requestErr.Details))
			log.Fatal().Msg("Please fix the above issue and try again")
		} else {
			log.Fatal().Err(err).Msg("Something unexpected happened, please let developers know")
		}
	}
	communication.LoadingSuccess()
	return siteID
}

func connectViaSSH(gatewayEndpoint lm.Endpoint, username string, authMethod ssh.AuthMethod) *ssh.Client {
	var serverSSHConnHTTPS *ssh.Client
	sshConfigHTTPS := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	var sshSuccess bool = false
	var sshRetries int = 5
	var err error
	for i := 0; i < sshRetries && !sshSuccess; i++ { // Connection retries in case of reconnect during gateway shutdown
		communication.StartLoading("Initializing secure tunnel... ")
		serverSSHConnHTTPS, err = ssh.Dial("tcp", gatewayEndpoint.URI(), sshConfigHTTPS)
		if err != nil {
			communication.LoadingFailure()
			log.Info().Msg(fmt.Sprintf("SSH Connection failed, retrying in 10 seconds... (Attempt %d/%d)", i+1, sshRetries))
			time.Sleep(10 * time.Second)
		} else {
			sshSuccess = true
		}
	}
	if !sshSuccess {
		communication.WriteRed("An error occured while dialing into SSH. If your connection has been running for a while")
		communication.WriteRed("this might be caused by the server shutting down your connection.")
		log.Fatal().Err(err).Msg("Dialing SSH Gateway for HTTPS failed.")
	}
	if el := log.Debug(); el.Enabled() {
		fmt.Println()
		el.Msg("Dialing SSH Gateway for HTTPS succeeded")
	}
	communication.LoadingSuccess()
	return serverSSHConnHTTPS
}

func createTLSReverseProxy(localEndpoint lm.Endpoint, siteID string, basicAuthUsername string, basicAuthPassword string, displayOptions lm.DisplayOptions) *http.Server {
	communication.StartLoading("Starting local TLS proxy server")
	serverBuilder := httpserver.New().
		WithHostname(siteID).
		Proxy().
		ToEndpoint(localEndpoint)

	if basicAuthUsername != "" && basicAuthPassword != "" {
		serverBuilder = serverBuilder.
			WithBasicAuth(basicAuthUsername, basicAuthPassword)
	}
	if displayOptions.DisableProxyErrorPage {
		serverBuilder = serverBuilder.
			DisableProxyErrorPage()
	}

	if el := log.Debug(); el.Enabled() {
		el.
			Str("target", localEndpoint.URI()).
			Msg("Proxy via http created")
	}
	server, err := serverBuilder.Build()
	if err != nil {
		communication.LoadingFailure()
		log.Fatal().Err(err).Msg("Something went wrong while creating server")
	}
	return server
}

func startLocalHTTPServer(server *http.Server) *lm.Endpoint {
	communication.StartLoading("Starting local proxy server... ")
	if el := log.Debug(); el.Enabled() {
		fmt.Println()
		el.Msg("Server for proxy created")
	}
	localListener, err := net.Listen("tcp", ":0")
	if err != nil {
		communication.LoadingFailure()
		log.Fatal().Err(err).Msg("Failed to listen on TLS proxy for HTTPS")
	}
	localListenerEndpoint := &lm.Endpoint{
		Host: "127.0.0.1",
		Port: int32(localListener.Addr().(*net.TCPAddr).Port),
	}
	if el := log.Debug(); el.Enabled() {
		el.
			Int32("port", localListenerEndpoint.Port).
			Msg("Proxy listener for HTTPS started")
	}
	go func() {
		err := server.ServeTLS(localListener, "", "")
		if err != nil {
			communication.LoadingFailure()
			log.Fatal().Msg("Failed to start TLS server")
		}
	}()
	if el := log.Debug(); el.Enabled() {
		el.Msg("Started server TLS server")
	}
	communication.LoadingSuccess()
	return localListenerEndpoint
}

func startRemoteForwardServer(serverSSHConnHTTPS *ssh.Client) net.Listener {
	listenerHTTPSOverSSH, err := serverSSHConnHTTPS.Listen("tcp", remoteEndpoint.URI())
	if err != nil {
		communication.LoadingFailure()
		log.Fatal().Err(err).Msg("Listening on remote endpoint for HTTPS failed")
	}
	if el := log.Debug(); el.Enabled() {
		fmt.Println()
		el.Msg("Listening on remote endpoint for HTTPS succeeded")
	}
	return listenerHTTPSOverSSH
}

func parsePublicKey(identityFile string) (ssh.AuthMethod, ssh.PublicKey) {
	publicKeyAuthMethod, publicKey, err := keys.ParsePublicKey(identityFile)
	if err != nil {
		communication.LoadingFailure()
		log.Fatal().Err(err).Msg("No public key available")
	}

	return publicKeyAuthMethod, publicKey
}

func getStaticFileServer(path string, siteID string, basicAuthUsername string, basicAuthPassword string) *http.Server {
	communication.StartLoading("Starting local file server")
	serverBuilder := httpserver.New().
		WithHostname(siteID).
		ServeStatic().
		FromDirectory(path)

	if basicAuthUsername != "" && basicAuthPassword != "" {
		serverBuilder = serverBuilder.
			WithBasicAuth(basicAuthUsername, basicAuthPassword)
	}

	communication.LoadingSuccess()
	server, err := serverBuilder.Build()
	if err != nil {
		communication.LoadingFailure()
		log.Fatal().Err(err).Msg("Something went wrong while creating server")
	}
	return server
}

func getWebdavServer(path string, siteID string, basicAuthUsername string, basicAuthPassword string) *http.Server {
	communication.StartLoading("Starting WebDav server")
	serverBuilder := httpserver.New().
		WithHostname(siteID).
		ServeWebdav().
		FromDirectory(path)

	if basicAuthUsername != "" && basicAuthPassword != "" {
		serverBuilder = serverBuilder.
			WithBasicAuth(basicAuthUsername, basicAuthPassword)
	}

	communication.LoadingSuccess()
	server, err := serverBuilder.Build()
	if err != nil {
		communication.LoadingFailure()
		log.Fatal().Err(err).Msg("Something went wrong while creating server")
	}
	return server
}

func listenOnRemoteEndpoint(serverSSHConnHTTPS *ssh.Client) net.Listener {
	listenerHTTPSOverSSH, err := serverSSHConnHTTPS.Listen("tcp", remoteEndpoint.URI())
	if err != nil {
		communication.LoadingFailure()
		log.Fatal().Err(err).Msg("Listening on remote endpoint for HTTPS failed")
	}
	return listenerHTTPSOverSSH
}

// ForwardPort is used to forward external URL to locally available port
func ForwardPort(config lm.ExposeHttpConfig) {
	communication.PrintWelcomeMessage()

	protocol := "http"
	if config.Local.HTTPS {
		protocol = "https"
	}
	localEndpoint := lm.Endpoint{
		Protocol: protocol,
		Host:     config.Local.Host,
		Port:     config.Local.Port,
	}

	publicKeyAuthMethod, publicKey := parsePublicKey(config.Remote.IdentityFile)
	siteID := registerDomain(config.Remote.APIEndpoint.URI(), &publicKey, config.Remote.SiteID)
	server := createTLSReverseProxy(localEndpoint, siteID, config.Remote.BasicAuthUsername, config.Remote.BasicAuthPassword, config.Display)
	forward(config.Remote, config.Display, publicKeyAuthMethod, siteID, server, localEndpoint.URI(), []string{"https"})
}

// ForwardDirectory is used to expose local directory via HTTP (download only)
func ForwardDirectory(config lm.ExposeDirectoryConfig) {
	communication.PrintWelcomeMessage()

	publicKeyAuthMethod, publicKey := parsePublicKey(config.Remote.IdentityFile)
	siteID := registerDomain(config.Remote.APIEndpoint.URI(), &publicKey, config.Remote.SiteID)
	server := getStaticFileServer(config.Local.Path, siteID, config.Remote.BasicAuthUsername, config.Remote.BasicAuthPassword)

	forward(config.Remote, config.Display, publicKeyAuthMethod, siteID, server, config.Local.Path, []string{"https"})
}

// ForwardDirectoryViaWebdav is used to expose local directory via Webdav (upload and download)
func ForwardDirectoryViaWebdav(config lm.ExposeWebdavConfig) {
	communication.PrintWelcomeMessage()

	publicKeyAuthMethod, publicKey := parsePublicKey(config.Remote.IdentityFile)
	siteID := registerDomain(config.Remote.APIEndpoint.URI(), &publicKey, config.Remote.SiteID)
	server := getWebdavServer(config.Local.Path, siteID, config.Remote.BasicAuthUsername, config.Remote.BasicAuthPassword)

	forward(config.Remote, config.Display, publicKeyAuthMethod, siteID, server, config.Local.Path, []string{"https", "davs", "webdav"})
}

func forward(remoteEndpointSpecs lm.RemoteEndpointSpecs, displayOptions lm.DisplayOptions,
	authMethod ssh.AuthMethod, siteID string, server *http.Server, localEndpoint string,
	protocols []string) {
	localListenerEndpoint := startLocalHTTPServer(server)

	serverSSHConnHTTPS := connectViaSSH(remoteEndpointSpecs.GatewayEndpoint, siteID, authMethod)
	defer serverSSHConnHTTPS.Close()
	listenerHTTPSOverSSH := listenOnRemoteEndpoint(serverSSHConnHTTPS)
	defer listenerHTTPSOverSSH.Close()

	go func() {
		log.Info().Msg("Issuing request to provision certificate")
		var netTransport = &http.Transport{
			Dial: (&net.Dialer{
				Timeout: 30 * time.Second,
			}).Dial,
			TLSHandshakeTimeout: 30 * time.Second,
		}
		var netClient = &http.Client{
			Timeout:   time.Second * 30,
			Transport: netTransport,
		}
		_, err := netClient.Get(urlmaker.GetSiteUrl("https", siteID))

		if err != nil {
			log.Error().Msg("TLS Certificate failed to provision. Will be obtained with first request made by any client, therefore first execution may be slower")
		} else {
			log.Info().Msg("TLS Certificate successfully provisioned")
		}
	}()

	communication.PrintTunnelSuccessMessage(siteID, protocols, localEndpoint, displayOptions.QR)

	for {
		client, err := listenerHTTPSOverSSH.Accept()
		if err == io.EOF {
			log.Info().Err(err).Msg("Connection dropped, reconnecting...")
			listenerHTTPSOverSSH.Close()
			serverSSHConnHTTPS = connectViaSSH(remoteEndpointSpecs.GatewayEndpoint, siteID, authMethod)
			defer serverSSHConnHTTPS.Close()
			listenerHTTPSOverSSH = listenOnRemoteEndpoint(serverSSHConnHTTPS)
			defer listenerHTTPSOverSSH.Close()
			continue
		} else if err != nil {
			log.Warn().Err(err).Msg("Failed to accept connection over HTTPS")
			continue
		}
		closehandler.SuccessfulConnectionOccured()
		go func() {
			log.Info().Msg("Succeeded to accept connection over HTTPS")
			local, err := net.Dial("tcp", localListenerEndpoint.URI())
			if err != nil {
				log.Fatal().Err(err).Msg("Dialing into local proxy for HTTPS failed")
			}
			if el := log.Debug(); el.Enabled() {
				el.Msg("Dialing into local proxy for HTTPS succeeded")
			}
			handleClient(client, local)
		}()
	}
}
