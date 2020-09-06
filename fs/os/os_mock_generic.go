package os

import "github.com/stretchr/testify/mock"

type AdapterMock struct {
	mock.Mock
}

func (m *AdapterMock) Mount(source string, target string, fsType string, flags uintptr, data string) error {
	args := m.Called(source, target, fsType, flags, data)
	return args.Error(0)
}

func (m *AdapterMock) MkNode(path string, mode uint32, dev int) error {
	args := m.Called(path, mode, dev)
	return args.Error(0)
}

func (m *AdapterMock) PivotRoot(newRoot string, putOld string) error {
	args := m.Called(newRoot, putOld)
	return args.Error(0)
}

func (m *AdapterMock) Unmount(target string, flags int) error {
	args := m.Called(target, flags)
	return args.Error(0)
}
