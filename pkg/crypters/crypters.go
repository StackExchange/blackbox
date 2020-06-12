package crypters

import (
	"sort"
	"strings"

	"github.com/StackExchange/blackbox/v2/models"
)

// Crypter is the handle
type Crypter interface {
	models.Crypter
}

// NewFnSig function signature needed by reg.
type NewFnSig func(debug bool) (Crypter, error)

// Item stores one item
type Item struct {
	Name     string
	New      NewFnSig
	Priority int
}

// Catalog is the list of registered vcs's.
var Catalog []*Item

// SearchByName returns a Crypter handle for name.
// The search is case insensitive.
func SearchByName(name string, debug bool) Crypter {
	name = strings.ToLower(name)
	for _, v := range Catalog {
		//fmt.Printf("Trying %v %v\n", v.Name)
		if strings.ToLower(v.Name) == name {
			chandle, err := v.New(debug)
			if err != nil {
				return nil // No idea how that would happen.
			}
			//fmt.Printf("USING! %v\n", v.Name)
			return chandle
		}
	}
	return nil
}

// Register a new VCS.
func Register(name string, priority int, newfn NewFnSig) {
	//fmt.Printf("CRYPTER registered: %v\n", name)
	item := &Item{
		Name:     name,
		New:      newfn,
		Priority: priority,
	}
	Catalog = append(Catalog, item)

	// Keep the list sorted.
	sort.Slice(Catalog, func(i, j int) bool { return Catalog[j].Priority < Catalog[i].Priority })
}
