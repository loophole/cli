package cache

import (
	"os"
	"path"

	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/mitchellh/go-homedir"
)

// GetLocalStorageDir returns local directory for loophole cache purposes
func GetLocalStorageDir(directoryName string) string {
	home, err := homedir.Dir()
	if err != nil {
		communication.LogFatalErr("Error reading user home directory ", err)
	}

	dirName := path.Join(home, ".loophole", directoryName)
	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		communication.LogFatalErr("Error creating local cache directory", err)
	}
	return dirName
}

// GetLocalStorageFile returns local file for loophole cache purposes
func GetLocalStorageFile(fileName string, directoryName string) string {
	home, err := homedir.Dir()
	if err != nil {
		communication.LogFatalErr("Error reading user home directory ", err)
	}
	dirName := path.Join(home, ".loophole", directoryName)
	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		communication.LogFatalErr("Error creating local cache directory", err)
	}

	return path.Join(dirName, fileName)
}
