package fs

type tmp struct {
	Filesystem
}

func NewTmp() Filesystem {
	return &tmp{&fs{}}
}

func (fs *tmp) onMount(path Path) error {
	if err := fs.Filesystem.onMount(path); err != nil {
		return err
	}
	return adapter.Mount("", fs.Path().String(), "tmpfs", 0, "")
}
