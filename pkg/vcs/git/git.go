package git

import (
	"path/filepath"

	"github.com/StackExchange/blackbox/v2/pkg/bbutil"
	"github.com/StackExchange/blackbox/v2/pkg/vcs"
)

func init() {
	vcs.Register("GIT", 100, newGit)
}

// VcsHandle is the handle
type VcsHandle struct {
}

func newGit() (vcs.Vcs, error) {
	return &VcsHandle{}, nil
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
