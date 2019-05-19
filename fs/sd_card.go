package fs

const (
	PARTITIONS int = 8
)

type sdCard struct {
	path Path
}

func NewSdCard(card int, partition int, path Path) (Device, error) {
	device := (card * PARTITIONS) | partition
	if err := adapter.MkNode(path.String(), uint32(BLOCK|RW_USR|RW_GRP), MMC|device); err != nil {
		return nil, err
	}
	return &sdCard{path}, nil
}

func (sd sdCard) Path() Path {
	return sd.path
}
