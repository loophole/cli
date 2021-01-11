// +build desktop

package communication

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	tm "github.com/loophole/cli/internal/pkg/token/models"
	"github.com/loophole/cli/internal/pkg/urlmaker"
)

type websocketLogger struct {
	wsClient     *websocket.Conn
	messageType  int
	messageMutex sync.Mutex
}

type websocketMessage interface {
	Serialize() string
}

type logEntryMessage struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Class   string `json:"class"`
}

func (m *logEntryMessage) Serialize() []byte {
	m.Type = "Log"
	message, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return message
}

type tunnelMetadataMessage struct {
	Type      string   `json:"type"`
	SiteID    string   `json:"siteId"`
	SiteAddrs []string `json:"siteAddrs"`
	LocalAddr string   `json:"localAddr"`
}

func (m *tunnelMetadataMessage) Serialize() []byte {
	m.Type = "TunnelMetadata"
	message, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return message
}

type tunnelShutdownMessage struct {
	Type string `json:"type"`
}

func (m *tunnelShutdownMessage) Serialize() []byte {
	m.Type = "TunnelShutDown"
	message, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return message
}

type authInfoMessage struct {
	Type     string `json:"type"`
	LoggedIn bool   `json:"loggedIn"`
}

func (m *authInfoMessage) Serialize() []byte {
	m.Type = "AuthorizationInfo"
	message, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return message
}

type authInstructionMessage struct {
	Type                    string `json:"type"`
	DeviceCode              string `json:"deviceCode"`
	UserCode                string `json:"userCode"`
	VerificationURI         string `json:"verificationUri"`
	VerificationURIComplete string `json:"verificationUriComplete"`
}

func (m *authInstructionMessage) Serialize() []byte {
	m.Type = "AuthorizationInstructions"
	message, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return message
}

func NewWebsocketLogger(wsClient *websocket.Conn, messageType int) Mechanism {
	logger := websocketLogger{
		wsClient:    wsClient,
		messageType: messageType,
	}

	return &logger
}

func (l *websocketLogger) PrintLoginPrompt(deviceCodeSpec tm.DeviceCodeSpec) {

	websocketMessage := authInstructionMessage{
		DeviceCode:              deviceCodeSpec.DeviceCode,
		UserCode:                deviceCodeSpec.UserCode,
		VerificationURI:         deviceCodeSpec.VerificationURI,
		VerificationURIComplete: deviceCodeSpec.VerificationURIComplete,
	}

	err := l.wsClient.WriteMessage(l.messageType, []byte(websocketMessage.Serialize()))
	if err != nil {
		fmt.Println("write:", err)
	}
}

func (l *websocketLogger) PrintWelcomeMessage(loggedIn bool) {
	l.WriteWithClass(fmt.Sprintf("Application status: {loggedIn: %t}", loggedIn), "info")
	websocketMessage := authInfoMessage{
		LoggedIn: loggedIn,
	}

	err := l.wsClient.WriteMessage(l.messageType, []byte(websocketMessage.Serialize()))
	if err != nil {
		fmt.Println("write:", err)
	}
}

func (l *websocketLogger) PrintTunnelSuccessMessage(siteID string, protocols []string, localAddr string, displayQR bool) {
	l.WriteWithClass("Tunnel started", "success")
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()

	siteAddrs := []string{}
	for _, protocol := range protocols {
		siteAddrs = append(siteAddrs, urlmaker.GetSiteURL(protocol, siteID))
	}

	websocketMessage := tunnelMetadataMessage{
		SiteID:    siteID,
		SiteAddrs: siteAddrs,
		LocalAddr: localAddr,
	}

	err := l.wsClient.WriteMessage(l.messageType, []byte(websocketMessage.Serialize()))
	if err != nil {
		fmt.Println("write:", err)
	}
}

func (l *websocketLogger) PrintGoodbyeMessage() {
	l.WriteWithClass("Goodbye", "success")
	websocketMessage := tunnelShutdownMessage{}

	err := l.wsClient.WriteMessage(l.messageType, []byte(websocketMessage.Serialize()))
	if err != nil {
		fmt.Println("write:", err)
	}
}

func (l *websocketLogger) PrintFeedbackMessage(feedbackFormURL string) {
	l.WriteWithClass(fmt.Sprintf("Thank you for using Loophole. Please give us your feedback via %s and help us improve our services.", feedbackFormURL), "info")
}

func (l *websocketLogger) StartLoading(message string) {
	l.WriteWithClass(message, "info")
}

func (l *websocketLogger) LoadingSuccess() {
	l.WriteWithClass("Success!", "success")
}

func (l *websocketLogger) LoadingFailure() {
	l.WriteWithClass("Failure!", "danger")
}

func (l *websocketLogger) LogInfo(message string) {
	l.WriteWithClass(message, "info")
}

func (l *websocketLogger) LogWarnErr(message string, err error) {
	l.WriteWithClass(fmt.Sprintf("%s: %s", message, err.Error()), "warning")
}

func (l *websocketLogger) LogFatalErr(message string, err error) {
	l.WriteWithClass(fmt.Sprintf("%s: %s", message, err.Error()), "danger")
	websocketMessage := tunnelShutdownMessage{}

	errSending := l.wsClient.WriteMessage(l.messageType, []byte(websocketMessage.Serialize()))
	if errSending != nil {
		fmt.Println("write:", errSending)
	}
}

func (l *websocketLogger) LogFatalMsg(message string) {
	l.WriteWithClass(message, "danger")
	websocketMessage := tunnelShutdownMessage{}

	errSending := l.wsClient.WriteMessage(l.messageType, []byte(websocketMessage.Serialize()))
	if errSending != nil {
		fmt.Println("write:", errSending)
	}
}

func (l *websocketLogger) LogDebug(message string) {
	l.WriteWithClass(message, "danger")
}

func (l *websocketLogger) NewLine() {

}

func (l *websocketLogger) Write(message string) {
	l.WriteWithClass(message, "")
}

func (l *websocketLogger) WriteWithClass(message string, class string) {
	l.messageMutex.Lock()
	defer l.messageMutex.Unlock()
	websocketMessage := logEntryMessage{
		Message: message,
		Class:   class,
	}

	err := l.wsClient.WriteMessage(l.messageType, []byte(websocketMessage.Serialize()))
	if err != nil {
		fmt.Println("write:", err)
	}
}

func (l *websocketLogger) WriteRed(message string) {
	l.WriteWithClass(message, "hdanger")
}
func (l *websocketLogger) WriteGreen(message string) {
	l.WriteWithClass(message, "success")
}

func (l *websocketLogger) WriteCyan(message string) {
	l.WriteWithClass(message, "info")
}

func (l *websocketLogger) WriteItalic(message string) {
	l.WriteWithClass(message, "info")
}

func (l *websocketLogger) QRCode(siteAddr string) {

}
