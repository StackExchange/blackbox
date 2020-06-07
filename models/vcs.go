package models

// Vcs is git/hg/etc.
type Vcs interface {
	Discover() bool
}
