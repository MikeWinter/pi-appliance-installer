package tmp

import (
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/fmt"
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/path"
	"io/ioutil"
	"os"
)

type Directory path.Path

func NewDirectory(parent path.Path, fs fmt.Filesystem) (Directory, error) {
	base, err := createParentIfAbsent(parent, fs.Path())
	if err != nil {
		return "", err
	}

	name, err := ioutil.TempDir(base.String(), "tmp")
	if err != nil {
		return "", err
	}

	return Directory(name), nil
}

func createParentIfAbsent(parent path.Path, root path.Path) (path.Path, error) {
	p := root.Join(parent)
	if err := os.MkdirAll(p.String(), os.ModePerm); err != nil && !os.IsExist(err) {
		return "", err
	}
	return p, nil
	//os.MkdirAll(p.String(), os.ModePerm)
	//segments, err := parent.Split()
	//if err != nil {
	//	return "", err
	//}
	//
	//base := root
	//for _, segment := range segments {
	//	base = base.Join(path.Path(segment))
	//	if err := os.Mkdir(base.String(), os.ModePerm); err != nil && !os.IsExist(err) {
	//		return "", err
	//	}
	//}
	//return base, nil
}

func (d Directory) Delete() error {
	return os.RemoveAll(d.String())
}

func (d Directory) String() string {
	return string(d)
}
