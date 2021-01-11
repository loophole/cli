package communication

import (
	tm "github.com/loophole/cli/internal/pkg/token/models"
)

var communicationMechanism Mechanism = NewStdOutLogger()

type Mechanism interface {
	PrintLoginPrompt(tm.DeviceCodeSpec)
	PrintWelcomeMessage(bool)
	PrintTunnelSuccessMessage(siteID string, protocols []string, localAddr string, displayQR bool)
	PrintGoodbyeMessage()
	PrintFeedbackMessage(feedbackFormURL string)
	StartLoading(message string)
	LoadingSuccess()
	LoadingFailure()
	LogInfo(message string)
	LogWarnErr(message string, err error)
	LogFatalErr(message string, err error)
	LogFatalMsg(message string)
	LogDebug(message string)
	NewLine()
	Write(message string)
	WriteRed(message string)
	WriteGreen(message string)
	WriteCyan(message string)
	WriteItalic(message string)
	QRCode(siteAddr string)
}

func SetCommunicationMechanism(mechanism Mechanism) {
	communicationMechanism = mechanism
}

func PrintLoginPrompt(deviceCodeSpec tm.DeviceCodeSpec) {
	communicationMechanism.PrintLoginPrompt(deviceCodeSpec)
}

func PrintWelcomeMessage(loggedIn bool) {
	communicationMechanism.PrintWelcomeMessage(loggedIn)
}

func PrintTunnelSuccessMessage(siteID string, protocols []string, localAddr string, displayQR bool) {
	communicationMechanism.PrintTunnelSuccessMessage(siteID, protocols, localAddr, displayQR)
}

func PrintGoodbyeMessage() {
	communicationMechanism.PrintGoodbyeMessage()
}

func PrintFeedbackMessage(feedbackFormURL string) {
	communicationMechanism.PrintFeedbackMessage(feedbackFormURL)
}

func StartLoading(message string) {
	communicationMechanism.StartLoading(message)
}

func LoadingSuccess() {
	communicationMechanism.LoadingSuccess()
}

func LoadingFailure() {
	communicationMechanism.LoadingFailure()
}

func LogInfo(message string) {
	communicationMechanism.LogInfo(message)
}

func LogWarnErr(message string, err error) {
	communicationMechanism.LogWarnErr(message, err)
}
func LogFatalErr(message string, err error) {
	communicationMechanism.LogFatalErr(message, err)
}

func LogFatalMsg(message string) {
	communicationMechanism.LogFatalMsg(message)
}

func LogDebug(message string) {
	communicationMechanism.LogDebug(message)
}

func NewLine() {
	communicationMechanism.NewLine()
}

func Write(message string) {
	communicationMechanism.Write(message)
}

func WriteRed(message string) {
	communicationMechanism.WriteRed(message)
}
func WriteGreen(message string) {
	communicationMechanism.WriteGreen(message)
}

func WriteCyan(message string) {
	communicationMechanism.WriteCyan(message)
}

func WriteItalic(message string) {
	communicationMechanism.WriteItalic(message)
}

func QRCode(siteAddr string) {
	communicationMechanism.QRCode(siteAddr)
}
