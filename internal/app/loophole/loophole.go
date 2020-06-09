package loophole

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	lm "github.com/loophole/cli/internal/app/loophole/models"
	"golang.org/x/crypto/ssh"
)

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
// Handle local client connections and tunnel data to the remote server
// Will use io.Copy - http://golang.org/pkg/io/#Copy
func handleClient(client net.Conn, remote net.Conn) {
	defer client.Close()
	chDone := make(chan bool)

	// Start remote -> local data transfer
	go func() {
		_, err := io.Copy(client, remote)
		if err != nil {
			log.Printf("Error copying data: remote->local: %s", err)
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			log.Printf("Error copying data: local->remote: %s", err)
		}
		chDone <- true
	}()

	<-chDone
}

func getSiteID(apiURL string, publicKey ssh.PublicKey) (string, error) {
	publicKeyString := publicKey.Type() + " " + base64.StdEncoding.EncodeToString(publicKey.Marshal())
	data := map[string]string{
		"key": publicKeyString,
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
		log.Fatalf("Error converting siteId to string")
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
	localEndpoint := lm.Endpoint{
		Host: config.Host,
		Port: config.Port,
	}
	publicKeyAuthMethod, publicKey, err := getPublicKey(config.IdentityFile)
	if err != nil {
		log.Fatalf("No public key available: %s", err)
	}

	siteID, err := getSiteID(config.APIURL, publicKey)
	if err != nil {
		log.Fatalf("Error getting site id from API: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: siteID,
		Auth: []ssh.AuthMethod{
			publicKeyAuthMethod,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to SSH remote server using GatewayEndpoint
	serverConn, err := ssh.Dial("tcp", config.GatewayEndpoint.String(), sshConfig)

	if err != nil {
		log.Fatalf("Dial INTO remote server error: %s", err)
	}

	// Listen on remote server port
	listener, err := serverConn.Listen("tcp", remoteEndpoint.String())
	if err != nil {
		log.Fatalf("Listen open port ON remote server error: %s", err)
	}

	fmt.Println("Loophole")
	fmt.Println()
	fmt.Printf("Forwarding http://%s.loophole.site -> %s:%d\n", siteID, config.Host, config.Port)
	fmt.Printf("Forwarding https://%s.loophole.site -> %s:%d\n", siteID, config.Host, config.Port)
	defer listener.Close()

	// handle incoming connections on reverse forwarded tunnel
	for {
		// Open a (local) connection to localEndpoint whose content will be forwarded so serverEndpoint
		local, err := net.Dial("tcp", localEndpoint.String())
		if err != nil {
			log.Fatalf("Dial INTO local service error: %s", err)
		}

		client, err := listener.Accept()
		if err != nil {
			log.Fatalf("Failed to accept connection: %v", err)
		}

		handleClient(client, local)
	}
}
