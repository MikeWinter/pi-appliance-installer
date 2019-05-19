package fmt

import (
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/os"
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/path"
)

type tmp struct {
	Filesystem
}

func NewTmp() Filesystem {
	return &tmp{&fs{}}
}

func (fs *tmp) onMount(path path.Path) error {
	if err := fs.Filesystem.onMount(path); err != nil {
		return err
	}
	return os.Adapter.Mount("", fs.Path().String(), "tmpfs", 0, "")
}
