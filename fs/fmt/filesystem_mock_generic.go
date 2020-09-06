package fmt

import (
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/path"
	"github.com/stretchr/testify/mock"
)

type FilesystemMock struct {
	mock.Mock
}

func (m *FilesystemMock) AsRoot() Filesystem {
	args := m.Called()
	return args.Get(0).(Filesystem)
}

func (m *FilesystemMock) Mount(child Filesystem, at path.Path) error {
	args := m.Called(child, at)
	return args.Error(0)
}

func (m *FilesystemMock) Path() path.Path {
	args := m.Called()
	return args.Get(0).(path.Path)
}

func (m *FilesystemMock) Pivot(to Filesystem, old path.Path) (Filesystem, error) {
	args := m.Called(to, old)
	return args.Get(0).(Filesystem), args.Error(1)
}

func (m *FilesystemMock) Remount() error {
	args := m.Called()
	return args.Error(0)
}

func (m *FilesystemMock) Unmount() error {
	args := m.Called()
	return args.Error(0)
}

func (m *FilesystemMock) setPath(path path.Path) {
}

func (m *FilesystemMock) setParent(parent Filesystem) {
}

func (m *FilesystemMock) onMount(at path.Path) error {
	args := m.Called(at)
	return args.Error(0)
}
