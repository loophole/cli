package cache

import (
	"fmt"
	"os"
	"testing"

	"github.com/mitchellh/go-homedir"
)

func TestGetLocalStorageDirReturnsCorrectPath(t *testing.T) {
	home, err := homedir.Dir()
	if err != nil {
		t.Fatal(err)
	}
	dirname := "expected-dirname"
	expectedPath := fmt.Sprintf("%s/.loophole/%s", home, dirname)
	createdDir := GetLocalStorageDir(dirname)

	if expectedPath != createdDir {
		t.Fatalf("Created directory path '%s' is different than expected: '%s'", createdDir, expectedPath)
	}

}

func TestGetLocalStorageDirCreatesDirectory(t *testing.T) {
	dirname := "expected-dirname"
	createdDir := GetLocalStorageDir(dirname)

	info, err := os.Stat(createdDir)
	if os.IsNotExist(err) {
		t.Fatalf("Directory '%s' doesn't exist", createdDir)
	}
	if !info.IsDir() {
		t.Fatalf("Path '%s' doesn't point to directory", createdDir)
	}
}

func TestGetLocalStorageFileReturnsCorrectPath(t *testing.T) {
	home, err := homedir.Dir()
	if err != nil {
		t.Fatal(err)
	}
	filename := "expected-filename"

	expectedPath := fmt.Sprintf("%s/.loophole/%s", home, filename)
	filePath := GetLocalStorageFile(filename, "")

	if expectedPath != filePath {
		t.Fatalf("Created directory path '%s' is different than expected: '%s'", filePath, expectedPath)
	}
}

func TestGetLocalStorageDirDoesntCreateFile(t *testing.T) {
	filename := "expected-filename"
	filePath := GetLocalStorageFile(filename, "")

	_, err := os.Stat(filePath)
	if !os.IsNotExist(err) {
		t.Fatalf("Error while executing stat on '%s'", filePath)
	}
}
