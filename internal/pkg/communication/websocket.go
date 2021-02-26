// +build desktop

package communication

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/loophole/cli/config"
	coreModels "github.com/loophole/cli/internal/app/loophole/models"
	authModels "github.com/loophole/cli/internal/pkg/token/models"
	"github.com/loophole/cli/internal/pkg/urlmaker"
	"github.com/mitchellh/go-homedir"
	"github.com/ncruces/zenity"
)

// MessageType is websocket message type
type MessageType string

const (
	// MessageTypeLog represents application level log messages
	MessageTypeLog MessageType = "MT_Log"
	// MessageTypeTunnelLog represents tunnel level log messages
	MessageTypeTunnelLog MessageType = "MT_TunnelLog"

	MessageTypeAppStart            MessageType = "MT_ApplicationStart"
	MessageTypeAppStop             MessageType = "MT_ApplicationStop"
	MessageTypeNewVersionAvailable MessageType = "MT_ApplicationNewVersionAvailable"

	MessageTypeLogin        MessageType = "MT_Login"
	MessageTypeLoginSuccess MessageType = "MT_LoginSuccess"
	MessageTypeLoginFailure MessageType = "MT_LoginFailure"

	MessageTypeLogout        MessageType = "MT_Logout"
	MessageTypeLogoutSuccess MessageType = "MT_LogoutSuccess"
	MessageTypeLogoutFailure MessageType = "MT_LogoutFailure"

	MessageTypeTunnelStart        MessageType = "MT_TunnelStart"
	MessageTypeTunnelStartSuccess MessageType = "MT_TunnelStartSuccess"
	MessageTypeTunnelStartFailure MessageType = "MT_TunnelStartFailure"

	MessageTypeTunnelStop MessageType = "MT_TunnelStop"

	MessageTypeLoadingStart   MessageType = "MT_LoadingStart"
	MessageTypeLoadingSuccess MessageType = "MT_LoadingSuccess"
	MessageTypeLoadingFailure MessageType = "MT_LoadingFailure"
)

// NewWebsocketLogger is websocket mechanism constructor
func NewWebsocketLogger(wsClient *websocket.Conn) Mechanism {
	logger := websocketLogger{
		wsClient: wsClient,
	}

	return &logger
}

type websocketLogger struct {
	wsClient     *websocket.Conn
	messageMutex sync.Mutex
}

type logMessage struct {
	Type    MessageType `json:"type"`
	Message string      `json:"message"`
	Class   string      `json:"class"`
}

type tunnelLogMessage struct {
	Type     MessageType `json:"type"`
	Message  string      `json:"message"`
	Class    string      `json:"class"`
	TunnelID string      `json:"tunnelId"`
}

type appStartMessage struct {
	Type          MessageType          `json:"type"`
	LoggedIn      bool                 `json:"loggedIn"`
	DisplayConfig config.DisplayConfig `json:"displayConfig"`
	HomeDirectory string               `json:"homeDirectory"`
	IDToken       string               `json:"idToken"`
	Version       string               `json:"version"`
	CommitHash    string               `json:"commitHash"`
}

type appStopMessage struct {
	Type MessageType `json:"type"`
}

type appNewVersionMessage struct {
	Type    MessageType `json:"type"`
	Version string      `json:"version"`
}

type loginMessage struct {
	Type                    MessageType `json:"type"`
	DeviceCode              string      `json:"deviceCode"`
	UserCode                string      `json:"userCode"`
	VerificationURI         string      `json:"verificationUri"`
	VerificationURIComplete string      `json:"verificationUriComplete"`
}

type loginSuccessMessage struct {
	Type    MessageType `json:"type"`
	IDToken string      `json:"idToken"`
}

type loginFailureMessage struct {
	Type  MessageType `json:"type"`
	Error string      `json:errror"`
}

type logoutSuccessMessage struct {
	Type MessageType `json:"type"`
}

type logoutFailureMessage struct {
	Type  MessageType `json:"type"`
	Error string      `json:errror"`
}

type tunnelStartMessage struct {
	Type     MessageType `json:"type"`
	TunnelID string      `json:"tunnelId"`
}

