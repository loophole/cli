package models

// DisplayOptions represents configuration used to display certain things in CLI
type DisplayOptions struct {
	Verbose               bool   `json:"verbose"`
	QR                    bool   `json:"qr"`
	FeedbackFormURL       string `json:"feedbackFormUrl"`
	DisableProxyErrorPage bool   `json:"disableProxyErrorPage"`
	Version               string `json:"version"`
}
