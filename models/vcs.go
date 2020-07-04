package models

import "github.com/StackExchange/blackbox/v2/pkg/commitlater"

// Vcs is git/hg/etc.
type Vcs interface {
	// Name returns the plug-in's canonical name.
	Name() string
	// Discover returns true if we are a repo of this type; along with the Abs path to the repo root (or "" if we don't know).
	Discover() (bool, string)

	// SetFileTypeUnix informs the VCS that files should maintain unix-style line endings.
	SetFileTypeUnix(repobasedir string, files ...string) error
	// IgnoreAnywhere tells the VCS to ignore these files anywhere in the repo.
	IgnoreAnywhere(repobasedir string, files []string) error
	// IgnoreAnywhere tells the VCS to ignore these files, rooted in the base of the repo.
	IgnoreFiles(repobasedir string, files []string) error

	// CommitTitle sets the title of the next commit.
	CommitTitle(title string)
	// NeedsCommit queues up commits for later execution.
	NeedsCommit(message string, repobasedir string, names []string)
	// DebugCommits dumps a list of future commits.
	DebugCommits() commitlater.List
	// FlushCommits informs the VCS to do queued up commits.
	FlushCommits() error

	// TestingInitRepo initializes a repo of this type (for use by integration tests)
	TestingInitRepo() error
}
