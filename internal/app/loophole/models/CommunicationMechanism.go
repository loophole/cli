package models

import authModels "github.com/loophole/cli/internal/pkg/token/models"

// Mechanism is a type defining interface for loophole communication
type CommunicationMechanism interface {
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
	ApplicationStop(feedbackFormURL ...string)

	TunnelStart(tunnelID string)

	TunnelStartSuccess(remoteConfig RemoteEndpointSpecs, localEndpoint string, qr ...bool)
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
