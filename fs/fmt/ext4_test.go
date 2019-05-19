package fmt

import (
	"errors"
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/dev"
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/os"
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/path"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"testing"
)

type Ext4TestSuite struct {
	suite.Suite

	tests []struct {
		Name         string
		DevicePath   path.Path
		MountPath    path.Path
		ExpectedArgs []interface{}
		ExpectedErr  error
	}

	adapterMock *os.OsAdapterMock
	deviceMock  *dev.DeviceMock
	rootFs      Filesystem
	fs          Filesystem
}

func TestExt4TestSuite(t *testing.T) {
	suite.Run(t, new(Ext4TestSuite))
}

func (s *Ext4TestSuite) SetupSuite() {
	s.tests = []struct {
		Name         string
		DevicePath   path.Path
		MountPath    path.Path
		ExpectedArgs []interface{}
		ExpectedErr  error
	}{
		{
			Name:         "Mount returns nil on success",
			ExpectedArgs: []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
		},
		{
			Name:         "Mount forwards error on failure",
			ExpectedArgs: []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
			ExpectedErr:  errors.New("mount error"),
		},
		{
			Name:         "Mount uses device as source",
			DevicePath:   "/dev/device-name",
			ExpectedArgs: []interface{}{"/dev/device-name", mock.Anything, mock.Anything, mock.Anything, mock.Anything},
		},
		{
			Name:         "Mount uses absolute path as target",
			MountPath:    "/mnt/path",
			ExpectedArgs: []interface{}{mock.Anything, "/root/mnt/path", mock.Anything, mock.Anything, mock.Anything},
		},
		{
			Name:         "Mount specifies ext4 filesystem format",
			ExpectedArgs: []interface{}{mock.Anything, mock.Anything, "ext4", mock.Anything, mock.Anything},
		},
		{
			Name:         "Mount sets no flags",
			ExpectedArgs: []interface{}{mock.Anything, mock.Anything, mock.Anything, uintptr(0), mock.Anything},
		},
		{
			Name:         "Mount sets no options",
			ExpectedArgs: []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, ""},
		},
	}
}

func (s *Ext4TestSuite) SetupTest() {
	s.deviceMock = new(dev.DeviceMock)
	s.rootFs = &fs{path: "/root"}
	s.fs = NewExt4(s.deviceMock)

	s.adapterMock = new(os.OsAdapterMock)
	os.Adapter = s.adapterMock
}

func (s *Ext4TestSuite) TestMount() {
	for _, test := range s.tests {
		s.Run(test.Name, func() {
			s.SetupTest()
			s.adapterMock.
				On("Mount", test.ExpectedArgs...).
				Return(test.ExpectedErr)
			s.deviceMock.
				On("Path").
				Return(test.DevicePath)

			err := s.rootFs.Mount(s.fs, test.MountPath)

			s.Exactly(test.ExpectedErr, err)
			s.adapterMock.AssertExpectations(s.T())
		})
	}
}
