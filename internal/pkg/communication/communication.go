package communication

import (
	coreModels "github.com/loophole/cli/internal/app/loophole/models"
	authModels "github.com/loophole/cli/internal/pkg/token/models"
)

var defaultLogger = NewStdOutLogger()
var communicationMechanism Mechanism = defaultLogger

// Mechanism is a type defining interface for loophole communication
type Mechanism interface {
	Debug(message string)
	Info(message string)
	Warn(message string)
	Error(message string)
	Fatal(message string)

	TunnelDebug(tunnelID string, message string)
	TunnelInfo(tunnelID string, message string)
	TunnelWarn(tunnelID string, message string)
	TunnelError(tunnelID string, message string)

	ApplicationStart(loggedIn bool, idToken string)
	ApplicationStop()

	TunnelStart(tunnelID string)

	TunnelStartSuccess(remoteConfig coreModels.RemoteEndpointSpecs, localEndpoint string)
	TunnelStartFailure(tunnelID string, err error)

	TunnelStopSuccess(tunnelID string)

	LoginStart(authModels.DeviceCodeSpec)
	LoginSuccess(idToken string)
	LoginFailure(err error)

	LogoutSuccess()
	LogoutFailure(err error)

	LoadingStart(tunnelID string, loaderMessage string)
	LoadingSuccess(tunnelID string)
	LoadingFailure(tunnelID string, err error)

	NewVersionAvailable(availableVersion string)
}

// SetCommunicationMechanism is communication mechanism switcher
func SetCommunicationMechanism(mechanism Mechanism) {
	communicationMechanism = mechanism
}

// TunnelDebug is debug level logger in context of a tunnel
func TunnelDebug(tunnelID string, message string) {
	communicationMechanism.TunnelDebug(tunnelID, message)
}

// TunnelInfo is info level logger in context of a tunnel
func TunnelInfo(tunnelID string, message string) {
	communicationMechanism.TunnelInfo(tunnelID, message)
}

// TunnelWarn is warn level logger in context of a tunnel
func TunnelWarn(tunnelID string, message string) {
	communicationMechanism.TunnelWarn(tunnelID, message)
}

// TunnelError is error level logger in context of a tunnel
func TunnelError(tunnelID string, message string) {
	communicationMechanism.TunnelError(tunnelID, message)
}

// Debug is debug level logger
func Debug(message string) {
	communicationMechanism.Debug(message)
}

// Info is info level logger
func Info(message string) {
	communicationMechanism.Info(message)
}

// Warn is warn level logger
func Warn(message string) {
	communicationMechanism.Warn(message)
}

// Error is error level logger
func Error(message string) {
	communicationMechanism.Error(message)
}

// Fatal is fatal level logger, which should cause application to stop
func Fatal(message string) {
	communicationMechanism.Fatal(message)
}

// ApplicationStart is the application startup welcome communicate
func ApplicationStart(loggedIn bool, idToken string) {
	communicationMechanism.ApplicationStart(loggedIn, idToken)
}

// ApplicationStop is the application startup goodbye communicate
func ApplicationStop() {
	communicationMechanism.ApplicationStop()
}

// LoginStart is the communicate to notify about login process being started
func LoginStart(deviceCodeSpec authModels.DeviceCodeSpec) {
	communicationMechanism.LoginStart(deviceCodeSpec)
}

// LoginSuccess is the application success login communicate
func LoginSuccess(idToken string) {
	communicationMechanism.LoginSuccess(idToken)
}

// LoginFailure is the application login failure communicate
func LoginFailure(err error) {
	communicationMechanism.LoginFailure(err)
}

// LogoutSuccess is the notification logout success communicate
func LogoutSuccess() {
	communicationMechanism.LogoutSuccess()
}

// LogoutFailure is the notification logout failure communicate
func LogoutFailure(err error) {
	communicationMechanism.LogoutFailure(err)
}

// TunnelStart is the notification about tunnel registration success
func TunnelStart(tunnelID string) {
	communicationMechanism.TunnelStart(tunnelID)
}

// TunnelStartSuccess is the notification about tunnel being started succesfully
func TunnelStartSuccess(remoteConfig coreModels.RemoteEndpointSpecs, localEndpoint string) {
	communicationMechanism.TunnelStartSuccess(remoteConfig, localEndpoint)
}

// TunnelStartFailure is the notification about tunnel failing to start
func TunnelStartFailure(tunnelID string, err error) {
	communicationMechanism.TunnelStartFailure(tunnelID, err)
}

// TunnelRestart is the notification about tunnel being restarted
func TunnelRestart(tunnelID string) {}

// TunnelStopSuccess is the notification about tunnel being shut down
func TunnelStopSuccess(tunnelID string) {
	communicationMechanism.TunnelStopSuccess(tunnelID)
}

// LoadingStart is the notification about some loading process being started
func LoadingStart(tunnelID string, loaderMessage string) {
	communicationMechanism.LoadingStart(tunnelID, loaderMessage)
}

// LoadingSuccess is the notification about started loading process being finished successfully
func LoadingSuccess(tunnelID string) {
	communicationMechanism.LoadingSuccess(tunnelID)
}

// LoadingFailure is the notification about started loading process being finished with failure
func LoadingFailure(tunnelID string, err error) {
	communicationMechanism.LoadingFailure(tunnelID, err)
}

// NewVersionAvailable is a communicate being sent if new version of the application is available
func NewVersionAvailable(availableVersion string) {
	communicationMechanism.NewVersionAvailable(availableVersion)
}
