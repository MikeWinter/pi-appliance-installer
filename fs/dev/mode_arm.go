package dev

import (
	"golang.org/x/sys/unix"
)

const (
	BLOCK  Mode = unix.S_IFBLK
	RW_USR      = unix.S_IRUSR | unix.S_IWUSR
	RW_GRP      = unix.S_IRGRP | unix.S_IWGRP
)
