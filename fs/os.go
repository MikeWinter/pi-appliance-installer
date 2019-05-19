package fs

var adapter osAdapter = noOpAdapter{}

type osAdapter interface {
	Mount(source string, target string, fsType string, flags uintptr, data string) error
	MkNode(path string, mode uint32, dev int) error
	PivotRoot(newRoot string, putOld string) error
	Unmount(target string, flags int) error
}

type noOpAdapter struct{}

func (noOpAdapter) Mount(source string, target string, fsType string, flags uintptr, data string) error {
	panic("Mount is not supported")
}

func (noOpAdapter) MkNode(path string, mode uint32, dev int) error {
	panic("MkNode is not supported")
}

func (noOpAdapter) PivotRoot(newRoot string, putOld string) error {
	panic("PivotRoot is not supported")
}

func (noOpAdapter) Unmount(target string, flags int) error {
	panic("Unmount is not supported")
}
