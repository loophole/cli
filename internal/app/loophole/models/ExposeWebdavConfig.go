package models

// ExposeWebdavConfig represents loophole configuration when directory is exposed via webdav
type ExposeWebdavConfig struct {
	Local   LocalDirectorySpecs
	Remote  RemoteEndpointSpecs
	Display DisplayOptions
}