type tunnelStartSuccessMessage struct {
	Type      MessageType `json:"type"`
	TunnelID  string      `json:"tunnelId"`
	SiteID    string      `json:"siteId"`
	Domain    string      `json:"domain"`
	SiteAddrs []string    `json:"siteAddrs"`
	LocalAddr string      `json:"localAddr"`
}
type tunnelStartFailureMessage struct {
	Type     MessageType `json:"type"`
	TunnelID string      `json:"tunnelId"`
	Error    string      `json:"error"`
}

type tunnelStopSuccessMessage struct {
	Type     MessageType `json:"type"`
	TunnelID string      `json:"tunnelId"`
}

type loadingStartMessage struct {
	Type     MessageType `json:"type"`
	TunnelID string      `json:"tunnelId"`
	Message  string      `json:"message"`
}
type loadingSuccessMessage struct {
	Type     MessageType `json:"type"`
	TunnelID string      `json:"tunnelId"`
}
type loadingFailureMessage struct {
	Type     MessageType `json:"type"`
	TunnelID string      `json:"tunnelId"`
	Error    string      `json:"error"`
}

func (l *websocketLogger) write(websocketMessage interface{}) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()

	err := l.wsClient.WriteJSON(websocketMessage)
	if err != nil {
		defaultLogger.Fatal("Communication over websocket is failing")
	}
}

func (l *websocketLogger) TunnelDebug(tunnelID string, message string) {
	if !config.Config.Display.Verbose {
		return
	}
	l.write(tunnelLogMessage{
		Type:     MessageTypeTunnelLog,
		Message:  message,
		Class:    "grey",
		TunnelID: tunnelID,
	})
}
func (l *websocketLogger) TunnelInfo(tunnelID string, message string) {
	l.write(tunnelLogMessage{
		Type:     MessageTypeTunnelLog,
		Message:  message,
		Class:    "info",
		TunnelID: tunnelID,
	})
}
func (l *websocketLogger) TunnelWarn(tunnelID string, message string) {
	l.write(tunnelLogMessage{
		Type:     MessageTypeTunnelLog,
		Message:  message,
		Class:    "warning",
		TunnelID: tunnelID,
	})
}
func (l *websocketLogger) TunnelError(tunnelID string, message string) {
	l.write(tunnelLogMessage{
		Type:     MessageTypeTunnelLog,
		Message:  message,
		Class:    "danger",
		TunnelID: tunnelID,
	})
}

func (l *websocketLogger) Debug(message string) {
	if !config.Config.Display.Verbose {
		return
	}
	l.write(logMessage{
		Type:    MessageTypeLog,
		Message: message,
		Class:   "grey",
	})
}
func (l *websocketLogger) Info(message string) {
	l.write(logMessage{
		Type:    MessageTypeLog,
		Message: message,
		Class:   "info",
	})
}
func (l *websocketLogger) Warn(message string) {
	l.write(logMessage{
		Type:    MessageTypeLog,
		Message: message,
		Class:   "warning",
	})
}
func (l *websocketLogger) Error(message string) {
	l.write(logMessage{
		Type:    MessageTypeLog,
		Message: message,
		Class:   "danger",
	})
}
func (l *websocketLogger) Fatal(message string) {
	l.write(logMessage{
		Type:    MessageTypeLog,
		Message: message,
		Class:   "danger",
	})

	zenity.Error(message)
	defaultLogger.Fatal(message)
}

func (l *websocketLogger) ApplicationStart(loggedIn bool, idToken string) {
	websocketMessage := appStartMessage{
		Type:          MessageTypeAppStart,
		LoggedIn:      loggedIn,
		DisplayConfig: config.Config.Display,
		IDToken:       idToken,
		Version:       config.Config.Version,
		CommitHash:    config.Config.CommitHash,
	}
	home, err := homedir.Dir()
	if err != nil {
		websocketMessage.HomeDirectory = ""
	} else {
		websocketMessage.HomeDirectory = home
	}

	l.write(websocketMessage)
}
func (l *websocketLogger) ApplicationStop() {
	l.write(appStopMessage{
		Type: MessageTypeAppStop,
	})
}

