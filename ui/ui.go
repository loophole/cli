//go:build desktop
// +build desktop

package ui

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net"
	"net/http"

	"github.com/ncruces/zenity"
	"github.com/rs/zerolog/log"
	"github.com/skratchdot/open-golang/open"
	"github.com/zserge/lorca"

	"github.com/gorilla/websocket"

	"github.com/loophole/cli/internal/app/loophole"
	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/loophole/cli/internal/pkg/token"
)

var upgrader = websocket.Upgrader{} // use default options
var authQuitChannel = make(chan bool)
var authAlreadyRan = false

var tunnelQuitChannels = make(map[string](chan bool))
var siteToRequestMapping = make(map[string]string)

func websocketHandler(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	communication.SetCommunicationMechanism(communication.NewWebsocketLogger(c))
	loggedIn := token.IsTokenSaved()
	idToken := token.GetIdToken()
	communication.ApplicationStart(loggedIn, idToken)
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Warn().Err(err).Msg("read:")
			break
		}
		var decodedMessage Message
		err = json.Unmarshal(message, &decodedMessage)
		if err != nil {
			communication.Warn("Error decoding message")
			communication.Warn(err.Error())
		}

		switch decodedMessage.Type {
		case MessageTypeStartTunnelHTTP:
			var exposeHTTPConfig lm.ExposeHTTPConfig
			err = json.Unmarshal(decodedMessage.Payload, &exposeHTTPConfig)
			if err != nil {
				communication.Warn("Error decoding message")
				communication.Warn(err.Error())
			}
			tunnelQuitChannel := make(chan bool)
			go func() {
				sshDir := cache.GetLocalStorageDir(".ssh") //getting our sshDir and creating it, if it doesn't exist
				exposeHTTPConfig.Remote.IdentityFile = fmt.Sprintf("%s/id_rsa", sshDir)

				communication.TunnelDebug(exposeHTTPConfig.Remote.TunnelID, fmt.Sprintf("Got request for SiteID: '%s'", exposeHTTPConfig.Remote.SiteID))
				if _, ok := siteToRequestMapping[exposeHTTPConfig.Remote.SiteID]; exposeHTTPConfig.Remote.SiteID != "" && ok {
					communication.TunnelStartFailure(exposeHTTPConfig.Remote.TunnelID,
						fmt.Errorf("Tunnel '%s' is already running", exposeHTTPConfig.Remote.SiteID))
					return
				}
				authMethod, err := loophole.RegisterTunnel(&exposeHTTPConfig.Remote)
				if err != nil {
					communication.TunnelStartFailure(exposeHTTPConfig.Remote.TunnelID, err)
					return
				}
				communication.TunnelDebug(exposeHTTPConfig.Remote.TunnelID, fmt.Sprintf("Obtained SiteID: '%s'", exposeHTTPConfig.Remote.SiteID))
				if _, ok := siteToRequestMapping[exposeHTTPConfig.Remote.SiteID]; ok {
					communication.TunnelStartFailure(exposeHTTPConfig.Remote.TunnelID,
						fmt.Errorf("Tunnel '%s' is already running", exposeHTTPConfig.Remote.SiteID))
					return
				}
				tunnelQuitChannels[exposeHTTPConfig.Remote.TunnelID] = tunnelQuitChannel
				siteToRequestMapping[exposeHTTPConfig.Remote.SiteID] = exposeHTTPConfig.Remote.TunnelID
				loophole.ForwardPort(exposeHTTPConfig, authMethod, tunnelQuitChannel)
			}()
		case MessageTypeStartTunnelDirectory:
			var exposeDirectoryConfig lm.ExposeDirectoryConfig
			err = json.Unmarshal(decodedMessage.Payload, &exposeDirectoryConfig)
			if err != nil {
				communication.Warn("Error decoding message")
				communication.Warn(err.Error())
			}

			tunnelQuitChannel := make(chan bool)
			go func() {
				sshDir := cache.GetLocalStorageDir(".ssh") //getting our sshDir and creating it, if it doesn't exist
				exposeDirectoryConfig.Remote.IdentityFile = fmt.Sprintf("%s/id_rsa", sshDir)

				communication.TunnelDebug(exposeDirectoryConfig.Remote.TunnelID, fmt.Sprintf("Got request for SiteID: '%s'", exposeDirectoryConfig.Remote.SiteID))
				if _, ok := siteToRequestMapping[exposeDirectoryConfig.Remote.SiteID]; exposeDirectoryConfig.Remote.SiteID != "" && ok {
					communication.TunnelStartFailure(exposeDirectoryConfig.Remote.TunnelID,
						fmt.Errorf("Tunnel '%s' is already running", exposeDirectoryConfig.Remote.SiteID))

					return
				}

				authMethod, err := loophole.RegisterTunnel(&exposeDirectoryConfig.Remote)
				if err != nil {
					communication.TunnelStartFailure(exposeDirectoryConfig.Remote.TunnelID, err)
					return

				}
				communication.TunnelDebug(exposeDirectoryConfig.Remote.TunnelID, fmt.Sprintf("Obtained SiteID: '%s'", exposeDirectoryConfig.Remote.SiteID))
				if _, ok := siteToRequestMapping[exposeDirectoryConfig.Remote.SiteID]; ok {
					communication.TunnelStartFailure(exposeDirectoryConfig.Remote.TunnelID,
						fmt.Errorf("Tunnel '%s' is already running", exposeDirectoryConfig.Remote.SiteID))
					return
				}
				tunnelQuitChannels[exposeDirectoryConfig.Remote.TunnelID] = tunnelQuitChannel
				siteToRequestMapping[exposeDirectoryConfig.Remote.SiteID] = exposeDirectoryConfig.Remote.TunnelID
				loophole.ForwardDirectory(exposeDirectoryConfig, authMethod, tunnelQuitChannel)
			}()
		case MessageTypeStartTunnelWebDav:
			var exposeWebdavConfig lm.ExposeWebdavConfig
			err = json.Unmarshal(decodedMessage.Payload, &exposeWebdavConfig)
			if err != nil {
				communication.Warn("Error decoding message")
				communication.Warn(err.Error())
			}

			tunnelQuitChannel := make(chan bool)
			go func() {
				sshDir := cache.GetLocalStorageDir(".ssh") //getting our sshDir and creating it, if it doesn't exist
				exposeWebdavConfig.Remote.IdentityFile = fmt.Sprintf("%s/id_rsa", sshDir)

				communication.TunnelDebug(exposeWebdavConfig.Remote.TunnelID, fmt.Sprintf("Got request for SiteID: '%s'", exposeWebdavConfig.Remote.SiteID))
				if _, ok := siteToRequestMapping[exposeWebdavConfig.Remote.SiteID]; exposeWebdavConfig.Remote.SiteID != "" && ok {
					communication.TunnelStartFailure(exposeWebdavConfig.Remote.TunnelID,
						fmt.Errorf("Tunnel '%s' is already running", exposeWebdavConfig.Remote.SiteID))
					return
				}

				authMethod, err := loophole.RegisterTunnel(&exposeWebdavConfig.Remote)
				if err != nil {
					communication.TunnelStartFailure(exposeWebdavConfig.Remote.TunnelID, err)
					return
				}

				communication.TunnelDebug(exposeWebdavConfig.Remote.TunnelID, fmt.Sprintf("Obtained SiteID: '%s'", exposeWebdavConfig.Remote.SiteID))
				if _, ok := siteToRequestMapping[exposeWebdavConfig.Remote.SiteID]; ok {
					communication.TunnelStartFailure(exposeWebdavConfig.Remote.TunnelID,
						fmt.Errorf("Tunnel '%s' is already running", exposeWebdavConfig.Remote.SiteID))
					return
				}

				tunnelQuitChannels[exposeWebdavConfig.Remote.TunnelID] = tunnelQuitChannel
				siteToRequestMapping[exposeWebdavConfig.Remote.SiteID] = exposeWebdavConfig.Remote.TunnelID
				loophole.ForwardDirectoryViaWebdav(exposeWebdavConfig, authMethod, tunnelQuitChannel)
			}()
		case MessageTypeStopTunnel:
			var stopTunnelMessage StopTunnelMessage
			err = json.Unmarshal(decodedMessage.Payload, &stopTunnelMessage)
			if err != nil {
				communication.Warn("Error decoding message")
				communication.Warn(err.Error())
			}
			tunnelQuitChannels[stopTunnelMessage.TunnelID] <- true
			siteID, ok := findKeyByValue(siteToRequestMapping, stopTunnelMessage.TunnelID)
			if ok {
				communication.Debug(fmt.Sprintf("Removing %s from the dict", siteID))
				delete(siteToRequestMapping, siteID)

			}
		case MessageTypeAuthorization:
			if authAlreadyRan {
				authQuitChannel <- true
			} else {
				authAlreadyRan = true
			}
			go func() {
				deviceCodeSpec, err := token.RegisterDevice()
				if err != nil {
					communication.LoginFailure(fmt.Errorf("Error obtaining device code: %s", err.Error()))
				}
				communication.LoginStart(*deviceCodeSpec)
				tokens, err := token.PollForToken(deviceCodeSpec.DeviceCode, deviceCodeSpec.Interval, authQuitChannel)
				if err != nil {
					communication.LoginFailure(fmt.Errorf("Error obtaining token: %s", err.Error()))
					return
				}
				err = token.SaveToken(tokens)
				if err != nil {
					communication.LoginFailure(fmt.Errorf("Error saving token: %s", err.Error()))
					return
				}
				communication.LoginSuccess(tokens.IDToken)
				communication.Info("Logged in successfully")
			}()
		case MessageTypeLogout:
			authAlreadyRan = false
			go func() {
				err := token.DeleteTokens()
				if err != nil {
					communication.LogoutFailure(fmt.Errorf("Error logging out: %s", err.Error()))
					return
				}
				communication.LogoutSuccess()
				communication.Info("Logged out successfully")
			}()
		case MessageTypeOpenBrowser:
			var openInBrowserMessage OpenInBrowserMessage
			err = json.Unmarshal(decodedMessage.Payload, &openInBrowserMessage)
			if err != nil {
				communication.Warn("Error decoding message")
				communication.Warn(err.Error())
			}
			go func() {
				err := open.Run(openInBrowserMessage.URL)
				if err != nil {
					communication.Warn("Error opening url")
					communication.Warn(err.Error())
				}
			}()
		default:
			communication.Warn(fmt.Sprintf("Unrecognized message type", decodedMessage.Type))
		}
	}
}

//go:embed desktop/build/*
var box embed.FS

// Display shows the main app window
func Display() {
	chromeLocation := lorca.LocateChrome()
	if chromeLocation == "" {
		message := "Chrome/Chromium >= 70 is required."
		zenity.Error(message)
		communication.Fatal(message)
	}

	subFS, err := fs.Sub(box, "desktop/build")
	if err != nil {
		panic(err)
	}

	http.Handle("/", http.FileServer(http.FS(subFS)))
	http.HandleFunc("/ws", websocketHandler)

	localListener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		panic(err)
	}
	go http.Serve(localListener, nil)

	ui, err := lorca.New(fmt.Sprintf("http://%s", localListener.Addr().String()), "", 980, 800)
	if err != nil {
		communication.Fatal(fmt.Sprintf("Unable to run Chrome/Chromium: %s", err.Error()))
	}
	defer ui.Close()

	<-ui.Done()
}

func findKeyByValue(m map[string]string, v string) (string, bool) {
	for k, x := range m {
		if x == v {
			return k, true
		}
	}
	return "", false
}
