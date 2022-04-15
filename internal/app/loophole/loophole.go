package loophole

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/loophole/cli/config"
	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/apiclient"
	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/loophole/cli/internal/pkg/httpserver"
	"github.com/loophole/cli/internal/pkg/keys"
	"github.com/loophole/cli/internal/pkg/urlmaker"
	"golang.org/x/crypto/ssh"
)

// TunnelType is used to define supported tunnel types
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

func handleClient(tunnelID string, client net.Conn, local net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	// Start local -> client data transfer
	go func() {
		nob, err := io.Copy(client, local)
		communication.TunnelDebug(tunnelID, fmt.Sprintf("Transfered out %d bytes", nob))
		if err != nil {
			if err != io.EOF {
				communication.TunnelWarn(tunnelID, fmt.Sprintf("Error copying local -> client: %s", err.Error()))
			} else {
				communication.TunnelDebug(tunnelID, fmt.Sprintf("Error copying local -> client: %s", err.Error()))
			}
		}
		chDone <- true
	}()

	// Start client -> local data transfer
	go func() {
		nob, err := io.Copy(local, client)
		communication.TunnelDebug(tunnelID, fmt.Sprintf("Received %d bytes", nob))
		if err != nil {
			if err != io.EOF {
				communication.TunnelWarn(tunnelID, fmt.Sprintf("Error copying client -> local: %s", err.Error()))
			} else {
				communication.TunnelDebug(tunnelID, fmt.Sprintf("Error copying client -> local: %s", err.Error()))
			}
		}
		chDone <- true
	}()

	<-chDone
}

func registerDomain(publicKey *ssh.PublicKey, requestedSiteID string, tunnelID string) (*apiclient.RegistrationSuccessResponse, error) {
	communication.LoadingStart(tunnelID, "Registering your domain...")
	registrationResult, err := apiclient.RegisterSite(*publicKey, requestedSiteID)
	if err != nil {
		communication.LoadingFailure(tunnelID, err)
		if requestErr, ok := err.(apiclient.RequestError); ok {
			communication.TunnelError(tunnelID, fmt.Sprintf("Request ended with status code %d", requestErr.StatusCode))
			communication.TunnelError(tunnelID, requestErr.Message)
			communication.TunnelError(tunnelID, fmt.Sprintf("Details: %s", requestErr.Details))
			communication.TunnelError(tunnelID, "Please fix the above issue and try again")
		} else {
			communication.TunnelError(tunnelID, "Something unexpected happened, please let developers know")
		}
		return nil, err
	}
	communication.LoadingSuccess(tunnelID)
	return registrationResult, nil
}