func (l *websocketLogger) TunnelStart(tunnelID string) {
	l.write(tunnelStartMessage{
		Type:     MessageTypeTunnelStart,
		TunnelID: tunnelID,
	})
}

func (l *websocketLogger) TunnelStartSuccess(remoteConfig coreModels.RemoteEndpointSpecs, localEndpoint string) {
	siteAddrs := []string{}
	siteAddrs = append(siteAddrs, urlmaker.GetSiteURL("https", remoteConfig.SiteID, remoteConfig.Domain))

	l.write(tunnelStartSuccessMessage{
		Type:      MessageTypeTunnelStartSuccess,
		TunnelID:  remoteConfig.TunnelID,
		SiteID:    remoteConfig.SiteID,
		Domain:    remoteConfig.Domain,
		SiteAddrs: siteAddrs,
		LocalAddr: localEndpoint,
	})
}
func (l *websocketLogger) TunnelStartFailure(tunnelID string, err error) {
	l.write(tunnelStartFailureMessage{
		Type:     MessageTypeTunnelStartFailure,
		TunnelID: tunnelID,
		Error:    err.Error(),
	})
}

func (l *websocketLogger) TunnelStopSuccess(tunnelID string) {
	l.write(tunnelStopSuccessMessage{
		Type:     MessageTypeTunnelStop,
		TunnelID: tunnelID,
	})
}

func (l *websocketLogger) LoginStart(deviceCodeSpec authModels.DeviceCodeSpec) {
	l.write(loginMessage{
		Type:                    MessageTypeLogin,
		DeviceCode:              deviceCodeSpec.DeviceCode,
		UserCode:                deviceCodeSpec.UserCode,
		VerificationURI:         deviceCodeSpec.VerificationURI,
		VerificationURIComplete: deviceCodeSpec.VerificationURIComplete,
	})
}
func (l *websocketLogger) LoginSuccess(idToken string) {
	l.write(loginSuccessMessage{
		Type:    MessageTypeLoginSuccess,
		IDToken: idToken,
	})
}
func (l *websocketLogger) LoginFailure(err error) {
	l.write(loginFailureMessage{
		Type:  MessageTypeLoginFailure,
		Error: err.Error(),
	})
}

func (l *websocketLogger) LogoutSuccess() {
	l.write(logoutSuccessMessage{
		Type: MessageTypeLogoutSuccess,
	})
}
func (l *websocketLogger) LogoutFailure(err error) {
	l.write(logoutFailureMessage{
		Type:  MessageTypeLogoutFailure,
		Error: err.Error(),
	})
}

func (l *websocketLogger) LoadingStart(tunnelID string, loaderMessage string) {
	l.write(loadingStartMessage{
		Type:     MessageTypeLoadingStart,
		TunnelID: tunnelID,
		Message:  loaderMessage,
	})
}
func (l *websocketLogger) LoadingStartWithRequest(tunnelID string, loaderMessage string) {
	l.write(loadingStartMessage{
		Type:     MessageTypeLoadingStart,
		TunnelID: tunnelID,
		Message:  loaderMessage,
	})
}
func (l *websocketLogger) LoadingSuccess(tunnelID string) {
	l.write(loadingSuccessMessage{
		Type:     MessageTypeLoadingSuccess,
		TunnelID: tunnelID,
	})

}
func (l *websocketLogger) LoadingSuccessWithRequest(tunnelID string) {
	l.write(loadingSuccessMessage{
		Type:     MessageTypeLoadingSuccess,
		TunnelID: tunnelID,
	})

}
func (l *websocketLogger) LoadingFailure(tunnelID string, err error) {
	l.write(loadingFailureMessage{
		Type:     MessageTypeLoadingFailure,
		TunnelID: tunnelID,
		Error:    err.Error(),
	})
}
func (l *websocketLogger) LoadingFailureWithRequest(tunnelID string, err error) {
	l.write(loadingFailureMessage{
		Type:     MessageTypeLoadingFailure,
		TunnelID: tunnelID,
		Error:    err.Error(),
	})
}

func (l *websocketLogger) NewVersionAvailable(availableVersion string) {
	l.write(appNewVersionMessage{
		Type:    MessageTypeLoadingFailure,
		Version: availableVersion,
	})
}
