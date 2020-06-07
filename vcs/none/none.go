package none

import (
	"github.com/StackExchange/blackbox/vcs"
)

func init() {
	vcs.Register("NONE", 0, newNone)
}

// VcsHandle is
type VcsHandle struct {
	Age int
}

func newNone() (*VcsHandle, error) {
	return &VcsHandle{}, nil
}

// Discover returns false
func (v VcsHandle) Discover() bool {
	return false
}
