package main

import (
	"github.com/MikeWinter/pi-on-boot-provisioning/fs"
	"github.com/bobziuchkovski/cue"
	"github.com/bobziuchkovski/cue/collector"
	"io/ioutil"
	"os"
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

	bootFs := fs.NewVfat().AsRoot()
	if err := bootFs.Remount(); err != nil {
		panic(err)
	}

	tmpDir, er := ioutil.TempDir("/", "tmp")
	if er != nil {
		panic(er)
	}
	defer os.Remove(tmpDir)

	tmpFs := fs.NewTmp()
	if err := bootFs.Mount(tmpFs, fs.Path(tmpDir)); err != nil {
		panic(err)
	}
	defer ignoringError(tmpFs.Unmount)

	blockDevice, err := fs.NewSdCard(0, 2, "/tmp/mmcblk0p2")
	if err != nil {
		panic(err)
	}
	ext4Fs := fs.NewExt4(blockDevice)
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
