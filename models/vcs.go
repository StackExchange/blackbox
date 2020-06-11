package models

// Vcs is git/hg/etc.
type Vcs interface {
	// Name returns the plug-in's canonical name.
	Name() string
	// Discover returns true if the cwd is a VCS of this type.
	Discover(repobasedir string) bool

	// SetFileTypeUnix informs the VCS that files should maintain unix-style line endings.
	SetFileTypeUnix(repobasedir string, files ...string) error
	// IgnoreAnywhere tells the VCS to ignore these files anywhere rin the repo.
	IgnoreAnywhere(repobasedir string, files ...string) error
	// SuggestTracking tells the VCS to suggest the user commit these files.
	SuggestTracking(repobasedir string, message string, files ...string) error

	// TestingInitRepo initializes a repo of this type (for use by integration tests)
	TestingInitRepo() error
}
