package models

// ExposeDirectoryConfig represent loophole configuration when directory is exposed
type ExposeDirectoryConfig struct {
	Local   LocalDirectorySpecs
	Remote  RemoteEndpointSpecs
	Display DisplayOptions
}
