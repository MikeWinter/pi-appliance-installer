package tmp

import (
	str "fmt"
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/fmt"
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/path"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"testing"
)

type DirTestSuite struct {
	suite.Suite

	fs  *fmt.FilesystemMock
	tmp path.Path
}

func TestDirTestSuite(t *testing.T) {
	suite.Run(t, new(DirTestSuite))
}

func (s *DirTestSuite) SetupTest() {
	s.fs = new(fmt.FilesystemMock)

	tmp, err := ioutil.TempDir("", "DirTest-")
	if err != nil {
		s.FailNow("could not create test root directory", err)
	}
	s.tmp = path.Path(tmp)
	s.fs.On("Path").
		Return(s.tmp)
}

func (s *DirTestSuite) TearDownTest() {
	if err := os.RemoveAll(s.tmp.String()); err != nil {
		s.FailNow("failed to clean up test root directory", err)
	}
}

func (s *DirTestSuite) TestNewDirectoryIsDirectory() {
	dir, err := NewDirectory("/path", s.fs)

	s.NoError(err)
	s.DirExists(dir.String())
}

func (s *DirTestSuite) TestNewDirectoryCanBeDeeplyNested() {
	dir, err := NewDirectory("/path/with/many/parts", s.fs)

	s.NoError(err)
	s.DirExists(dir.String())
}

func (s *DirTestSuite) TestNewDirectoryIsRootedInFilesystem() {
	dir, _ := NewDirectory("/path", s.fs)

	s.Regexp(str.Sprintf("^%s/.+", s.tmp.Join("/path/")), dir)
}

func (s *DirTestSuite) TestNewDirectoriesHaveDifferentPaths() {
	first, _ := NewDirectory("/", s.fs)
	second, _ := NewDirectory("/", s.fs)

	s.NotEqual(first, second)
}

func (s *DirTestSuite) TestDeleteRemovesEmptyDirectory() {
	dir, _ := NewDirectory("/", s.fs)

	err := dir.Delete()

	s.NoError(err)
	_, statErr := os.Lstat(dir.String())
	s.True(os.IsNotExist(statErr))
}

func (s *DirTestSuite) TestDeleteRemovesDirectorySubtree() {
	dir, _ := NewDirectory("/", s.fs)
	if _, err := os.Create(path.Path(dir).Join("file").String()); err != nil {
		s.FailNow("unable to create file at " + dir.String())
	}

	err := dir.Delete()

	s.NoError(err)
	_, statErr := os.Lstat(dir.String())
	s.True(os.IsNotExist(statErr))
}
