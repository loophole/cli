package loophole

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/briandowns/spinner"
	"github.com/logrusorgru/aurora"
	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/loophole/cli/internal/pkg/client"
	"github.com/loophole/cli/internal/pkg/token"
	"github.com/mattn/go-colorable"
	"github.com/mdp/qrterminal"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	apiURL = "https://api.loophole.cloud"
)

// remote forwarding port (on remote SSH server network)
var remoteEndpoint = lm.Endpoint{
	Host: "127.0.0.1",
	Port: 80,
}

var colorableOutput = colorable.NewColorableStdout()

func parsePublicKey(file string) (ssh.AuthMethod, ssh.PublicKey, error) {
	key, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, nil, err
	}

	var passwordError *ssh.PassphraseMissingError
	signer, err := ssh.ParsePrivateKey(key)

	if err != nil {
		if errors.As(err, &passwordError) {
			fmt.Fprint(colorableOutput, "Enter SSH password:")
			password, _ := terminal.ReadPassword(int(os.Stdin.Fd()))
			fmt.Println()
			signer, err = ssh.ParsePrivateKeyWithPassphrase(key, []byte(password))
			if err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	}

	return ssh.PublicKeys(signer), signer.PublicKey(), nil
}

// From https://sosedoff.com/2015/05/25/ssh-port-forwarding-with-go.html
// Handle local client connections and tunnel data to the remote serverq
// Will use io.Copy - http://golang.org/pkg/io/#Copy
func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	// Start remote -> local data transfer
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			if el := log.Debug(); el.Enabled() {
				el.Err(err).Msg("Error copying remote->local:")
			}
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			if el := log.Debug(); el.Enabled() {
				el.Err(err).Msg("Error copying local->remote:")
			}
		}
		chDone <- true
	}()

	<-chDone
}

func printWelcomeMessage() {
	fmt.Fprint(colorableOutput, aurora.Cyan("Loophole"))
	fmt.Fprint(colorableOutput, aurora.Italic(" - End to end TLS encrypted TCP communication between you and your clients"))
	fmt.Println()
	fmt.Println()
}

func startLoading(loader *spinner.Spinner, message string) {
	if el := log.Debug(); !el.Enabled() {
		loader.Prefix = fmt.Sprintf("%s ", message)
		loader.Start()
	} else {
		fmt.Println(message)
	}
}

func loadingSuccess(loader *spinner.Spinner) {
	if el := log.Debug(); !el.Enabled() {
		loader.FinalMSG = fmt.Sprintf("%s%s\n", loader.Prefix, aurora.Green("Success!"))
		loader.Stop()
	}
}

func loadingFailure(loader *spinner.Spinner) {
	if el := log.Debug(); !el.Enabled() {
		loader.FinalMSG = fmt.Sprintf("%s%s\n", loader.Prefix, aurora.Red("Error!"))
		loader.Stop()
	}
}

