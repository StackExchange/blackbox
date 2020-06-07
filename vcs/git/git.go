package git

import (
	"github.com/StackExchange/blackbox/vcs"
)

func init() {
	vcs.Register("GIT", 1, newGit)
}

// VcsHandle is the handle
type VcsHandle struct {
}

func newGit() (*VcsHandle, error) {
	return &VcsHandle{}, nil
}

// Discover returns false.
func (v VcsHandle) Discover() bool {
	return false
}
