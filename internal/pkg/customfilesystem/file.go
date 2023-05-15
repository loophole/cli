//taken from https://stackoverflow.com/questions/52697277/simples-way-to-make-a-byte-into-a-virtual-file-object-in-golang
package customfilesystem

import (
	"bytes"
	"io"
	"os"
	"time"
)

type File interface {
	io.Closer
	io.Reader
	io.Seeker
	Readdir(count int) ([]os.FileInfo, error)
	Stat() (os.FileInfo, error)
}

type myFileInfo struct {
	name string
	data []byte
}

func (mif myFileInfo) Name() string       { return mif.name }
func (mif myFileInfo) Size() int64        { return int64(len(mif.data)) }
func (mif myFileInfo) Mode() os.FileMode  { return 0444 }        // Read for all
func (mif myFileInfo) ModTime() time.Time { return time.Time{} } // Return anything
func (mif myFileInfo) IsDir() bool        { return false }
func (mif myFileInfo) Sys() interface{}   { return nil }

type MyFile struct {
	*bytes.Reader
	mif myFileInfo
}

func (mf *MyFile) Close() error { return nil } // Noop, nothing to do

func (mf *MyFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil // We are not a directory but a single file
}

func (mf *MyFile) Stat() (os.FileInfo, error) {
	return mf.mif, nil
}
