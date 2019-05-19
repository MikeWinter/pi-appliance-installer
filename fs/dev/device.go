package dev

import (
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/path"
)

const (
	MMC int = 179 << 8
)

type Device interface {
	Path() path.Path
}
