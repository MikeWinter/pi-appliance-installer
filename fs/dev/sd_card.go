package dev

import (
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/os"
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/path"
)

const (
	PARTITIONS int = 8
)

type sdCard struct {
	path path.Path
}

func NewSdCard(card int, partition int, path path.Path) (Device, error) {
	device := (card * PARTITIONS) | partition
	if err := os.Adapter.MkNode(path.String(), uint32(BLOCK|RW_USR|RW_GRP), MMC|device); err != nil {
		return nil, err
	}
	return &sdCard{path}, nil
}

func (sd sdCard) Path() path.Path {
	return sd.path
}
