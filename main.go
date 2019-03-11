package main

import (
	"errors"
	"fmt"
	"golang.org/x/sys/unix"
	"os"
	"path/filepath"
)

var logFile *os.File

func main() {
	defer abortStartupOnError()

	ensureBootPartitionIsWritable()
	mountRootPartition()
	fatally(createLogFile())

	configuration := loadConfiguration()

	defer startOs(configuration)
	defer closeLogFile()
	defer convertErrorToLogMessage()

	// copy each subdirectory in /boot/provisioner recursively to /
	fatally(provisionApps("/boot/provisioner"))
}

func createLogFile() (err error) {
	logFile, err = os.Create("/boot/provisioner/log")
	return
}

func closeLogFile() {
	_ = logFile.Close()
}

func abortStartupOnError() {
	if e := recover(); e != nil {
		log("an unrecoverable error occurred during boot: %v\n", e)
		os.Exit(1)
	}
}

func loadConfiguration() *Configuration {
	configuration, e := Load("/boot/provisioner/settings.conf")
	if e != nil {
		panic(e)
	}
	return configuration
}

func startOs(configuration *Configuration) {
	if configuration.FirstBoot {
		configuration.FirstBoot = false
		fatally(configuration.Save("/boot/provisioner/settings.conf"))
	}

	// Control should pass to the invoked executable and never return.
	_ = unix.Exec(configuration.OnBoot(), []string{"init"}, os.Environ())
	// If the call does return, it always represents an error.
	os.Exit(1)
}

func convertErrorToLogMessage() {
	if e := recover(); e != nil {
		log("an error occurred: %v\n", e)
	}
}

func log(format string, args ...interface{}) {
	message := fmt.Sprintf(filepath.Base(os.Args[0]) + ": " + format, args...)
	_, _ = fmt.Fprint(os.Stderr, message)
	_, _ = fmt.Fprint(logFile, message)
}

func stoppingOnError(walkFn func(path string, info os.FileInfo) error) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		return walkFn(path, info)
	}
}

func provisionApps(rootPath string) error {
	log("Provisioning apps...\n")
	return filepath.Walk(
		rootPath,
		stoppingOnError(func(path string, info os.FileInfo) error {
			if !info.IsDir() || path == rootPath {
				return nil
			}

			log("Copying app '%s'...\n", filepath.Base(path))
			if err := os.Chdir(path); err != nil {
				return err
			}
			if err := filepath.Walk(".", stoppingOnError(copyApp)); err != nil {
				log("An error occurred while copying: %v\n", err)
			}
			return filepath.SkipDir
		}))
}

func copyApp(path string, info os.FileInfo) error {
	absolutePath := filepath.Join("/", path)
	if info.IsDir() {
		return failingUnlessPathExists(createDirectory(absolutePath))
	}

	ensureFileDoesNotExist(absolutePath)
	return linkFile(path, absolutePath)
}

func failingUnlessPathExists(err error) error {
	if os.IsExist(err) {
		return nil
	}
	return err
}

func ensureFileDoesNotExist(path string) {
	if exists(path) && os.Remove(path) != nil {
		log("Unable to remove existing entry at %s\n", path)
	}
}

func linkFile(source string, destination string) error {
	if absoluteSourcePath, err := filepath.Abs(source); err != nil {
		return errors.New(fmt.Sprintf("Could not resolve '%s' to an absolute path (cause: %v)\n", source, err))
	} else if err := os.Symlink(absoluteSourcePath, destination); err != nil {
		return errors.New(fmt.Sprintf("Could not symlink '%s' to '%s' (cause: %v)\n", absoluteSourcePath, destination, err))
	} else {
		log("Linked '%s' to '%s'\n", absoluteSourcePath, destination)
	}
	return nil
}

func exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	} else if err != nil {
		panic(err)
	}
	return true
}
