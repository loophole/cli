//necessary bits and pieces for being able to serve e.g. the custom notification html
// without needing to create files in the users folders
package customfilesystem

import (
	"bytes"
	"errors"
	"io/fs"
	"net/http"
	"path/filepath"
)

var DirectoryListingDisabledPage = []byte("<!DOCTYPE html><html><body><img src=\"https://raw.githubusercontent.com/loophole/website/master/static/img/logo.png\" alt=\"https://raw.githubusercontent.com/loophole/website/master/static/img/logo.png\" class=\"transparent shrinkToFit\" width=\"400\" height=\"88\"><p>Directory index listing has been disabled. Please enter the path of a file.</p></body></html>")

type CustomFileSystem struct {
	FS http.FileSystem
}

//the file cannot be reused since it's io.Reader can only be read from once,
// so we need a reusable way to create it
func writeDirectoryListingDisabledPageFile(pageFile *MyFile) {
	*pageFile = MyFile{
		Reader: bytes.NewReader(DirectoryListingDisabledPage),
		mif: myFileInfo{
			name: "customIndex.html",
			data: DirectoryListingDisabledPage,
		},
	}
}

func (cfs CustomFileSystem) Open(path string) (http.File, error) {
	f, err := _Open(path, cfs)

	if err != nil {
		var pathErrorInstance error = &fs.PathError{
			Err: errors.New(""),
		}
		if errors.As(err, &pathErrorInstance) {
			return nil, err
		}
		var pageFile *MyFile = &MyFile{}
		writeDirectoryListingDisabledPageFile(pageFile)
		return pageFile, nil
	}
	return f, nil
}

//if there is an elegant way to integrate the following into the function above without
// using labeled breaks or adding even more control structures let me know
func _Open(path string, cfs CustomFileSystem) (http.File, error) {
	f, err := cfs.FS.Open(path)
	if err != nil {
		if path == "/" {
			var pageFile *MyFile = &MyFile{}
			writeDirectoryListingDisabledPageFile(pageFile)
			return pageFile, nil
		}
		return nil, err
	}

	s, err := f.Stat()
	if err != nil {
		return nil, err
	}

	if s.IsDir() {
		index := filepath.Join(path, "index.html")
		if _, err := cfs.FS.Open(index); err != nil {
			var pageFile *MyFile = &MyFile{}
			writeDirectoryListingDisabledPageFile(pageFile)
			return pageFile, nil
		}
	}
	return f, nil
}
