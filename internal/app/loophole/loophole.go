package loophole

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/loophole/cli/internal/pkg/token"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/ssh"
)

const (
	apiURL = "https://api.loophole.cloud"
)

var logger *zap.Logger

func init() {
	atomicLevel := zap.NewAtomicLevel()
	encoderCfg := zap.NewProductionEncoderConfig()
	logger = zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atomicLevel,
	))

	atomicLevel.SetLevel(zap.DebugLevel)
}

func parsePublicKey(file string) (ssh.AuthMethod, ssh.PublicKey, error) {
	key, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)

	if err != nil {

		return nil, nil, err
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
			logger.Debug("Error copying remote->local:", zap.Error(err))
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			logger.Debug("Error copying local->remote:", zap.Error(err))
		}
		chDone <- true
	}()

	<-chDone
}

func registerSite(apiURL string, publicKey ssh.PublicKey, siteID string) (string, error) {
	publicKeyString := publicKey.Type() + " " + base64.StdEncoding.EncodeToString(publicKey.Marshal())

	if !token.IsTokenSaved() {
		logger.Fatal("Please log in before using Loophole")
	}

	accessToken, err := token.GetAccessToken()
	if err != nil {
		logger.Fatal("There was a problem reading token", zap.Error(err))
	}

	data := map[string]string{
		"key": publicKeyString,
	}
	if siteID != "" {
		data["id"] = siteID
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		logger.Fatal("There was a problem encoding request body", zap.Error(err))
	}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/site", apiURL), bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Fatal("There was a problem creating request body", zap.Error(err))
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

	logger.Debug("Response", zap.Any("result", result))
	defer resp.Body.Close()

	siteID, ok := result["id"].(string)
	if !ok {
		logger.Fatal("Error converting siteId to string")
	}
	return siteID, nil
}

// remote forwarding port (on remote SSH server network)
var remoteEndpoint = lm.Endpoint{
	Host: "127.0.0.1",
	Port: 80,
}

// Start starts the tunnel on specified host and port
func Start(config lm.Config) {
	defer logger.Sync()

	logger.Debug("Starting the tunnels..")

	localEndpoint := lm.Endpoint{
		Host: config.Host,
		Port: config.Port,
	}
	publicKeyAuthMethod, publicKey, err := parsePublicKey(config.IdentityFile)
	if err != nil {
		logger.Fatal("No public key available", zap.Error(err))
	}

	siteID, err := registerSite(apiURL, publicKey, config.SiteID)
	if err != nil {
		logger.Debug("Failed to register site", zap.Error(err))
		logger.Debug("Trying to refresh token")
		err := token.RefreshToken()
		if err != nil {
			logger.Fatal("Failed to refresh token", zap.Error(err))
		}
		siteID, err = registerSite(apiURL, publicKey, config.SiteID)
		if err != nil {
			logger.Fatal("Failed to register site", zap.Error(err))
		}
	}

	sshConfigHTTPS := &ssh.ClientConfig{
		User: fmt.Sprintf("%s_https", siteID),
		Auth: []ssh.AuthMethod{
			publicKeyAuthMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to SSH remote server using GatewayEndpoint
	serverSSHConnHTTPS, err := ssh.Dial("tcp", config.GatewayEndpoint.String(), sshConfigHTTPS)

	if err != nil {
		logger.Fatal("Dialing SSH Gateway for HTTPS failed", zap.Error(err))
	}
	logger.Debug("Dialing SSH Gateway for HTTPS succeded")

	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(fmt.Sprintf("%s.loophole.site", siteID), "abc.loophole.site"), //Your domain here
		Cache:      autocert.DirCache(cache.GetLocalStorageDir("certs")),                                 //Folder for storing certificates
		Email:      fmt.Sprintf("%s@loophole.main.dev", siteID),
	}
	logger.Debug("Cert Manager created")

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "http",
		Host:   localEndpoint.String(),
	})
	logger.Debug("Proxy via http created", zap.String("target", localEndpoint.String()))

	server := &http.Server{
		Handler:   proxy,
		TLSConfig: certManager.TLSConfig(),
	}
	logger.Debug("Server for proxy created")
	proxyListenerHTTPS, err := net.Listen("tcp", ":0")
	if err != nil {
		logger.Fatal("Failed to listen on TLS proxy for HTTPS", zap.Error(err))
	}
	logger.Debug("Proxy server for HTTPS listening", zap.Int("port", proxyListenerHTTPS.Addr().(*net.TCPAddr).Port))
	go server.ServeTLS(proxyListenerHTTPS, "", "")
	logger.Debug("Started servers go routines")

	listenerHTTPSOverSSH, err := serverSSHConnHTTPS.Listen("tcp", remoteEndpoint.String())
	if err != nil {
		logger.Fatal("Listening on remote endpoint for HTTPS failed", zap.Error(err))
	}
	logger.Debug("Listening on remote endpoint for HTTPS succeded")

	defer listenerHTTPSOverSSH.Close()

	proxiedEndpointHTTPS := &lm.Endpoint{
		Host: "127.0.0.1",
		Port: int32(proxyListenerHTTPS.Addr().(*net.TCPAddr).Port),
	}

	logger.Debug("Printing user friendly info about startup")
	fmt.Println("Loophole")
	fmt.Println()
	fmt.Printf("Forwarding http://%s.loophole.site -> %s:%d\n", siteID, config.Host, config.Port)
	fmt.Printf("Forwarding https://%s.loophole.site -> %s:%d\n", siteID, config.Host, config.Port)
	logger.Debug("Printed user friendly info about startup")
	for {
		client, err := listenerHTTPSOverSSH.Accept()
		if err != nil {
			logger.Debug("Failed to accept connection for HTTPS", zap.Error(err))
			continue
		}
		logger.Debug("Succeded to accept connection for HTTPS")
		// Open a (local) connection to proxiedEndpointHTTPS whose content will be forwarded to serverEndpoint
		local, err := net.Dial("tcp", proxiedEndpointHTTPS.String())
		if err != nil {
			logger.Fatal("Dialing into local proxy for HTTPS failed", zap.Error(err))
		}
		logger.Debug("Dialing into local proxy for HTTPS succeded")
		go handleClient(client, local)
	}
}
