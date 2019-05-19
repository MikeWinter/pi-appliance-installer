package fs

const (
	MMC int = 179 << 8
)

type Device interface {
	Path() Path
}
