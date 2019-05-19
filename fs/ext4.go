package fs

type ext4 struct {
	Filesystem
	device Device
}

func NewExt4(device Device) Filesystem {
	return &ext4{&fs{}, device}
}

func (fs *ext4) onMount(path Path) error {
	if err := fs.Filesystem.onMount(path); err != nil {
		return err
	}
	return adapter.Mount(fs.device.Path().String(), fs.Path().String(), "ext4", 0, "")
}
