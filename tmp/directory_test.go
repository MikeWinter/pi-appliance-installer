package tmp

import (
	"github.com/stretchr/testify/suite"
	"os"
	"testing"
)

type TmpDirectoryTestSuite struct {
	suite.Suite
}

func TestTmpDirectoryTestSuite(t *testing.T) {
	suite.Run(t, new(TmpDirectoryTestSuite))
}

func (s *TmpDirectoryTestSuite) TestNewDirectoryHasPath() {
	parentDir := os.TempDir()

	dir, err := NewDirectory(parentDir)

	s.NoError(err)
	s.Contains(dir.Path(), parentDir)
}

func (s *TmpDirectoryTestSuite) TestNewDirectoryIsDirectory() {
	dir, _ := NewDirectory(os.TempDir())

	s.DirExists(dir.Path())
}

func (s *TmpDirectoryTestSuite) TestNewDirectoryIsNotTheSameAsTheParentDirectory() {
	parentDir := os.TempDir()

	dir, _ := NewDirectory(parentDir)

	s.NotEqual(dir.Path(), parentDir)
}

func (s *TmpDirectoryTestSuite) TestNewDirectoriesAreNotTheSame() {
	dir1, _ := NewDirectory(os.TempDir())
	dir2, _ := NewDirectory(os.TempDir())

	s.NotEqual(dir1.Path(), dir2.Path())
}

func (s *TmpDirectoryTestSuite) TestDeleteRemovesDirectory() {
	dir, _ := NewDirectory(os.TempDir())

	err := dir.Delete()

	s.NoError(err)
	_, statErr := os.Stat(dir.Path())
	s.True(os.IsNotExist(statErr))
}
