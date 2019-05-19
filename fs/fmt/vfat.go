package fmt

type vfat struct {
	Filesystem
}

func NewVfat() Filesystem {
	return &vfat{&fs{}}
}
