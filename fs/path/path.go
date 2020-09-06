package path

import (
	"path/filepath"
)

type Path string

func (p Path) String() string {
	return string(p)
}

func (p Path) Join(other ...Path) Path {
	parts := []string{p.String()}
	for _, path := range other {
		parts = append(parts, path.String())
	}
	return Path(filepath.Join(parts...))
}

//func (p Path) Split() ([]string, error) {
//	relative, err := filepath.Rel("/", p.String())
//	if err != nil {
//		return nil, err
//	}
//	return strings.Split(relative, "/"), nil
//}
