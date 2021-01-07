package config

type OAuthConfig struct {
	DeviceCodeURL string `json:"deviceCodeUrl"`
	TokenURL      string `json:"tokenUrl"`
	ClientID      string `json:"clientId"`
	Scope         string `json:"scope"`
	Audience      string `json:"audience"`
}

type DisplayConfig struct {
	Version    string
	CommitHash string

	Verbose               bool   `json:"verbose"`
	QR                    bool   `json:"qr"`
	FeedbackFormURL       string `json:"feedbackFormUrl"`
	DisableProxyErrorPage bool   `json:"disableProxyErrorPage"`
}

type ApplicationConfig struct {
	OAuthConfig   OAuthConfig   `json:"oauthConfig"`
	DisplayConfig DisplayConfig `json:"displayConfig"`
}

var (
	// Will be filled in during build
	version = "development"
	commit  = "unknown"
)

var Config = ApplicationConfig{
	OAuthConfig: OAuthConfig{
		DeviceCodeURL: "https://loophole.eu.auth0.com/oauth/device/code",
		TokenURL:      "https://loophole.eu.auth0.com/oauth/token",
		ClientID:      "9ocnSAnfJSb6C52waL8xcPidCkRhUwBs",
		Scope:         "openid offline_access",
		Audience:      "https://api.loophole.cloud",
	},
	DisplayConfig: DisplayConfig{
		Version:         version,
		CommitHash:      commit,
		FeedbackFormURL: "https://bit.ly/3mvmZBA",
	},
}
