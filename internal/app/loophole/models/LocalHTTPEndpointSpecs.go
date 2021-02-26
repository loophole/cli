package models

import "fmt"

// LocalHTTPEndpointSpecs is collection of parameters used to describe
// configuration for local port to be exposed
type LocalHTTPEndpointSpecs struct {
	Port  int32  `json:"port"`
	Host  string `json:"host"`
	HTTPS bool   `json:"https"`
	Path  string `json:"path"`
}

func Validate(options *LocalHTTPEndpointSpecs) error {
	if options.Port <= 0 {
		return fmt.Errorf("Port not set")
	}
	if options.Host == "" {
		return fmt.Errorf("Host not set")
	}
	return nil
}
