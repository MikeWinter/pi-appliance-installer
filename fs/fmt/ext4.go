package fmt

import (
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/dev"
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/os"
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/path"
)

type ext4 struct {
	Filesystem
	device dev.Device
}

func NewExt4(device dev.Device) Filesystem {
	return &ext4{&fs{}, device}
}

func (fs *ext4) onMount(path path.Path) error {
	if err := fs.Filesystem.onMount(path); err != nil {
		return err
	}
	return os.Adapter.Mount(fs.device.Path().String(), fs.Path().String(), "ext4", 0, "")
}