func generateListener(config lm.Config, publicKeyAuthMethod *ssh.AuthMethod, publicKey *ssh.PublicKey, siteSpecs client.SiteSpecification) (net.Listener, *lm.Endpoint, client.SiteSpecification) {

	loader := spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(colorable.NewColorableStdout()))

	localEndpoint := lm.Endpoint{
		Host: config.Host,
		Port: config.Port,
	}

	if el := log.Debug(); el.Enabled() {
		el.Msg("Checking public key availability")
	}

	var err error
	if *publicKey == nil {
		*publicKeyAuthMethod, *publicKey, err = parsePublicKey(config.IdentityFile)
		if err != nil {
			log.Fatal().Err(err).Msg("No public key available")
		}
	}

	startLoading(loader, "Registering your domain...")

	if el := log.Debug(); el.Enabled() {
		fmt.Println()
		el.Msg("Registering site")
	}

	if siteSpecs.ResultCode != 0 { //checking whether siteSpecshas been used yet
		log.Info().Msg("Trying to reuse old hostname...")
	} else {
		siteSpecs, err = client.RegisterSite(apiURL, *publicKey, config.SiteID)
		if err != nil {
			if siteSpecs.ResultCode == 401 {
				if el := log.Debug(); el.Enabled() {
					fmt.Println()
					el.Err(err).Msg("Failed to register site")
				}
				if el := log.Debug(); el.Enabled() {
					el.Msg("Trying to refresh token")
				}
				if err := token.RefreshToken(); err != nil {
					loadingFailure(loader)
					log.Fatal().Err(err).Msg("Failed to refresh token, try logging in again")
				}
				siteSpecs, err = client.RegisterSite(apiURL, *publicKey, config.SiteID)
				if err != nil {
					loadingFailure(loader)
					log.Fatal().Err(err).Msg("Failed to register site, try logging in again")
				}
			} else if siteSpecs.ResultCode == 403 {
				log.Fatal().Err(err).Msg("You don't have required permissions to establish tunnel with given parameters")
			} else if siteSpecs.ResultCode == 600 || siteSpecs.ResultCode == 601 {
				log.Fatal().Err(err).Msg("Looks like you're not logged in")
			} else {
				log.Fatal().Err(err).Msg("Something unexpected happened, please let developers know")
			}
		}
	}
	loadingSuccess(loader)

	sshConfigHTTPS := &ssh.ClientConfig{
		User: fmt.Sprintf("%s_https", siteSpecs.SiteID),
		Auth: []ssh.AuthMethod{
			*publicKeyAuthMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	startLoading(loader, "Initializing secure tunnel... ")

	if el := log.Debug(); el.Enabled() {
		fmt.Println()
		el.Msg("Dialing gateway to establish the tunnel..")
	}
	serverSSHConnHTTPS, err := ssh.Dial("tcp", config.GatewayEndpoint.String(), sshConfigHTTPS)
	if err != nil {
		loadingFailure(loader)
		fmt.Fprintln(colorableOutput, aurora.BgRed("An error occured while dialing into SSH. If your connection has been running for a while, this might be caused by the server shutting down your connection."))
		log.Fatal().Err(err).Msg("Dialing SSH Gateway for HTTPS failed.")
	}
	if el := log.Debug(); el.Enabled() {
		fmt.Println()
		el.Msg("Dialing SSH Gateway for HTTPS succeeded")
	}
	loadingSuccess(loader)

	startLoading(loader, "Obtaining TLS certificate provider... ")

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(fmt.Sprintf("%s.loophole.site", siteSpecs.SiteID), "abc.loophole.site"), //Your domain here
		Cache:      autocert.DirCache(cache.GetLocalStorageDir("certs")),                                           //Folder for storing certificates
		Email:      fmt.Sprintf("%s@loophole.main.dev", siteSpecs.SiteID),
	}
	if el := log.Debug(); el.Enabled() {
		fmt.Println()
		el.Msg("Cert Manager created")
	}

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   localEndpoint.String(),
	})
	if el := log.Debug(); el.Enabled() {
		el.
			Str("target", localEndpoint.String()).
			Msg("Proxy via http created")
	}
	server := &http.Server{
		Handler:   proxy,
		TLSConfig: certManager.TLSConfig(),
	}
	loadingSuccess(loader)

	startLoading(loader, "Starting the server... ")

	if el := log.Debug(); el.Enabled() {
		fmt.Println()
		el.Msg("Server for proxy created")
	}
	proxyListenerHTTPS, err := net.Listen("tcp", ":0")
	if err != nil {
		loadingFailure(loader)
		log.Fatal().Err(err).Msg("Failed to listen on TLS proxy for HTTPS")
	}
	if el := log.Debug(); el.Enabled() {
		el.
			Int("port", proxyListenerHTTPS.Addr().(*net.TCPAddr).Port).
			Msg("Proxy listener for HTTPS started")
	}
	go func() {
		err := server.ServeTLS(proxyListenerHTTPS, "", "")
		if err != nil {
			loadingFailure(loader)
			log.Fatal().Msg("Failed to start TLS server")
		}
	}()
	if el := log.Debug(); el.Enabled() {
		el.Msg("Started server TLS server")
	}
	listenerHTTPSOverSSH, err := serverSSHConnHTTPS.Listen("tcp", remoteEndpoint.String())
	if err != nil {
		loadingFailure(loader)
		log.Fatal().Err(err).Msg("Listening on remote endpoint for HTTPS failed")
	}
	if el := log.Debug(); el.Enabled() {
		fmt.Println()
		el.Msg("Listening on remote endpoint for HTTPS succeeded")
	}

	loadingSuccess(loader)

	proxiedEndpointHTTPS := &lm.Endpoint{
		Host: "127.0.0.1",
		Port: int32(proxyListenerHTTPS.Addr().(*net.TCPAddr).Port),
	}

	fmt.Println()
	fmt.Fprint(colorableOutput, "Forwarding ")
	fmt.Fprint(colorableOutput, aurora.Green(fmt.Sprintf("http://%s.loophole.site", siteSpecs.SiteID)))
	fmt.Fprint(colorableOutput, " -> ")
	fmt.Fprint(colorableOutput, aurora.Green(fmt.Sprintf("%s:%d", config.Host, config.Port)))
	fmt.Println()
	fmt.Fprint(colorableOutput, "Forwarding ")
	fmt.Fprint(colorableOutput, aurora.Green(fmt.Sprintf("https://%s.loophole.site", siteSpecs.SiteID)))
	fmt.Fprint(colorableOutput, " -> ")
	fmt.Fprint(colorableOutput, aurora.Green(fmt.Sprintf("%s:%d", config.Host, config.Port)))
	fmt.Println()
	if config.QR {
		QRconfig := qrterminal.Config{
			Level:     qrterminal.L,
			Writer:    colorableOutput,
			BlackChar: qrterminal.WHITE,
			WhiteChar: qrterminal.BLACK,
			QuietZone: 1,
		}
		qrterminal.GenerateWithConfig(fmt.Sprintf("http://%s.loophole.site", siteSpecs.SiteID), QRconfig)
	}
	fmt.Fprint(colorableOutput, fmt.Sprintf("%s", aurora.Italic("TLS Certificate will be obtained with first request to the above address, therefore first execution may be slower\n")))
	fmt.Println()
	fmt.Fprint(colorableOutput, fmt.Sprintf("%s", aurora.Cyan("Press CTRL + C to stop the service\n")))
	fmt.Println()
	fmt.Fprint(colorableOutput, fmt.Sprint("Logs:\n"))

	log.Info().Msg("Awaiting connections...")
	return listenerHTTPSOverSSH, proxiedEndpointHTTPS, siteSpecs
}

