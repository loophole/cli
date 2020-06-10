package models

// Config represent loophole configuration
type Config struct {
	Port            int32
	Host            string
	IdentityFile    string
	GatewayEndpoint Endpoint
	APIURL          string
	APIPort         int32
	SiteID          string
}
