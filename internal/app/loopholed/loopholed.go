package loopholed

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/beevik/guid"
	"github.com/loophole/cli/internal/app/loophole"
	"github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/loophole/cli/internal/pkg/communication"
)

const (

	// name of the service
	name        = "loophole"
	description = "Loophole daemon"

	// port which daemon should be listen
	port = ":9977"
)

// dependencies that are NOT required by the service, but might be used
var dependencies = []string{ /*"dummy.service"*/ }

var tunnels = make(map[string]string)

// LoopholeService implements the daemon.Executable interface
// and represents the actual service behavior
type LoopholeService struct {
	listen chan net.Conn
}

// Start gets the service up
func (svc *LoopholeService) Start() {
	// Set up listener for defined host and port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		errlog.Println("Possibly was a problem with the port binding", err)
		return
	}

	// set up channel on which to send accepted connections
	svc.listen = make(chan net.Conn, 100)
	go acceptConnection(listener, svc.listen)

	// loop work cycle with accept connections or interrupt
	// by system signal
	go func() {
		for {
			select {
			case conn, ok := <-svc.listen:
				if !ok {
					stdlog.Println("Closing connections")
					listener.Close()
					return
				}
				go handleClient(conn)
			}
		}
	}()
}

// Accept a client connection and collect it in a channel
func acceptConnection(listener net.Listener, listen chan<- net.Conn) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		listen <- conn
	}
}

func handleClient(client net.Conn) {
	reader := bufio.NewReader(client)
	for {
		actionType, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			fmt.Println("Error while reading user order")
			client.Close()
			return
		}
		communication.Info(fmt.Sprintf("Received order %s", actionType))
		sshDir := cache.GetLocalStorageDir(".ssh")

		actions := strings.Split(actionType, ",")
		switch strings.TrimSpace(actions[0]) {
		// assume HTTP PORT LOCAL_HOST DOMAIN HTTPS USERNAME PASSWORD
		case "HTTP":
			if len(actions) < 2 {
				communication.Error("Not enough arguments to start anything")
			}
			port, err := strconv.Atoi(strings.TrimSpace(actions[1]))
			localHost := "127.0.0.1"
			if actions[2] != "" {
				localHost = strings.TrimSpace(actions[2])
			}
			siteID := ""
			if len(actions) >= 4 {
				siteID = strings.TrimSpace(actions[3])
			}
			if err != nil {
				communication.Error("Invalid port")
			}
			exposeHTTPConfig := models.ExposeHTTPConfig{
				Remote: models.RemoteEndpointSpecs{
					IdentityFile: fmt.Sprintf("%s/id_rsa", sshDir),
					TunnelID:     guid.New().String(),
					SiteID:       siteID,
				},
				Local: models.LocalHTTPEndpointSpecs{
					Host: localHost,
					Port: int32(port),
				},
			}

			communication.Info(fmt.Sprintf("Starting tunnel with details: type: HTTP, localPort: %d, localHost: %s, remoteDomain: %s", port, localHost, siteID))

			authMethod, err := loophole.RegisterTunnel(&exposeHTTPConfig.Remote)
			if err != nil {
				communication.TunnelStartFailure(exposeHTTPConfig.Remote.TunnelID, err)
				return
			}
			tunnelQuitChannel := make(chan bool)
			tunnels[exposeHTTPConfig.Remote.TunnelID] = exposeHTTPConfig.Remote.SiteID
			go loophole.ForwardPort(exposeHTTPConfig, authMethod, tunnelQuitChannel)

			client.Write([]byte("Tunnel started\n"))
			client.Close()
		case "PS":
			communication.Info("Listing tunnels...")
			if len(tunnels) == 0 {
				client.Write([]byte("There are no running tunnels"))
			} else {
				client.Write([]byte("TunnelID\tSiteID"))
				for tID, sID := range tunnels {
					client.Write([]byte(fmt.Sprintf("\n%s\t%s", tID, sID)))
				}
			}
			client.Write([]byte("\n"))
			client.Close()
			return
		default:
			communication.Warn("Unknown message")
			client.Write([]byte("Unknown message"))
			client.Close()
			return
		}
	}
}

// Stop shuts down the service
func (svc *LoopholeService) Stop() {
	close(svc.listen)
}

// Run is invoked when the service is run in interective mode
// (ie during development). On Windows it is never invoked
func (svc *LoopholeService) Run() {
	svc.Start()
	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	fmt.Println("Daemon awaiting orders...")
	// loop work cycle with accept connections or interrupt
	// by system signal
loop:
	for {
		select {
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			if killSignal == os.Interrupt {
				stdlog.Println("Daemon was interrupted by system signal")
			}
			stdlog.Println("Daemon was killed")
			break loop
		}
	}

	svc.Stop()
}
