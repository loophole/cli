package models

// ExposeWebdavConfig represents loophole configuration when directory is exposed via webdav
type ExposeWebdavConfig struct {
	Local  LocalDirectorySpecs `json:"local"`
	Remote RemoteEndpointSpecs `json:"remote"`
}
