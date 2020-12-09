package models

// RemoteEndpointSpecs is collection of parameters used to describe
// configuration for public endpoint
type RemoteEndpointSpecs struct {
	GatewayEndpoint   Endpoint
	APIEndpoint       Endpoint
	IdentityFile      string
	SiteID            string
	BasicAuthUsername string
	BasicAuthPassword string
}
