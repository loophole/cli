package models

// LocalHttpEndpointSpecs is collection of parameters used to describe
// configuration for local port to be exposed
type LocalHttpEndpointSpecs struct {
	Port  int32
	Host  string
	HTTPS bool
}
