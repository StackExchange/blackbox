package none

import (
	"github.com/StackExchange/blackbox/v2/pkg/vcs"
)

var pluginName = "GIT"

func init() {
	vcs.Register(pluginName, 0, newNone)
}

// VcsHandle is
type VcsHandle struct {
	Age int
}

func newNone() (vcs.Vcs, error) {
	return &VcsHandle{}, nil
}

// Name returns my name.
func (v VcsHandle) Name() string {
	return pluginName
}

// Discover returns true
func (v VcsHandle) Discover(repobasedir string) bool {
	return true
}

// The following are "secret" functions only used by the integration testing system.

// TestingInitRepo initializes a repo.
func (v VcsHandle) TestingInitRepo() error {
	return nil
}
