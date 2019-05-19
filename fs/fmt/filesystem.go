package fmt

import (
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/os"
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/path"
)

type fs struct {
	parent Filesystem
	path   path.Path
}

type Filesystem interface {
	AsRoot() Filesystem
	Mount(child Filesystem, at path.Path) error
	Path() path.Path
	Pivot(to Filesystem, old path.Path) (Filesystem, error)
	Remount() error
	Unmount() error

	setPath(path path.Path)
	setParent(parent Filesystem)
	onMount(at path.Path) error
}

func (fs *fs) AsRoot() Filesystem {
	fs.parent = nil
	fs.path = "/"
	return Filesystem(fs)
}

func (fs *fs) Mount(child Filesystem, at path.Path) error {
	child.setPath(at)
	child.setParent(fs)
	return child.onMount(at)
}

func (fs fs) Path() path.Path {
	if fs.parent == nil {
		return fs.path
	}
	return fs.parent.Path().Join(fs.path)
}

func (fs *fs) Pivot(to Filesystem, old path.Path) (Filesystem, error) {
	if err := os.Adapter.PivotRoot(to.Path().String(), old.String()); err != nil {
		return nil, err
	}
	fs.parent = to
	fs.path = old
	return to.AsRoot(), nil
}

func (fs fs) Remount() error {
	return os.Adapter.Mount("", fs.Path().String(), "", uintptr(REMOUNT), "")
}

func (fs fs) Unmount() error {
	return os.Adapter.Unmount(fs.Path().String(), 0)
}

func (fs *fs) setPath(path path.Path) {
	fs.path = path
}

func (fs *fs) setParent(parent Filesystem) {
	fs.parent = parent
}

func (fs fs) onMount(path path.Path) error {
	return nil
}
