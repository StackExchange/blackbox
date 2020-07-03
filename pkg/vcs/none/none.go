package none

import (
	"fmt"

	"github.com/StackExchange/blackbox/v2/pkg/commitlater"
	"github.com/StackExchange/blackbox/v2/pkg/vcs"
)

var pluginName = "NONE"

func init() {
	vcs.Register(pluginName, 0, newNone)
}

// VcsHandle is
type VcsHandle struct {
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

// IgnoreAnywhere tells the VCS to ignore these files anywhere in the repo.
func (v VcsHandle) IgnoreAnywhere(repobasedir string, files []string) error {
	return nil
}

// IgnoreFiles tells the VCS to ignore these files anywhere in the repo.
func (v VcsHandle) IgnoreFiles(repobasedir string, files []string) error {
	return nil
}

// CommitTitle sets the title of the next commit.
func (v VcsHandle) CommitTitle(title string) {}

// NeedsCommit queues up commits for later execution.
func (v VcsHandle) NeedsCommit(message string, repobasedir string, names []string) {
	return
}

// DebugCommits dumps a list of future commits.
func (v VcsHandle) DebugCommits() commitlater.List {
	return commitlater.List{}
}

// FlushCommits informs the VCS to do queued up commits.
func (v VcsHandle) FlushCommits() error {
	return nil
}

// The following are "secret" functions only used by the integration testing system.

// TestingInitRepo initializes a repo.
func (v VcsHandle) TestingInitRepo() error {
	fmt.Println("VCS=none, TestingInitRepo")
	return nil
}
