package fs

type fs struct {
	parent Filesystem
	path   Path
}

type Filesystem interface {
	AsRoot() Filesystem
	Mount(child Filesystem, at Path) error
	Path() Path
	Pivot(to Filesystem, old Path) (Filesystem, error)
	Remount() error
	Unmount() error

	setPath(path Path)
	setParent(parent Filesystem)
	onMount(at Path) error
}

func (fs *fs) AsRoot() Filesystem {
	fs.parent = nil
	fs.path = "/"
	return Filesystem(fs)
}

func (fs *fs) Mount(child Filesystem, at Path) error {
	child.setPath(at)
	child.setParent(fs)
	return child.onMount(at)
}

func (fs fs) Path() Path {
	if fs.parent == nil {
		return fs.path
	}
	return fs.parent.Path().Join(fs.path)
}

func (fs *fs) Pivot(to Filesystem, old Path) (Filesystem, error) {
	if err := adapter.PivotRoot(to.Path().String(), old.String()); err != nil {
		return nil, err
	}
	fs.parent = to
	fs.path = old
	return to.AsRoot(), nil
}

func (fs fs) Remount() error {
	return adapter.Mount("", fs.Path().String(), "", uintptr(REMOUNT), "")
}

func (fs fs) Unmount() error {
	return adapter.Unmount(fs.Path().String(), 0)
}

func (fs *fs) setPath(path Path) {
	fs.path = path
}

func (fs *fs) setParent(parent Filesystem) {
	fs.parent = parent
}

func (fs fs) onMount(path Path) error {
	return nil
}
