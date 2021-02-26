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
