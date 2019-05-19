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
	for _, p := range other {
		parts = append(parts, p.String())
	}
	return Path(filepath.Join(parts...))
}
