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

type TunnelType string

const (
	// HTTP specifies HTTP tunnel type
	HTTP TunnelType = "Tunnel_HTTP"
	// Directory specifies local directory tunnel type (download only)
	Directory TunnelType = "Tunnel_Directory"
	// WebDav specifies local directory tunnel type (download+upload via WebDav)
	WebDav TunnelType = "Tunnel_WebDav"
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

func registerDomain(apiURL string, publicKey *ssh.PublicKey, requestedSiteID, version string) string {
	communication.StartLoading("Registering your domain...")
	siteID, err := apiclient.RegisterSite(apiURL, *publicKey, requestedSiteID, version)
	if err != nil {
		communication.LoadingFailure()
		if requestErr, ok := err.(apiclient.RequestError); ok {
			log.Error().Int("status", requestErr.StatusCode).Msg("Request ended")
			log.Error().Msg(requestErr.Message)
			log.Error().Msg(fmt.Sprintf("Details: %s", requestErr.Details))
			communication.LogFatalMsg("Please fix the above issue and try again")
		} else {
			communication.LogFatalErr("Something unexpected happened, please let developers know", err)
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
			communication.LogInfo(fmt.Sprintf("SSH Connection failed, retrying in 10 seconds... (Attempt %d/%d)", i+1, sshRetries))
			time.Sleep(10 * time.Second)
		} else {
			sshSuccess = true
		}
	}
	if !sshSuccess {
		communication.WriteRed("An error occured while dialing into SSH. If your connection has been running for a while")
		communication.WriteRed("this might be caused by the server shutting down your connection.")
		communication.LogFatalErr("Dialing SSH Gateway for HTTPS failed.", err)
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
		communication.LogFatalErr("Something went wrong while creating server", err)
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
		communication.LogFatalErr("Failed to listen on TLS proxy for HTTPS", err)
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
			communication.LogFatalMsg("Failed to start TLS server")
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
		communication.LogFatalErr("Listening on remote endpoint for HTTPS failed", err)
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
		communication.LogFatalErr("No public key available", err)
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
		communication.LogFatalErr("Something went wrong while creating server", err)
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
		communication.LogFatalErr("Something went wrong while creating server", err)
	}
	return server
}

func listenOnRemoteEndpoint(serverSSHConnHTTPS *ssh.Client) *net.Listener {
	listenerHTTPSOverSSH, err := serverSSHConnHTTPS.Listen("tcp", remoteEndpoint.URI())
	if err != nil {
		communication.LoadingFailure()
		communication.LogFatalErr("Listening on remote endpoint for HTTPS failed", err)
	}
	return &listenerHTTPSOverSSH
}

// ForwardPort is used to forward external URL to locally available port
func ForwardPort(config lm.ExposeHTTPConfig, quitChannel <-chan bool) {
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
	siteID := registerDomain(config.Remote.APIEndpoint.URI(), &publicKey, config.Remote.SiteID, config.Display.Version)
	server := createTLSReverseProxy(localEndpoint, siteID, config.Remote.BasicAuthUsername, config.Remote.BasicAuthPassword, config.Display)

	forward(config.Remote, config.Display, publicKeyAuthMethod, siteID, server, localEndpoint.URI(), []string{"https"}, quitChannel)
}

// ForwardDirectory is used to expose local directory via HTTP (download only)
func ForwardDirectory(config lm.ExposeDirectoryConfig, quitChannel <-chan bool) {
	publicKeyAuthMethod, publicKey := parsePublicKey(config.Remote.IdentityFile)
	siteID := registerDomain(config.Remote.APIEndpoint.URI(), &publicKey, config.Remote.SiteID, config.Display.Version)
	server := getStaticFileServer(config.Local.Path, siteID, config.Remote.BasicAuthUsername, config.Remote.BasicAuthPassword)

	forward(config.Remote, config.Display, publicKeyAuthMethod, siteID, server, config.Local.Path, []string{"https"}, quitChannel)
}

// ForwardDirectoryViaWebdav is used to expose local directory via Webdav (upload and download)
func ForwardDirectoryViaWebdav(config lm.ExposeWebdavConfig, quitChannel <-chan bool) {
	publicKeyAuthMethod, publicKey := parsePublicKey(config.Remote.IdentityFile)
	siteID := registerDomain(config.Remote.APIEndpoint.URI(), &publicKey, config.Remote.SiteID, config.Display.Version)
	server := getWebdavServer(config.Local.Path, siteID, config.Remote.BasicAuthUsername, config.Remote.BasicAuthPassword)

	forward(config.Remote, config.Display, publicKeyAuthMethod, siteID, server, config.Local.Path, []string{"https", "davs", "webdav"}, quitChannel)
}

func forward(remoteEndpointSpecs lm.RemoteEndpointSpecs, displayOptions lm.DisplayOptions,
	authMethod ssh.AuthMethod, siteID string, server *http.Server, localEndpoint string,
	protocols []string, quitChannel <-chan bool) {
	localListenerEndpoint := startLocalHTTPServer(server)

	serverSSHConnHTTPS := connectViaSSH(remoteEndpointSpecs.GatewayEndpoint, siteID, authMethod)
	defer serverSSHConnHTTPS.Close()
	listenerHTTPSOverSSH := listenOnRemoteEndpoint(serverSSHConnHTTPS)
	defer (*listenerHTTPSOverSSH).Close()

	go func() {
		communication.LogInfo("Issuing request to provision certificate")
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
		_, err := netClient.Get(urlmaker.GetSiteURL("https", siteID))

		if err != nil {
			log.Error().Msg("TLS Certificate failed to provision. Will be obtained with first request made by any client, therefore first execution may be slower")
		} else {
			communication.LogInfo("TLS Certificate successfully provisioned")
		}
	}()

	communication.PrintTunnelSuccessMessage(siteID, protocols, localEndpoint, displayOptions.QR)

	acceptedClients := make(chan net.Conn)
	tunnelTerminatedOnPurpose := false

	go func(l *net.Listener, tunnelTerminatedOnPurpose *bool) {
		for {
			log.Debug().Msg("Waiting to accept")
			client, err := (*l).Accept()
			log.Debug().Msg("Accepted")
			if err == io.EOF {
				if !(*tunnelTerminatedOnPurpose) {
					communication.LogInfo(err.Error() + " Connection dropped, reconnecting...")
					(*l).Close()
					serverSSHConnHTTPS = connectViaSSH(remoteEndpointSpecs.GatewayEndpoint, siteID, authMethod)
					defer serverSSHConnHTTPS.Close()
					l = listenOnRemoteEndpoint(serverSSHConnHTTPS)
					defer (*l).Close()
					continue
				}
			} else if err != nil {
				communication.LogWarnErr("Failed to accept connection over HTTPS", err)
				continue
			}
			log.Debug().Msg("Sending client trough channel")
			acceptedClients <- client
		}
	}(listenerHTTPSOverSSH, &tunnelTerminatedOnPurpose)

	for {
		log.Debug().Msg("For loop cycle")
		select {
		case <-quitChannel:
			tunnelTerminatedOnPurpose = true
			communication.PrintGoodbyeMessage()
			return
		case client := <-acceptedClients:
			log.Debug().Msg("Handling client")
			closehandler.SuccessfulConnectionOccured()
			go func() {
				communication.LogInfo("Succeeded to accept connection over HTTPS")
				local, err := net.Dial("tcp", localListenerEndpoint.URI())
				if err != nil {
					communication.LogFatalErr("Dialing into local proxy for HTTPS failed", err)
				}
				if el := log.Debug(); el.Enabled() {
					el.Msg("Dialing into local proxy for HTTPS succeeded")
				}
				handleClient(client, local)
			}()
		}
	}
}
