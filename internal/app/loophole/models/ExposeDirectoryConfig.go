package models

// ExposeDirectoryConfig represents loophole configuration when directory is exposed
type ExposeDirectoryConfig struct {
	Local  LocalDirectorySpecs `json:"local"`
	Remote RemoteEndpointSpecs `json:"remote"`
}
