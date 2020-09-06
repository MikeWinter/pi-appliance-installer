package main

import (
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/dev"
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/fmt"
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/path"
	"github.com/MikeWinter/pi-on-boot-provisioning/fs/tmp"
	"github.com/bobziuchkovski/cue"
	"github.com/bobziuchkovski/cue/collector"
)

var log = cue.NewLogger("main")

func init() {
	cue.Collect(cue.INFO, collector.Terminal{
		ErrorsToStderr: true,
	}.New())
	cue.Collect(cue.INFO, collector.File{
		Path: "/boot/provisioning/log",
	}.New())
}

func main() {
	// TODO: Mount current root (/boot) partition as read-write
	// TODO: Create (and delete) temporary directories for /tmp and /mnt
	// TODO: Create device to represent Raspbian OS
	// TODO: Mount Raspbian OS partition
	// TODO: Replace current root with Raspbian OS

	bootFs := fmt.NewVfat().AsRoot()
	if err := bootFs.Remount(); err != nil {
		panic(err)
	}

	mountPoint, err := tmp.NewDirectory("/", bootFs)
	if err != nil {
		panic(err)
	}
	defer ignoringError(mountPoint.Delete)

	tmpFs := fmt.NewTmp()
	if err := bootFs.Mount(tmpFs, path.Path(mountPoint)); err != nil {
		panic(err)
	}
	defer ignoringError(tmpFs.Unmount)

	blockDevice, err := dev.NewSdCard(0, 2, "/tmp/mmcblk0p2")
	if err != nil {
		panic(err)
	}
	ext4Fs := fmt.NewExt4(blockDevice)
	if err := bootFs.Mount(ext4Fs, "/mnt"); err != nil {
		panic(err)
	}

	if _, err := bootFs.Pivot(ext4Fs, "/boot"); err != nil {
		panic(err)
	}
}

func ignoringError(fn func() error) {
	_ = fn()
}