func connectViaSSH(siteID string, tunnelID string, authMethod ssh.AuthMethod) (*ssh.Client, error) {
	var serverSSHConnHTTPS *ssh.Client
	sshConfigHTTPS := &ssh.ClientConfig{
		User: siteID,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	var sshSuccess bool = false
	var sshRetries int = 5
	var err error
	for i := 0; i < sshRetries && !sshSuccess; i++ { // Connection retries in case of reconnect during gateway shutdown
		communication.LoadingStart(tunnelID, "Initializing secure tunnel... ")
		serverSSHConnHTTPS, err = ssh.Dial("tcp", config.Config.GatewayEndpoint.Hostname(), sshConfigHTTPS)
		if err != nil {
			communication.LoadingFailure(tunnelID, err)
			communication.TunnelInfo(tunnelID, fmt.Sprintf("SSH Connection failed, retrying in 10 seconds... (Attempt %d/%d)", i+1, sshRetries))
			time.Sleep(10 * time.Second)
		} else {
			sshSuccess = true
		}
	}
	if !sshSuccess {
		communication.TunnelError(tunnelID, "An error occured while dialing into SSH. If your connection has been running for a while, "+
			"this might be caused by the server shutting down your connection. Dialing SSH Gateway for HTTPS failed.")
		err = fmt.Errorf("failed %d", 5)
		return nil, err
	}
	communication.TunnelDebug(tunnelID, "Dialing SSH Gateway for HTTPS succeeded")
	communication.LoadingSuccess(tunnelID)
	return serverSSHConnHTTPS, nil
}

func createTLSReverseProxy(localEndpoint lm.Endpoint, remoteConfig lm.RemoteEndpointSpecs) (*http.Server, error) {
	communication.LoadingStart(remoteConfig.TunnelID, "Starting local TLS proxy server")
	serverBuilder := httpserver.New().
		WithSiteID(remoteConfig.SiteID).
		WithDomain(remoteConfig.Domain).
		DisableOldCiphers(remoteConfig.DisableOldCiphers).
		Proxy().
		ToEndpoint(localEndpoint)

	if remoteConfig.BasicAuthUsername != "" && remoteConfig.BasicAuthPassword != "" {
		serverBuilder = serverBuilder.
			WithBasicAuth(remoteConfig.BasicAuthUsername, remoteConfig.BasicAuthPassword)
	}
	if remoteConfig.DisableProxyErrorPage {
		serverBuilder = serverBuilder.
			DisableProxyErrorPage()
	}
	if localEndpoint.Protocol == "https" {
		serverBuilder = serverBuilder.
			EnableInsecureHTTPSBackend()
	}

	communication.TunnelDebug(remoteConfig.TunnelID, fmt.Sprintf("Proxy via http to %s created", localEndpoint.URI()))
	server, err := serverBuilder.Build()
	if err != nil {
		communication.LoadingFailure(remoteConfig.TunnelID, err)
		communication.TunnelError(remoteConfig.TunnelID, "Something went wrong while creating server")
		communication.TunnelStartFailure(remoteConfig.TunnelID, err)
		return nil, err
	}
	return server, nil
}

func startLocalHTTPServer(tunnelID string, server *http.Server) (*lm.Endpoint, error) {
	communication.LoadingStart(tunnelID, "Starting local proxy server... ")

	communication.TunnelDebug(tunnelID, "Server for proxy created")
	localListener, err := net.Listen("tcp", ":0")
	if err != nil {
		communication.LoadingFailure(tunnelID, err)
		communication.TunnelError(tunnelID, "Failed to listen on TLS proxy for HTTPS")
		return nil, err
	}
	localListenerEndpoint := &lm.Endpoint{
		Host: "127.0.0.1",
		Port: int32(localListener.Addr().(*net.TCPAddr).Port),
	}
	communication.TunnelDebug(tunnelID, fmt.Sprintf("Proxy listener for HTTPS started on port %d", localListenerEndpoint.Port))
	go func() {
		err := server.ServeTLS(localListener, "", "")
		if err != nil {
			communication.LoadingFailure(tunnelID, err)
			communication.TunnelStartFailure(tunnelID, err)
		}
	}()
	communication.TunnelDebug(tunnelID, "Started server TLS server")
	communication.LoadingSuccess(tunnelID)
	return localListenerEndpoint, nil
}

func startRemoteForwardServer(tunnelID string, serverSSHConnHTTPS *ssh.Client) (net.Listener, error) {
	listenerHTTPSOverSSH, err := serverSSHConnHTTPS.Listen("tcp", remoteEndpoint.URI())
	if err != nil {
		communication.LoadingFailure(tunnelID, err)
		communication.TunnelError(tunnelID, "Listening on remote endpoint for HTTPS failed")
		return nil, err
	}
	communication.TunnelDebug(tunnelID, "Listening on remote endpoint for HTTPS succeeded")
	return listenerHTTPSOverSSH, nil
}

func parsePublicKey(tunnelID string, identityFile string) (ssh.AuthMethod, ssh.PublicKey, error) {
	publicKeyAuthMethod, publicKey, err := keys.ParsePublicKey(identityFile)
	if err != nil {
		communication.LoadingFailure(tunnelID, err)
		communication.TunnelError(tunnelID, "No public key available")
		return nil, nil, err
	}

	return publicKeyAuthMethod, publicKey, nil
}

func getStaticFileServer(exposeDirectoryConfig lm.ExposeDirectoryConfig) (*http.Server, error) {
	communication.LoadingStart(exposeDirectoryConfig.Remote.TunnelID, "Starting local file server")
	serverBuilder := httpserver.New().
		WithSiteID(exposeDirectoryConfig.Remote.SiteID).
		WithDomain(exposeDirectoryConfig.Remote.Domain).
		DisableOldCiphers(exposeDirectoryConfig.Remote.DisableOldCiphers).
		ServeStatic().
		FromDirectory(exposeDirectoryConfig.Local.Path)

	if exposeDirectoryConfig.Remote.BasicAuthUsername != "" && exposeDirectoryConfig.Remote.BasicAuthPassword != "" {
		serverBuilder = serverBuilder.
			WithBasicAuth(exposeDirectoryConfig.Remote.BasicAuthUsername, exposeDirectoryConfig.Remote.BasicAuthPassword)
	}

	communication.LoadingSuccess(exposeDirectoryConfig.Remote.TunnelID)
	server, err := serverBuilder.Build()
	if err != nil {
		communication.LoadingFailure(exposeDirectoryConfig.Remote.TunnelID, err)
		communication.TunnelError(exposeDirectoryConfig.Remote.TunnelID, "Something went wrong while creating server")
		return nil, err
	}
	return server, nil
}

func getWebdavServer(exposeWebDavConfig lm.ExposeWebdavConfig) (*http.Server, error) {
	communication.LoadingStart(exposeWebDavConfig.Remote.TunnelID, "Starting WebDav server")
	serverBuilder := httpserver.New().
		WithSiteID(exposeWebDavConfig.Remote.SiteID).
		WithDomain(exposeWebDavConfig.Remote.Domain).
		DisableOldCiphers(exposeWebDavConfig.Remote.DisableOldCiphers).
		ServeWebdav().
		FromDirectory(exposeWebDavConfig.Local.Path)

	if exposeWebDavConfig.Remote.BasicAuthUsername != "" && exposeWebDavConfig.Remote.BasicAuthPassword != "" {
		serverBuilder = serverBuilder.
			WithBasicAuth(exposeWebDavConfig.Remote.BasicAuthUsername, exposeWebDavConfig.Remote.BasicAuthPassword)
	}

	communication.LoadingSuccess(exposeWebDavConfig.Remote.TunnelID)
	server, err := serverBuilder.Build()
	if err != nil {
		communication.LoadingFailure(exposeWebDavConfig.Remote.TunnelID, err)
		communication.TunnelError(exposeWebDavConfig.Remote.TunnelID, "Something went wrong while creating server")
		return nil, err
	}
	return server, nil
}

func listenOnRemoteEndpoint(tunnelID string, serverSSHConnHTTPS *ssh.Client) (*net.Listener, error) {
	listenerHTTPSOverSSH, err := serverSSHConnHTTPS.Listen("tcp", remoteEndpoint.URI())
	if err != nil {
		communication.LoadingFailure(tunnelID, err)
		communication.TunnelError(tunnelID, "Listening on remote endpoint for HTTPS failed")
		return nil, err
	}
	return &listenerHTTPSOverSSH, nil
}

// RegisterTunnel is used to register tunnel in loophole API and grant user access to connect to it
func RegisterTunnel(remoteConfig *lm.RemoteEndpointSpecs) (ssh.AuthMethod, error) {
	publicKeyAuthMethod, publicKey, err := parsePublicKey(remoteConfig.TunnelID, remoteConfig.IdentityFile)
	if err != nil {
		return nil, err
	}
	registrationResult, err := registerDomain(&publicKey, remoteConfig.SiteID, remoteConfig.TunnelID)
	if err != nil {
		return nil, err
	}
	remoteConfig.SiteID = registrationResult.SiteID
	remoteConfig.Domain = registrationResult.Domain
	communication.TunnelStart(remoteConfig.TunnelID)

	return publicKeyAuthMethod, nil
}

// ForwardPort is used to forward external URL to locally available port
func ForwardPort(exposeHTTPConfig lm.ExposeHTTPConfig, publicKeyAuthMethod ssh.AuthMethod, quitChannel <-chan bool) error {
	protocol := "http"
	if exposeHTTPConfig.Local.HTTPS {
		protocol = "https"
	}
	localEndpoint := lm.Endpoint{
		Protocol: protocol,
		Host:     exposeHTTPConfig.Local.Host,
		Port:     exposeHTTPConfig.Local.Port,
		Path:     exposeHTTPConfig.Local.Path,
	}

	server, err := createTLSReverseProxy(localEndpoint, exposeHTTPConfig.Remote)
	if err != nil {
		return err
	}
	return forward(exposeHTTPConfig.Remote, publicKeyAuthMethod, server, localEndpoint.URI(), []string{"https"}, quitChannel)
}

// ForwardDirectory is used to expose local directory via HTTP (download only)
func ForwardDirectory(exposeDirectoryConfig lm.ExposeDirectoryConfig, publicKeyAuthMethod ssh.AuthMethod, quitChannel <-chan bool) error {
	server, err := getStaticFileServer(exposeDirectoryConfig)
	if err != nil {
		return err
	}
	return forward(exposeDirectoryConfig.Remote, publicKeyAuthMethod, server, exposeDirectoryConfig.Local.Path, []string{"https"}, quitChannel)
}

// ForwardDirectoryViaWebdav is used to expose local directory via Webdav (upload and download)
func ForwardDirectoryViaWebdav(exposeWebdavConfig lm.ExposeWebdavConfig, publicKeyAuthMethod ssh.AuthMethod, quitChannel <-chan bool) error {
	server, err := getWebdavServer(exposeWebdavConfig)
	if err != nil {
		return err
	}

	return forward(exposeWebdavConfig.Remote, publicKeyAuthMethod, server, exposeWebdavConfig.Local.Path, []string{"https", "davs", "webdav"}, quitChannel)
}

func forward(remoteEndpointSpecs lm.RemoteEndpointSpecs,
	authMethod ssh.AuthMethod, server *http.Server, localEndpoint string,
	protocols []string, quitChannel <-chan bool) error {

	localListenerEndpoint, err := startLocalHTTPServer(remoteEndpointSpecs.TunnelID, server)
	if err != nil {
		communication.TunnelStartFailure(remoteEndpointSpecs.TunnelID, err)
		return err
	}
	serverSSHConnHTTPS, err := connectViaSSH(remoteEndpointSpecs.SiteID, remoteEndpointSpecs.TunnelID, authMethod)
	if err != nil {
		communication.TunnelStartFailure(remoteEndpointSpecs.TunnelID, err)
		return err
	} else {
		defer serverSSHConnHTTPS.Close()
	}
	listenerHTTPSOverSSH, err := listenOnRemoteEndpoint(remoteEndpointSpecs.TunnelID, serverSSHConnHTTPS)
	if err != nil {
		communication.TunnelStartFailure(remoteEndpointSpecs.TunnelID, err)
		return err
	} else {
		defer (*listenerHTTPSOverSSH).Close()
	}

	go func() {
		communication.TunnelDebug(remoteEndpointSpecs.TunnelID, "Issuing request to provision certificate")
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
		_, err := netClient.Get(urlmaker.GetSiteURL("https", remoteEndpointSpecs.SiteID, remoteEndpointSpecs.Domain))

		if err != nil {
			communication.TunnelError(remoteEndpointSpecs.TunnelID, "TLS Certificate failed to provision. Will be obtained with first request made by any client, therefore first execution may be slower")
		} else {
			communication.TunnelInfo(remoteEndpointSpecs.TunnelID, "TLS Certificate successfully provisioned")
		}
	}()

	communication.TunnelStartSuccess(remoteEndpointSpecs, localEndpoint)

	acceptedClients := make(chan net.Conn)
	tunnelTerminatedOnPurpose := false

	go func(l *net.Listener, tunnelTerminatedOnPurpose *bool) {
		for {
			communication.TunnelDebug(remoteEndpointSpecs.TunnelID, "Waiting to accept")
			client, err := (*l).Accept()
			communication.TunnelDebug(remoteEndpointSpecs.TunnelID, "Accepted")
			if err == io.EOF {
				if !(*tunnelTerminatedOnPurpose) {
					communication.TunnelInfo(remoteEndpointSpecs.TunnelID, err.Error()+" Connection dropped, reconnecting...")
					(*l).Close()
					serverSSHConnHTTPS, err = connectViaSSH(remoteEndpointSpecs.SiteID, remoteEndpointSpecs.TunnelID, authMethod)
					if err != nil {
						communication.TunnelStartFailure(remoteEndpointSpecs.TunnelID, err)
						return
					} else {
						defer serverSSHConnHTTPS.Close()
					}
					defer serverSSHConnHTTPS.Close()
					l, err = listenOnRemoteEndpoint(remoteEndpointSpecs.TunnelID, serverSSHConnHTTPS)
					if err != nil {
						communication.TunnelStartFailure(remoteEndpointSpecs.TunnelID, err)
						return
					} else {
						defer (*l).Close()
					}
					continue
				}
			} else if err != nil {
				communication.TunnelWarn(remoteEndpointSpecs.TunnelID, "Failed to accept connection over HTTPS")
				communication.TunnelWarn(remoteEndpointSpecs.TunnelID, err.Error())
				continue
			}
			communication.TunnelDebug(remoteEndpointSpecs.TunnelID, "Sending client trough channel")
			acceptedClients <- client
		}
	}(listenerHTTPSOverSSH, &tunnelTerminatedOnPurpose)

	for {
		communication.TunnelDebug(remoteEndpointSpecs.TunnelID, "For loop cycle")
		select {
		case <-quitChannel:
			tunnelTerminatedOnPurpose = true
			communication.TunnelStopSuccess(remoteEndpointSpecs.TunnelID)
			return nil
		case client := <-acceptedClients:
			communication.TunnelDebug(remoteEndpointSpecs.TunnelID, "Handling client")
			go func() {
				communication.TunnelInfo(remoteEndpointSpecs.TunnelID, "Succeeded to accept connection over HTTPS")
				communication.TunnelDebug(remoteEndpointSpecs.TunnelID, fmt.Sprintf("Dialing into local proxy for HTTPS: %s", localListenerEndpoint.URI()))
				local, err := net.Dial("tcp", localListenerEndpoint.URI())
				if err != nil {
					communication.TunnelError(remoteEndpointSpecs.TunnelID, "Dialing into local proxy for HTTPS failed")
				}
				communication.TunnelDebug(remoteEndpointSpecs.TunnelID, "Dialing into local proxy for HTTPS succeeded")
				handleClient(remoteEndpointSpecs.TunnelID, client, local)
			}()
		}
	}
}
