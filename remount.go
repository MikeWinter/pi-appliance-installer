package main

import (
	"golang.org/x/sys/unix"
	"os"
	"path/filepath"
	"syscall"
)

const (
	ReadWriteUserPermissions      = syscall.S_IRUSR | syscall.S_IWUSR
	ReadWriteGroupPermissions     = syscall.S_IRGRP | syscall.S_IWGRP
	BootPartitionPath             = "/"
	BootPartitionType             = "vfat"
	BootPartitionFlags            = syscall.MS_REMOUNT
	RootPartitionName             = "mmcblk0p2"
	RootPartitionType             = "ext4"
	RootPartitionFlags            = 0
	TemporaryPartitionPermissions = ReadWriteUserPermissions | ReadWriteGroupPermissions
	TemporaryPartitionPath        = "/tmp"
	TemporaryPartitionType        = "tmpfs"
	TemporaryPartitionFlags       = 0
	TemporaryPartitionDeviceMode  = TemporaryPartitionPermissions | syscall.S_IFBLK
	MountingPath                  = "/mnt"
	BootPath                      = "/boot"
	NoData                        = ""
	NoSource                      = ""
)

func fatally(err error) {
	if err != nil {
		panic(err)
	}
}

func ensureBootPartitionIsWritable() {
	fatally(unix.Mount(BootPartitionPath, BootPartitionPath, BootPartitionType, BootPartitionFlags, NoData))
}

func mountRootPartition() {
	rootPath := "/"
	directories := [...]string{TemporaryPartitionPath, MountingPath}
	for _, directory := range directories {
		// Only delete directories if they were explicitly created.
		if missing(directory) {
			fatally(createDirectory(directory))
			//noinspection GoDeferInLoop
			defer func(root *string, path string) {
				deleteDirectory(filepath.Join(*root, path))
			}(&rootPath, directory)
		}
	}

	fatally(mountTemporaryDirectory())
	defer func(root *string) {
		unmountDirectory(filepath.Join(*root, TemporaryPartitionPath))
	}(&rootPath)
	fatally(createDeviceNode())
	fatally(mountBlockDevice())
	fatally(exchangeRoot())
	rootPath = "/boot"
}

func missing(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return true
	} else if err != nil {
		panic(err)
	}
	return false
}

func createDirectory(path string) error {
	return os.Mkdir(path, os.FileMode(0755))
}

func deleteDirectory(path string) {
	fatally(os.Remove(path))
}

func mountTemporaryDirectory() error {
	return unix.Mount(NoSource, TemporaryPartitionPath, TemporaryPartitionType, TemporaryPartitionFlags, NoData)
}

func unmountDirectory(path string) {
	fatally(unix.Unmount(path, 0))
}

func createDeviceNode() error {
	path := filepath.Join(TemporaryPartitionPath, RootPartitionName)
	return unix.Mknod(path, TemporaryPartitionDeviceMode, makeDevice(179, 2))
}

func makeDevice(major int, minor int) int {
	return major << 8 | minor
}

func mountBlockDevice() error {
	path := filepath.Join(TemporaryPartitionPath, RootPartitionName)
	return unix.Mount(path, MountingPath, RootPartitionType, RootPartitionFlags, NoData)
}

func exchangeRoot() error {
	rootPath := MountingPath
	relativeBootPath := filepath.Join(MountingPath, BootPath)
	return unix.PivotRoot(rootPath, relativeBootPath)
}
