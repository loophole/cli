package cache

import (
	"fmt"
	"os"
	"path"

	"github.com/loophole/cli/internal/pkg/communication"
	"github.com/mitchellh/go-homedir"
)

// GetLocalStorageDir returns local directory for loophole cache purposes
func GetLocalStorageDir(directoryName string) string {
	home, err := homedir.Dir()
	if err != nil {
		communication.Fatal(fmt.Sprintf("Error reading user home directory: %s", err.Error()))
	}

	dirName := path.Join(home, ".loophole", directoryName)
	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		communication.Fatal(fmt.Sprintf("Error creating local cache directory: %s", err.Error()))
	}
	return dirName
}

// GetLocalStorageFile returns local file for loophole cache purposes
func GetLocalStorageFile(fileName string, directoryName string) string {
	home, err := homedir.Dir()
	if err != nil {
		communication.Fatal(fmt.Sprintf("Error reading user home directory: %s", err.Error()))
	}
	dirName := path.Join(home, ".loophole", directoryName)
	err = os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		communication.Fatal(fmt.Sprintf("Error creating local cache directory: %s", err.Error()))
	}

	return path.Join(dirName, fileName)
}
