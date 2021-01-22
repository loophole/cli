package config

import (
	"github.com/loophole/cli/internal/app/loophole/models"
)

// OAuthConfig defined OAuth settings shape
type OAuthConfig struct {
	DeviceCodeURL string `json:"deviceCodeUrl"`
	TokenURL      string `json:"tokenUrl"`
	ClientID      string `json:"clientId"`
	Scope         string `json:"scope"`
	Audience      string `json:"audience"`
}

// DisplayConfig defines the display switches shape
type DisplayConfig struct {
	Verbose bool `json:"verbose"`
	QR      bool `json:"qr"`
}

// ApplicationConfig defines the application config shape
type ApplicationConfig struct {
	Version    string `json:"version"`
	CommitHash string `json:"commitHash"`
	ClientMode string `json:"clientMode"`

	FeedbackFormURL string `json:"feedbackFormUrl"`

	OAuth   OAuthConfig   `json:"oauthConfig"`
	Display DisplayConfig `json:"displayConfig"`

	APIEndpoint     models.Endpoint `json:"apiConfig"`
	GatewayEndpoint models.Endpoint `json:"gatewayConfig"`
}

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
