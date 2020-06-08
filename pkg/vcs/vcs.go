package vcs

import (
	"sort"

	"github.com/StackExchange/blackbox/v2/models"
)

// Vcs is the handle
type Vcs interface {
	models.Vcs
}

// NewFnSig function signature needed by reg.
type NewFnSig func() (Vcs, error)

// Item stores one item
type Item struct {
	Name     string
	New      NewFnSig
	Priority int
}

// Catalog is the list of registered vcs's.
var Catalog []*Item

// Register a new VCS.
func Register(name string, priority int, newfn NewFnSig) {
	//fmt.Printf("VCS registered: %v\n", name)
	item := &Item{
		Name:     name,
		New:      newfn,
		Priority: priority,
	}
	Catalog = append(Catalog, item)

	// Keep the list sorted.
	sort.Slice(Catalog, func(i, j int) bool { return Catalog[j].Priority < Catalog[i].Priority })
}