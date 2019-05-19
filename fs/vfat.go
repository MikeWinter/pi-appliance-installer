package fs

type vfat struct {
	Filesystem
}

func NewVfat() Filesystem {
	return &vfat{&fs{}}
}
