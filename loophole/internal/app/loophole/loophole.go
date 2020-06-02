package loophole

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"

	"golang.org/x/crypto/ssh"
)

func getPublicKey(file string) (ssh.AuthMethod, error) {
	key, err := ioutil.ReadFile(file)

	if err != nil {
		return nil, err
	}

	signer, err := ssh.ParsePrivateKey(key)

	if err != nil {

		return nil, err
	}

	return ssh.PublicKeys(signer), nil
}

// Endpoint is representing host address
type Endpoint struct {
	Host string
	Port int
}

func (endpoint *Endpoint) String() string {
	if endpoint != nil {
		return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
	}
	return ""
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
			log.Printf("error while copy remote->local: %s", err)
		}
		chDone <- true
	}()

	// Start local -> remote data transfer
	go func() {
		_, err := io.Copy(remote, client)
		if err != nil {
			log.Printf("error while copy local->remote: %s", err)
		}
		chDone <- true
	}()

	<-chDone
}

// remote SSH server
var serverEndpoint = Endpoint{
	Host: "gateway.loophole.cloud",
	Port: 8022,
}

// remote forwarding port (on remote SSH server network)
var remoteEndpoint = Endpoint{
	Host: "127.0.0.1",
	Port: 80,
}

// Start starts the tunnel on specified host and port
func Start(port int, host string, secure bool, identityFile string) {
	localEndpoint := Endpoint{
		Host: host,
		Port: port,
	}
	publicKey, err := getPublicKey(identityFile)
	if err != nil {
		log.Fatalf("no public key available: %s", err)
	}

	sshConfig := &ssh.ClientConfig{
		User: "whatever_http",
		Auth: []ssh.AuthMethod{
			publicKey,
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	// Connect to SSH remote server using serverEndpoint
	serverConn, err := ssh.Dial("tcp", serverEndpoint.String(), sshConfig)

	if err != nil {
		log.Fatalf("Dial INTO remote server error: %s", err)
	}

	// Listen on remote server port
	listener, err := serverConn.Listen("tcp", remoteEndpoint.String())
	if err != nil {
		log.Fatalf("Listen open port ON remote server error: %s", err)
	}
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