// Start starts the tunnel on specified host and port
func Start(config lm.Config) {
	printWelcomeMessage()

	var publicKeyAuthMethod *ssh.AuthMethod = new(ssh.AuthMethod)
	var publicKey *ssh.PublicKey = new(ssh.PublicKey)
	var siteSpecs client.SiteSpecification

	listenerHTTPSOverSSH, proxiedEndpointHTTPS, siteSpecs := generateListener(config, publicKeyAuthMethod, publicKey, siteSpecs)
	defer listenerHTTPSOverSSH.Close()

	for {
		client, err := listenerHTTPSOverSSH.Accept()
		if err == io.EOF {
			log.Info().Err(err).Msg("Connection dropped, reconnecting...")
			listenerHTTPSOverSSH.Close()
			listenerHTTPSOverSSH, _, _ = generateListener(config, publicKeyAuthMethod, publicKey, siteSpecs)
			continue
		} else if err != nil {
			log.Info().Err(err).Msg("Failed to accept connection over HTTPS")
			continue
		}
		go func() {
			log.Info().Msg("Succeeded to accept connection over HTTPS")
			// Open a (local) connection to proxiedEndpointHTTPS whose content will be forwarded to serverEndpoint
			local, err := net.Dial("tcp", proxiedEndpointHTTPS.String())
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
