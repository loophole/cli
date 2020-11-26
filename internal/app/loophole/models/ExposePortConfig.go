package models

// ExposeHttpConfig represent loophole configuration when port is exposed
type ExposeHttpConfig struct {
	Local   LocalHttpEndpointSpecs
	Remote  RemoteEndpointSpecs
	Display DisplayOptions
}
