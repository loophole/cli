// +build desktop

package ui

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	"github.com/markbates/pkger"
	"github.com/rs/zerolog/log"
	"github.com/zserge/lorca"

	"github.com/gorilla/websocket"

	"github.com/loophole/cli/config"
	"github.com/loophole/cli/internal/app/loophole"
	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/loophole/cli/internal/pkg/token"
)

var upgrader = websocket.Upgrader{} // use default options
var tunnelQuitChannel = make(chan bool)
var authQuitChannel = make(chan bool)
var authAlreadyRan = false

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	communication.SetCommunicationMechanism(communication.NewWebsocketLogger(c, websocket.TextMessage))
	loggedIn := token.IsTokenSaved()
	communication.PrintWelcomeMessage(loggedIn)
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Warn().Err(err).Msg("read:")
			break
		}
		var decodedMessage Message
		err = json.Unmarshal(message, &decodedMessage)
		if err != nil {
			err = c.WriteMessage(mt, []byte(fmt.Sprintf("Error decoding message: %s", err.Error())))
			if err != nil {
				log.Warn().Err(err).Msg("write:")
				break
			}
		}

		switch decodedMessage.MessageType {
		case MessageTypeStartTunnel:
			switch decodedMessage.StartTunnelMessage.TunnelType {
			case loophole.HTTP:
				go func() {
					startHTTPTunnel(decodedMessage.StartTunnelMessage.ExposeHTTPConfig, tunnelQuitChannel)
				}()
			case loophole.Directory:
				go func() {
					startDirectoryTunnel(decodedMessage.StartTunnelMessage.ExposeDirectoryConfig, tunnelQuitChannel)
				}()
			case loophole.WebDav:
				go func() {
					startWebdavTunnel(decodedMessage.StartTunnelMessage.ExposeWebdavConfig, tunnelQuitChannel)
				}()
			default:
				err = c.WriteMessage(mt, []byte(fmt.Sprintf("Unrecognized tunnel type: %s", decodedMessage.StartTunnelMessage.TunnelType)))
				if err != nil {
					log.Warn().Err(err).Msg("write:")
					break
				}
			}
		case MessageTypeStopTunnel:
			tunnelQuitChannel <- true
		case MessageTypeAuthorization:
			if authAlreadyRan {
				authQuitChannel <- true
			} else {
				authAlreadyRan = true
			}
			go func() {
				startAuthorizationProcess(authQuitChannel)
			}()
		case MessageTypeLogout:
			authAlreadyRan = false
			go func() {
				logout()
			}()
		case MessageTypeOpenBrowser:
			go func() {
				openInBrowser(decodedMessage.OpenInBrowserMessage.URL)
			}()
		default:
			err = c.WriteMessage(mt, []byte(fmt.Sprintf("Unrecognized message type: %s", decodedMessage.MessageType)))
			if err != nil {
				log.Info().Err(err).Msg("write:")
				break
			}
		}
	}
}

// Display shows the main app window
func Display(appConfig config.ApplicationConfig) {
	// path is absolute with / set to module (repository) root
	box := pkger.Dir("/ui/desktop/build")

	http.Handle("/", http.FileServer(box))
	http.HandleFunc("/ws", websocketHandler)

	localListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	go http.Serve(localListener, nil)

	ui, _ := lorca.New(fmt.Sprintf("http://%s", localListener.Addr().String()), "", 1280, 960)
	defer ui.Close()

	<-ui.Done()
}
