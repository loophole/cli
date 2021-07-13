package communication

import (
	"github.com/loophole/cli/config"
	coreModels "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/logger"
	authModels "github.com/loophole/cli/internal/pkg/token/models"
)

// SetCommunicationMechanism is communication mechanism switcher
func SetCommunicationMechanism(mechanism coreModels.CommunicationMechanism) {
	logger.CommunicationMechanism = mechanism
}

// TunnelDebug is debug level logger in context of a tunnel
func TunnelDebug(tunnelID string, message string) {
	logger.CommunicationMechanism.TunnelDebug(tunnelID, message)
}

// TunnelInfo is info level logger in context of a tunnel
func TunnelInfo(tunnelID string, message string) {
	logger.CommunicationMechanism.TunnelInfo(tunnelID, message)
}

// TunnelWarn is warn level logger in context of a tunnel
func TunnelWarn(tunnelID string, message string) {
	logger.CommunicationMechanism.TunnelWarn(tunnelID, message)
}

// TunnelError is error level logger in context of a tunnel
func TunnelError(tunnelID string, message string) {
	logger.CommunicationMechanism.TunnelError(tunnelID, message)
}

// Debug is debug level logger
func Debug(message string) {
	logger.CommunicationMechanism.Debug(message)
}

// Info is info level logger
func Info(message string) {
	logger.CommunicationMechanism.Info(message)
}

// Warn is warn level logger
func Warn(message string) {
	logger.CommunicationMechanism.Warn(message)
}

// Error is error level logger
func Error(message string) {
	logger.CommunicationMechanism.Error(message)
}

// Fatal is fatal level logger, which should cause application to stop
func Fatal(message string) {
	logger.CommunicationMechanism.Fatal(message)
}

// ApplicationStart is the application startup welcome communicate
func ApplicationStart(loggedIn bool, idToken string) {
	logger.CommunicationMechanism.ApplicationStart(loggedIn, idToken)
}

// ApplicationStop is the application startup goodbye communicate
func ApplicationStop() {
	logger.CommunicationMechanism.ApplicationStop(config.Config.FeedbackFormURL)
}

// LoginStart is the communicate to notify about login process being started
func LoginStart(deviceCodeSpec authModels.DeviceCodeSpec) {
	logger.CommunicationMechanism.LoginStart(deviceCodeSpec)
}

// LoginSuccess is the application success login communicate
func LoginSuccess(idToken string) {
	logger.CommunicationMechanism.LoginSuccess(idToken)
}

// LoginFailure is the application login failure communicate
func LoginFailure(err error) {
	logger.CommunicationMechanism.LoginFailure(err)
}

// LogoutSuccess is the notification logout success communicate
func LogoutSuccess() {
	logger.CommunicationMechanism.LogoutSuccess()
}

// LogoutFailure is the notification logout failure communicate
func LogoutFailure(err error) {
	logger.CommunicationMechanism.LogoutFailure(err)
}

// TunnelStart is the notification about tunnel registration success
func TunnelStart(tunnelID string) {
	logger.CommunicationMechanism.TunnelStart(tunnelID)
}

// TunnelStartSuccess is the notification about tunnel being started succesfully
func TunnelStartSuccess(remoteConfig coreModels.RemoteEndpointSpecs, localEndpoint string) {
	logger.CommunicationMechanism.TunnelStartSuccess(remoteConfig, localEndpoint, config.Config.Display.QR)
}

// TunnelStartFailure is the notification about tunnel failing to start
func TunnelStartFailure(tunnelID string, err error) {
	logger.CommunicationMechanism.TunnelStartFailure(tunnelID, err)
}

// TunnelRestart is the notification about tunnel being restarted
func TunnelRestart(tunnelID string) {}

// TunnelStopSuccess is the notification about tunnel being shut down
func TunnelStopSuccess(tunnelID string) {
	logger.CommunicationMechanism.TunnelStopSuccess(tunnelID)
}

// LoadingStart is the notification about some loading process being started
func LoadingStart(tunnelID string, loaderMessage string) {
	logger.CommunicationMechanism.LoadingStart(tunnelID, loaderMessage)
}

// LoadingSuccess is the notification about started loading process being finished successfully
func LoadingSuccess(tunnelID string) {
	logger.CommunicationMechanism.LoadingSuccess(tunnelID)
}

// LoadingFailure is the notification about started loading process being finished with failure
func LoadingFailure(tunnelID string, err error) {
	logger.CommunicationMechanism.LoadingFailure(tunnelID, err)
}

// NewVersionAvailable is a communicate being sent if new version of the application is available
func NewVersionAvailable(availableVersion string) {
	logger.CommunicationMechanism.NewVersionAvailable(availableVersion)
}
