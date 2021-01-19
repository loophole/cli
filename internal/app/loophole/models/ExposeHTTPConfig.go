package models

// ExposeHTTPConfig represents loophole configuration when port is exposed
type ExposeHTTPConfig struct {
	Local  LocalHTTPEndpointSpecs `json:"local"`
	Remote RemoteEndpointSpecs    `json:"remote"`
}
