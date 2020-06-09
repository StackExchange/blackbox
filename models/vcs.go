package models

// Vcs is git/hg/etc.
type Vcs interface {
	// Name returns the plug-in's canonical name.
	Name() string
	// Discover returns true if the cwd is a VCS of this type.
	Discover(repobasedir string) bool
	// Initialize a repo of this type (for use by integration tests)
	TestingInitRepo() error
}
