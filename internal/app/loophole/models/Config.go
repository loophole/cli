package models

// Config represent loophole configuration
type Config struct {
	Port            int32
	Host            string
	IdentityFile    string
	GatewayEndpoint Endpoint
	SiteID          string
	LogLevel        string
}
