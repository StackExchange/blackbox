package bbnone

import (
	"log"
	"os"
)

// NoneInfo contains Git-specific info about this repository.
type NoneInfo struct {
}

// New is a factory; returns nil if this is not a Git repo.
func New() (*NoneInfo, error) {
	return new(NoneInfo), nil
}

// Name returns the name of this type of repo.
func (repo *NoneInfo) Name() string {
	return "unknown"
}

// RepoBaseDir returns the current working directory.
func (repo *NoneInfo) RepoBaseDir() string {
	d, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	return d
}
