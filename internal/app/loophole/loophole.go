package loophole

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
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
	"github.com/kyokomi/emoji"
	"github.com/logrusorgru/aurora"
	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/loophole/cli/internal/pkg/token"
	"github.com/mattn/go-colorable"
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

func registerSite(apiURL string, publicKey ssh.PublicKey, siteID string) (string, error) {
	publicKeyString := publicKey.Type() + " " + base64.StdEncoding.EncodeToString(publicKey.Marshal())

	if !token.IsTokenSaved() {
		return "", fmt.Errorf("Please log in before using Loophole")
	}

	accessToken, err := token.GetAccessToken()
	if err != nil {
		return "", fmt.Errorf("There was a problem reading token")
	}

	data := map[string]string{
		"key": publicKeyString,
	}
	if siteID != "" {
		data["id"] = siteID
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("There was a problem encoding request body")
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/site", apiURL), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("There was a problem creating request body")
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", accessToken))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("Site registration request ended with %d status and no message", resp.StatusCode)
		}
		return "", fmt.Errorf("Site registration request ended with %d status and message: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

	if el := log.Debug(); el.Enabled() {
		fmt.Println()
		el.Interface("result", result).Msg("Response")
	}
	defer resp.Body.Close()

	siteID, ok := result["id"].(string)
	if !ok {
		return "", fmt.Errorf("Error converting siteId to string")
	}
	return siteID, nil
}

func printWelcomeMessage() {
	fmt.Fprint(colorableOutput, aurora.Cyan("Loophole"))
	fmt.Fprint(colorableOutput, aurora.Italic(" - End to end TLS encrypted TCP communication between you and your clients"))
	fmt.Println()
	fmt.Println()
}

func startLoading(loader *spinner.Spinner, message string) {
	if el := log.Debug(); !el.Enabled() {
		loader.Prefix = emoji.Sprintf("%s ", message)
		loader.Start()
	} else {
		fmt.Println(emoji.Sprint(message))
	}
}

func loadingSuccess(loader *spinner.Spinner) {
	if el := log.Debug(); !el.Enabled() {
		loader.FinalMSG = emoji.Sprintf("%s%s\n", loader.Prefix, aurora.Green(":check_mark:"))
		loader.Stop()
	} else {
		fmt.Println(emoji.Sprint(loader.Prefix))
	}
}

func loadingFailure(loader *spinner.Spinner) {
	if el := log.Debug(); !el.Enabled() {
		loader.FinalMSG = emoji.Sprintf("%s%s\n", loader.Prefix, aurora.Red(":cross_mark:"))
		loader.Stop()
	} else {
		fmt.Fprintln(colorableOutput, emoji.Sprint(loader.Prefix))
	}
}

func generateListener(config lm.Config, publicKeyAuthMethod *ssh.AuthMethod, publicKey *ssh.PublicKey) (net.Listener, *lm.Endpoint) {

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

	startLoading(loader, ":wave: Registering your domain...")

	if el := log.Debug(); el.Enabled() {
		fmt.Println()
		el.Msg("Registering site")
	}
	siteID, err := registerSite(apiURL, *publicKey, config.SiteID)
	if err != nil {
		if el := log.Debug(); el.Enabled() {
			fmt.Println()
			el.Err(err).Msg("Failed to register site")
		}
		if el := log.Debug(); el.Enabled() {
			el.Msg("Trying to refresh token")
		}
		if err := token.RefreshToken(); err != nil {
			loadingFailure(loader)
			log.Fatal().Err(err).Msg("Failed to refresh token")
		}
		siteID, err = registerSite(apiURL, *publicKey, config.SiteID)
		if err != nil {
			loadingFailure(loader)
			log.Fatal().Err(err).Msg("Failed to register site")
		}
	}
	loadingSuccess(loader)

	sshConfigHTTPS := &ssh.ClientConfig{
		User: fmt.Sprintf("%s_https", siteID),
		Auth: []ssh.AuthMethod{
			*publicKeyAuthMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	startLoading(loader, ":lock: Initializing secure tunnel... ")

	if el := log.Debug(); el.Enabled() {
		fmt.Println()
		el.Msg("Dialing gateway to establish the tunnel..")
	}
	serverSSHConnHTTPS, err := ssh.Dial("tcp", config.GatewayEndpoint.String(), sshConfigHTTPS)
	if err != nil {
		loadingFailure(loader)
		log.Fatal().Err(err).Msg("Dialing SSH Gateway for HTTPS failed")
	}
	if el := log.Debug(); el.Enabled() {
		fmt.Println()
		el.Msg("Dialing SSH Gateway for HTTPS succeeded")
	}
	loadingSuccess(loader)

	startLoading(loader, ":key: Obtaining TLS certificate provider... ")

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(fmt.Sprintf("%s.loophole.site", siteID), "abc.loophole.site"), //Your domain here
		Cache:      autocert.DirCache(cache.GetLocalStorageDir("certs")),                                 //Folder for storing certificates
		Email:      fmt.Sprintf("%s@loophole.main.dev", siteID),
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

	startLoading(loader, ":cloud:  Starting the server... ")

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
	fmt.Fprint(colorableOutput, aurora.Green(fmt.Sprintf("http://%s.loophole.site", siteID)))
	fmt.Fprint(colorableOutput, " -> ")
	fmt.Fprint(colorableOutput, aurora.Green(fmt.Sprintf("%s:%d", config.Host, config.Port)))
	fmt.Println()
	fmt.Fprint(colorableOutput, "Forwarding ")
	fmt.Fprint(colorableOutput, aurora.Green(fmt.Sprintf("https://%s.loophole.site", siteID)))
	fmt.Fprint(colorableOutput, " -> ")
	fmt.Fprint(colorableOutput, aurora.Green(fmt.Sprintf("%s:%d", config.Host, config.Port)))
	fmt.Println()
	fmt.Println()
	fmt.Fprint(colorableOutput, emoji.Sprint(fmt.Sprintf(":nerd: %s", aurora.Italic("TLS Certificate will be obtained with first request to the above address, therefore first execution may be slower\n"))))
	fmt.Fprint(colorableOutput, emoji.Sprint(":newspaper: Logs:\n"))

	log.Info().Msg("Awaiting connections...")
	return listenerHTTPSOverSSH, proxiedEndpointHTTPS
}

// Start starts the tunnel on specified host and port
func Start(config lm.Config) {
	printWelcomeMessage()

	var publicKeyAuthMethod *ssh.AuthMethod = new(ssh.AuthMethod)
	var publicKey *ssh.PublicKey = new(ssh.PublicKey)

	listenerHTTPSOverSSH, proxiedEndpointHTTPS := generateListener(config, publicKeyAuthMethod, publicKey)
	defer listenerHTTPSOverSSH.Close()

	for {
		client, err := listenerHTTPSOverSSH.Accept()
		if err == io.EOF {
			log.Info().Err(err).Msg("Connection dropped, reconnecting...")
			listenerHTTPSOverSSH.Close()
			listenerHTTPSOverSSH, _ = generateListener(config, publicKeyAuthMethod, publicKey)
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
