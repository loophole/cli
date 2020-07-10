package cache

import (
	"os"
	"path"

	"github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog/log"
)

// GetLocalStorageDir returns local directory for loophole cache purposes
func GetLocalStorageDir(directoryName string) string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading user home directory ")
	}

	dirName := path.Join(home, ".loophole", directoryName)
	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating local cache directory")
	}
	return dirName
}

// GetLocalStorageFile returns local file for loophole cache purposes
func GetLocalStorageFile(fileName string) string {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal().Err(err).Msg("Error reading user home directory ")
	}
	dirName := path.Join(home, ".loophole")
	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating local cache directory")
	}

	return path.Join(dirName, fileName)
}
