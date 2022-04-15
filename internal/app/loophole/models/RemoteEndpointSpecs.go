package models

// RemoteEndpointSpecs is collection of parameters used to describe
// configuration for public endpoint
type RemoteEndpointSpecs struct {
	GatewayEndpoint       Endpoint `json:"gatewayEndpoint"`
	APIEndpoint           Endpoint `json:"apiEndpoint"`
	IdentityFile          string   `json:"identityFile"`
	SiteID                string   `json:"siteId"`
	Domain                string   `json:"domain"`
	TunnelID              string   `json:"tunnelId"`
	BasicAuthUsername     string   `json:"basicAuthUsername"`
	BasicAuthPassword     string   `json:"basicAuthPassword"`
	DisableProxyErrorPage bool     `json:"disableProxyErrorPage"`
	DisableOldCiphers     bool     `json:"disableOldCiphers"`
}
