// +build !dev

package config

import "github.com/loophole/cli/internal/app/loophole/models"

// Config is global application config
var Config = ApplicationConfig{
	Version:         "development",
	CommitHash:      "unknown",
	ClientMode:      "unknown",
	FeedbackFormURL: "https://bit.ly/3mvmZBA",

	OAuth: OAuthConfig{
		DeviceCodeURL: "https://loophole.eu.auth0.com/oauth/device/code",
		TokenURL:      "https://loophole.eu.auth0.com/oauth/token",
		ClientID:      "9ocnSAnfJSb6C52waL8xcPidCkRhUwBs",
		Scope:         "openid offline_access profile email",
		Audience:      "https://api.loophole.cloud",
	},
	Display: DisplayConfig{
		Verbose: false,
		QR:      false,
	},
	APIEndpoint: models.Endpoint{
		Protocol: "https",
		Host:     "api.loophole.cloud",
		Port:     443,
	},
	GatewayEndpoint: models.Endpoint{
		Protocol: "ssh",
		Host:     "gateway.loophole.host",
		Port:     8022,
	},
}
