package tmp

import (
	"io/ioutil"
	"os"
)

type TmpDir interface {
	Path() string
	Delete() error
}

func NewDirectory(path string) (TmpDir, error) {
	newDir, _ := ioutil.TempDir(path, `go-tmpdir`)
	return &tmpDir{newDir}, nil
}

type tmpDir struct {
	path string
}

func (t tmpDir) Path() string {
	return t.path
}

func (t tmpDir) Delete() error {
	return os.Remove(t.Path())
}
