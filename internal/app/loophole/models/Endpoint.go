package models

import "fmt"

// Endpoint is representing host address
type Endpoint struct {
	Protocol string
	Host     string
	Port     int32
}

// URI returns the full uri string protocol://host:port
func (endpoint *Endpoint) URI() string {
	if endpoint.Protocol != "" {
		return fmt.Sprintf("%s://%s:%d", endpoint.Protocol, endpoint.Host, endpoint.Port)
	}
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}

// Hostname returns the hostname part of endpoint (not including protocol)
func (endpoint *Endpoint) Hostname() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}
