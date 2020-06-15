package none

import (
	"fmt"

	"github.com/StackExchange/blackbox/v2/pkg/vcs"
)

var pluginName = "NONE"

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

// SetFileTypeUnix informs the VCS that files should maintain unix-style line endings.
func (v VcsHandle) SetFileTypeUnix(repobasedir string, files ...string) error {
	return nil
}

// IgnoreAnywhere tells the VCS to ignore these files anywhere rin the repo.
func (v VcsHandle) IgnoreAnywhere(repobasedir string, files ...string) error {
	return nil
}

// SuggestTracking tells the VCS to suggest the user commit these files.
func (v VcsHandle) SuggestTracking(repobasedir string, message string, files []string) error {
	return nil
}

// The following are "secret" functions only used by the integration testing system.

// TestingInitRepo initializes a repo.
func (v VcsHandle) TestingInitRepo() error {
	fmt.Println("VCS=none, TestingInitRepo")
	return nil
}
