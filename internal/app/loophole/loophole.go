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
	"path"

	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/mitchellh/go-homedir"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/crypto/acme/autocert"
	"golang.org/x/crypto/ssh"
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

func getLocalStorageDir() string {
	home, err := homedir.Dir()
	if err != nil {
		logger.Fatal("Error reading user home directory ", zap.Error(err))
	}

	return path.Join(home, ".local", "loophole")
}

func getPublicKey(file string) (ssh.AuthMethod, ssh.PublicKey, error) {
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
	data := map[string]string{
		"key": publicKeyString,
	}
	if siteID != "" {
		data["id"] = siteID
	}

	jsonData, _ := json.Marshal(data)
	resp, err := http.Post(fmt.Sprintf("%s/site", apiURL), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)

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
	publicKeyAuthMethod, publicKey, err := getPublicKey(config.IdentityFile)
	if err != nil {
		logger.Fatal("No public key available", zap.Error(err))
	}

	siteID, err := registerSite(config.APIURL, publicKey, config.SiteID)
	if err != nil {
		logger.Fatal("Failed to register site", zap.Error(err))
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
		Cache:      autocert.DirCache(getLocalStorageDir()),                                              //Folder for storing certificates
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
	go server.ServeTLS(proxyListenerHTTPS, "certs/server.crt", "certs/server.key")
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
