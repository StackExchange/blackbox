package none

import (
	"github.com/StackExchange/blackbox/v2/pkg/vcs"
)

func init() {
	vcs.Register("NONE", 0, newNone)
}

// VcsHandle is
type VcsHandle struct {
	Age int
}

func newNone() (vcs.Vcs, error) {
	return &VcsHandle{}, nil
}

// Discover returns true
func (v VcsHandle) Discover(repobasedir string) bool {
	return true
}
