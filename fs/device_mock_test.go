package fs

import (
	"github.com/stretchr/testify/mock"
)

type DeviceMock struct {
	mock.Mock
}

func (m DeviceMock) Path() Path {
	args := m.Called()
	return args.Get(0).(Path)
}
