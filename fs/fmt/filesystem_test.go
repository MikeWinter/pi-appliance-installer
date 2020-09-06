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

type FsTestSuite struct {
	suite.Suite

	adapterMock *os.AdapterMock
	rootFs      Filesystem
}

func TestFsTestSuite(t *testing.T) {
	suite.Run(t, new(FsTestSuite))
}

func (s *FsTestSuite) SetupTest() {
	s.adapterMock = new(os.AdapterMock)
	s.rootFs = new(fs).AsRoot()

	os.Adapter = s.adapterMock
}

func (s *FsTestSuite) TestRootFilesystemHasRootPath() {
	root := s.rootFs.Path()

	s.EqualValues("/", root)
}

func (s *FsTestSuite) TestMountsFilesystemRelativeToRoot() {
	fs := s.newFilesystem()

	require.NoError(s.T(), s.rootFs.Mount(fs, "/path"))

	s.EqualValues("/path", fs.Path())
}

func (s *FsTestSuite) TestMountsFilesystemRelativeToParent() {
	parentFs := s.newFilesystemAtPath("/parent")
	fs := s.newFilesystem()

	require.NoError(s.T(), parentFs.Mount(fs, "/child"))

	s.EqualValues("/parent/child", fs.Path())
}

var fsRemountTests = []struct {
	Name         string
	ExpectedArgs []interface{}
	ExpectedErr  error
}{
	{
		Name:         "no source",
		ExpectedArgs: []interface{}{"", mock.Anything, mock.Anything, mock.Anything, mock.Anything},
	},
	{
		Name:         "filesystem mount point",
		ExpectedArgs: []interface{}{mock.Anything, "/mnt", mock.Anything, mock.Anything, mock.Anything},
	},
	{
		Name:         "no format",
		ExpectedArgs: []interface{}{mock.Anything, mock.Anything, "", mock.Anything, mock.Anything},
	},
	{
		Name:         "remount flag",
		ExpectedArgs: []interface{}{mock.Anything, mock.Anything, mock.Anything, uintptr(REMOUNT), mock.Anything},
	},
	{
		Name:         "no data",
		ExpectedArgs: []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, ""},
	},
	{
		Name:         "forwards errors from adapter",
		ExpectedArgs: []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything},
		ExpectedErr:  errors.New("adapter error"),
	},
}

func (s *FsTestSuite) TestRemountPassesNecessaryArgumentsToAdapter() {
	for _, test := range fsRemountTests {
		s.Run(test.Name, func() {
			s.SetupTest()
			fs := s.newFilesystemAtPath("/mnt")
			s.adapterMock.
				On("Mount", test.ExpectedArgs...).
				Return(test.ExpectedErr)

			err := fs.Remount()

			if s.Equal(test.ExpectedErr, err) {
				s.adapterMock.AssertExpectations(s.T())
			}
		})
	}
}

func (s *FsTestSuite) TestPivotReturnsOnlyTargetedFilesystemOnSuccess() {
	fs := s.newFilesystem()
	_ = fs.Mount(s.rootFs, "/path")
	s.adapterMock.
		On("PivotRoot", mock.Anything, mock.Anything).
		Return(nil)

	newRoot, err := s.rootFs.Pivot(fs, "/old-root")

	s.NoError(err)
	s.Same(newRoot, fs)
}

func (s *FsTestSuite) TestPivotReturnsOnlyErrorOnFailure() {
	fs := s.newFilesystem()
	_ = fs.Mount(s.rootFs, "/path")
	s.adapterMock.
		On("PivotRoot", mock.Anything, mock.Anything).
		Return(errors.New("pivot error"))

	newRoot, err := s.rootFs.Pivot(fs, "/old-root")

	s.Nil(newRoot)
	s.Error(err)
}

func (s *FsTestSuite) TestPivotReturnsTargetedFilesystemAsRoot() {
	fs := s.newFilesystem()
	_ = fs.Mount(s.rootFs, "/path")
	s.adapterMock.
		On("PivotRoot", mock.Anything, mock.Anything).
		Return(nil)

	newRoot, _ := s.rootFs.Pivot(fs, "/old-root")

	s.EqualValues("/", newRoot.Path())
}

func (s *FsTestSuite) TestPivotSwapsParentChildRelationship() {
	fs := s.newFilesystem()
	s.adapterMock.
		On("PivotRoot", mock.Anything, mock.Anything).
		Return(nil)

	_, _ = s.rootFs.Pivot(fs, "/old-root")

	s.EqualValues("/old-root", s.rootFs.Path())
}

func (s *FsTestSuite) TestPivotPassesPreOperationPathToAdapter() {
	fs := s.newFilesystem()
	_ = s.rootFs.Mount(fs, "/new-root")
	s.adapterMock.
		On("PivotRoot", "/new-root", "/old-root").
		Return(nil)

	_, _ = s.rootFs.Pivot(fs, "/old-root")

	s.adapterMock.AssertExpectations(s.T())
}

func (s *FsTestSuite) TestUnmountPassesMountPointToAdapter() {
	fs := s.newFilesystemAtPath("/mnt")
	s.adapterMock.
		On("Unmount", "/mnt", mock.Anything).
		Return(nil)

	_ = fs.Unmount()

	s.adapterMock.AssertExpectations(s.T())
}

func (s *FsTestSuite) TestUnmountSetsNoFlags() {
	s.adapterMock.
		On("Unmount", mock.Anything, 0).
		Return(nil)

	_ = s.rootFs.Unmount()

	s.adapterMock.AssertExpectations(s.T())
}

func (s *FsTestSuite) TestUnmountForwardsErrorOnFailure() {
	s.adapterMock.
		On("Unmount", mock.Anything, mock.Anything).
		Return(errors.New("unmount error"))

	err := s.rootFs.Unmount()

	s.Error(err)
}

func (s FsTestSuite) newFilesystem() Filesystem {
	return s.newFilesystemAtPath("")
}

func (s FsTestSuite) newFilesystemAtPath(p path.Path) Filesystem {
	return Filesystem(&fs{path: p})
}
