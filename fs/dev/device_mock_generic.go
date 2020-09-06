package dev

import (
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/path"
	"github.com/stretchr/testify/mock"
)

type DeviceMock struct {
	mock.Mock
}

func (m *DeviceMock) Path() path.Path {
	args := m.Called()
	return args.Get(0).(path.Path)
}
