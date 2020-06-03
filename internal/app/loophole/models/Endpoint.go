package models

import "fmt"

// Endpoint is representing host address
type Endpoint struct {
	Host string
	Port int32
}

func (endpoint *Endpoint) String() string {
	if endpoint != nil {
		return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
	}
	return ""
}
