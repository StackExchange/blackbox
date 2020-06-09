package git

import (
	"path/filepath"

	"github.com/StackExchange/blackbox/v2/pkg/bbutil"
	"github.com/StackExchange/blackbox/v2/pkg/vcs"
)

var pluginName = "GIT"

func init() {
	vcs.Register(pluginName, 100, newGit)
}

// VcsHandle is the handle
type VcsHandle struct {
}

func newGit() (vcs.Vcs, error) {
	return &VcsHandle{}, nil
}

// Name returns my name.
func (v VcsHandle) Name() string {
	return pluginName
}

// Discover returns false.
func (v VcsHandle) Discover(repobasedir string) bool {
	n := filepath.Join(repobasedir, ".git")
	found, err := bbutil.DirExists(n)
	if err != nil {
		return false
	}
	return found
}

// The following are "secret" functions only used by the integration testing system.

// TestingInitRepo initializes a repo.

func (v VcsHandle) TestingInitRepo() error {
	bbutil.RunBash("git", "init")
	return nil
}
