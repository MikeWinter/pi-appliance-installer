package dev

import (
	"errors"
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/os"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type SdCardTestSuite struct {
	suite.Suite

	adapterMock *os.OsAdapterMock
}

func TestSdTestSuite(t *testing.T) {
	suite.Run(t, new(SdCardTestSuite))
}

func (s *SdCardTestSuite) SetupTest() {
	s.adapterMock = new(os.OsAdapterMock)

	os.Adapter = s.adapterMock
}

func (s *SdCardTestSuite) TestNewSdCardReturnsNilErrorOnSuccess() {
	s.adapterMock.
		On("MkNode", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	dev, err := NewSdCard(0, 0, "")

	s.NotNil(dev)
	s.NoError(err)
}

func (s *SdCardTestSuite) TestNewSdCardForwardsErrorOnFailure() {
	s.adapterMock.
		On("MkNode", mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("mknod error"))

	dev, err := NewSdCard(0, 0, "")

	s.Nil(dev)
	s.EqualError(err, "mknod error")
}

func (s *SdCardTestSuite) TestNewSdCardCreatesDeviceNodeAtTheGivenPath() {
	s.adapterMock.
		On("MkNode", "/dev/sdcard0", mock.Anything, mock.Anything).
		Return(nil)

	_, err := NewSdCard(0, 0, "/dev/sdcard0")

	require.NoError(s.T(), err)
	s.adapterMock.AssertExpectations(s.T())
}

func (s *SdCardTestSuite) TestNewSdCardCreatesDeviceNodeWithReadWritePermissions() {
	readWritePermissions := func(mode uint32) bool {
		return Mode(mode)&(RW_USR|RW_GRP) == (RW_USR | RW_GRP)
	}
	s.adapterMock.
		On("MkNode", mock.Anything, mock.MatchedBy(readWritePermissions), mock.Anything).
		Return(nil)

	_, err := NewSdCard(0, 0, "/dev")

	require.NoError(s.T(), err)
	s.adapterMock.AssertExpectations(s.T())
}

func (s *SdCardTestSuite) TestNewSdCardCreatesDeviceNodeWithBlockMode() {
	blockMode := func(mode uint32) bool {
		return Mode(mode)&BLOCK == BLOCK
	}
	s.adapterMock.
		On("MkNode", mock.Anything, mock.MatchedBy(blockMode), mock.Anything).
		Return(nil)

	_, err := NewSdCard(0, 0, "/dev")

	require.NoError(s.T(), err)
	s.adapterMock.AssertExpectations(s.T())
}

func (s *SdCardTestSuite) TestNewSdCardCreatesDeviceNodeWithMmcDeviceType() {
	mmcDevice := func(device int) bool {
		return device&MMC == MMC
	}
	s.adapterMock.
		On("MkNode", mock.Anything, mock.Anything, mock.MatchedBy(mmcDevice)).
		Return(nil)

	_, err := NewSdCard(0, 0, "/dev")

	require.NoError(s.T(), err)
	s.adapterMock.AssertExpectations(s.T())
}

func (s *SdCardTestSuite) TestNewSdCardCreatesDeviceNodeWithCardIndex() {
	card := 2
	cardIndex := func(device int) bool {
		expectedIndex := 2 * PARTITIONS
		return device&expectedIndex == expectedIndex
	}
	s.adapterMock.
		On("MkNode", mock.Anything, mock.Anything, mock.MatchedBy(cardIndex)).
		Return(nil)

	_, err := NewSdCard(card, 6, "/dev")

	require.NoError(s.T(), err)
	s.adapterMock.AssertExpectations(s.T())
}

func (s *SdCardTestSuite) TestNewSdCardCreatesDeviceNodeWithPartitionIndex() {
	partition := 7
	partitionIndex := func(device int) bool {
		expectedIndex := 7
		return device&expectedIndex == expectedIndex
	}
	s.adapterMock.
		On("MkNode", mock.Anything, mock.Anything, mock.MatchedBy(partitionIndex)).
		Return(nil)

	_, err := NewSdCard(1, partition, "/dev")

	require.NoError(s.T(), err)
	s.adapterMock.AssertExpectations(s.T())
}

func (s *SdCardTestSuite) TestSdCardHasCreatedPath() {
	s.adapterMock.
		On("MkNode", mock.Anything, mock.Anything, mock.Anything).
		Return(nil)

	dev, err := NewSdCard(1, 0, "/dev/sd-card0")

	require.NoError(s.T(), err)
	s.EqualValues("/dev/sd-card0", dev.Path())
}
