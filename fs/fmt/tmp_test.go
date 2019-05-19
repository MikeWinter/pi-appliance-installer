package fmt

import (
	"errors"
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/os"
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/path"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TmpTestSuite struct {
	suite.Suite

	tests []struct {
		Name         string
		Path         path.Path
		ExpectedArgs []interface{}
	}

	adapterMock *os.OsAdapterMock
	rootFs      Filesystem
	fs          Filesystem
}

func TestTmpTestSuite(t *testing.T) {
	suite.Run(t, new(TmpTestSuite))
}

func (s *TmpTestSuite) SetupSuite() {
	s.tests = []struct {
		Name         string
		Path         path.Path
		ExpectedArgs []interface{}
	}{
		{
			Name:         "passes tmpfs format type to adapter",
			ExpectedArgs: []interface{}{mock.Anything, mock.Anything, "tmpfs", mock.Anything, mock.Anything},
		},
		{
			Name:         "passes the absolute mount path to adapter",
			Path:         "/mount-point",
			ExpectedArgs: []interface{}{mock.Anything, "/root/mount-point", mock.Anything, mock.Anything, mock.Anything},
		},
		{
			Name:         "provides no source to adapter",
			ExpectedArgs: []interface{}{"", mock.Anything, mock.Anything, mock.Anything, mock.Anything},
		},
		{
			Name:         "provides no flags to adapter",
			ExpectedArgs: []interface{}{mock.Anything, mock.Anything, mock.Anything, uintptr(0), mock.Anything},
		},
		{
			Name:         "provides no options to adapter",
			ExpectedArgs: []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, ""},
		},
	}
}

func (s *TmpTestSuite) SetupTest() {
	s.rootFs = &fs{path: "/root"}
	s.fs = NewTmp()

	s.adapterMock = new(os.OsAdapterMock)
	os.Adapter = s.adapterMock
}

func (s *TmpTestSuite) TestMountReturnsNilOnSuccess() {
	s.adapterMock.
		On("Mount", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	err := s.rootFs.Mount(s.fs, "/tmp")

	s.NoError(err)
}

func (s *TmpTestSuite) TestMountForwardsErrorOnFailure() {
	s.adapterMock.
		On("Mount", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("mount error"))

	err := s.rootFs.Mount(s.fs, "/tmp")

	s.EqualError(err, "mount error")
}

func (s *TmpTestSuite) TestMountInvokesAdapterWithArguments() {
	for _, test := range s.tests {
		s.Run(test.Name, func() {
			s.SetupTest()
			s.adapterMock.
				On("Mount", test.ExpectedArgs...).
				Return(nil)

			err := s.rootFs.Mount(s.fs, test.Path)

			require.NoError(s.T(), err)
			s.adapterMock.AssertExpectations(s.T())
		})
	}
}
