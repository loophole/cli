// +build desktop

package ui

import (
	"fmt"

	"github.com/loophole/cli/internal/app/loophole"
	lm "github.com/loophole/cli/internal/app/loophole/models"
	"github.com/loophole/cli/internal/pkg/cache"
	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/loophole/cli/internal/pkg/token"
	"github.com/skratchdot/open-golang/open"
)

type MessageType string

const (
	MessageTypeStartTunnel   MessageType = "MT_StartTunnel"
	MessageTypeStopTunnel    MessageType = "MT_StopTunnel"
	MessageTypeAuthorization MessageType = "MT_Authorize"
	MessageTypeLogout        MessageType = "MT_Logout"
	MessageTypeOpenBrowser   MessageType = "MT_OpenInBrowser"
)

type Message struct {
	MessageType          MessageType          `json:"messageType"`
	StartTunnelMessage   StartTunnelMessage   `json:"startTunnelMessage"`
	StopTunnelMessage    StopTunnelMessage    `json:"stopTunnelMessage"`
	OpenInBrowserMessage OpenInBrowserMessage `json:"openInBrowserMessage"`
}

type StartTunnelMessage struct {
	TunnelType            loophole.TunnelType      `json:"tunnelType"`
	ExposeHTTPConfig      lm.ExposeHTTPConfig      `json:"exposeHttpConfig"`
	ExposeDirectoryConfig lm.ExposeDirectoryConfig `json:"exposeDirectoryConfig"`
	ExposeWebdavConfig    lm.ExposeWebdavConfig    `json:"exposeWebdavConfig"`
}

type StopTunnelMessage struct {
	SiteID string `json:"siteId"`
}

type OpenInBrowserMessage struct {
	URL string `json:"url"`
}

func startAuthorizationProcess(quitChannel <-chan bool) {
	deviceCodeSpec, err := token.RegisterDevice()
	if err != nil {
		communication.LogFatalErr("Error obtaining device code", err)
	}
	communication.PrintLoginPrompt(*deviceCodeSpec)
	tokens, err := token.PollForToken(deviceCodeSpec.DeviceCode, deviceCodeSpec.Interval, quitChannel)
	if err != nil {
		communication.LogWarnErr("Error obtaining token", err)
		communication.PrintWelcomeMessage(token.IsTokenSaved())
		return
	}
	err = token.SaveToken(tokens)
	if err != nil {
		communication.LogWarnErr("Error saving token", err)
		communication.PrintWelcomeMessage(token.IsTokenSaved())
		return
	}
	communication.PrintWelcomeMessage(token.IsTokenSaved())
	communication.LogInfo("Logged in successfully")
}

func logout() {
	err := token.DeleteTokens()
	if err != nil {
		communication.LogWarnErr("Error logging out", err)
		communication.PrintWelcomeMessage(token.IsTokenSaved())
		return
	}
	communication.PrintWelcomeMessage(token.IsTokenSaved())
	communication.LogInfo("Logged out successfully")
}

func startHTTPTunnel(config lm.ExposeHTTPConfig, quitChannel <-chan bool) {
	config.Remote.GatewayEndpoint.Host = "gateway.loophole.host"
	config.Remote.GatewayEndpoint.Port = 8022
	config.Remote.APIEndpoint.Protocol = "https"
	config.Remote.APIEndpoint.Host = "api.loophole.cloud"
	config.Remote.APIEndpoint.Port = 443

	fmt.Println(config)

	sshDir := cache.GetLocalStorageDir(".ssh") //getting our sshDir and creating it, if it doesn't exist
	config.Remote.IdentityFile = fmt.Sprintf("%s/id_rsa", sshDir)

	// err = config.Local.Validate()

	fmt.Println(config)
	loophole.ForwardPort(config, quitChannel)
}

func startDirectoryTunnel(config lm.ExposeDirectoryConfig, quitChannel <-chan bool) {
	config.Remote.GatewayEndpoint.Host = "gateway.loophole.host"
	config.Remote.GatewayEndpoint.Port = 8022
	config.Remote.APIEndpoint.Protocol = "https"
	config.Remote.APIEndpoint.Host = "api.loophole.cloud"
	config.Remote.APIEndpoint.Port = 443

	sshDir := cache.GetLocalStorageDir(".ssh") //getting our sshDir and creating it, if it doesn't exist
	config.Remote.IdentityFile = fmt.Sprintf("%s/id_rsa", sshDir)

	// err = config.Local.Validate()

	fmt.Println(config)
	loophole.ForwardDirectory(config, quitChannel)
}

func startWebdavTunnel(config lm.ExposeWebdavConfig, quitChannel <-chan bool) {
	config.Remote.GatewayEndpoint.Host = "gateway.loophole.host"
	config.Remote.GatewayEndpoint.Port = 8022
	config.Remote.APIEndpoint.Protocol = "https"
	config.Remote.APIEndpoint.Host = "api.loophole.cloud"
	config.Remote.APIEndpoint.Port = 443

	sshDir := cache.GetLocalStorageDir(".ssh") //getting our sshDir and creating it, if it doesn't exist
	config.Remote.IdentityFile = fmt.Sprintf("%s/id_rsa", sshDir)

	// err = config.Local.Validate()

	fmt.Println(config)
	loophole.ForwardDirectoryViaWebdav(config, quitChannel)
}

func openInBrowser(url string) {
	err := open.Run(url)
	if err != nil {
		communication.LogWarnErr("Error opening url", err)
	}
}
