package vcs

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

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

// Discover polls the VCS plug-ins to determine the VCS of directory.
// The first to succeed is returned.
// It never returns nil, since "NONE" is always valid.
func Discover() (Vcs, string) {
	for _, v := range Catalog {
		h, err := v.New()
		if err != nil {
			return nil, "" // No idea how that would happen.
		}
		if b, repodir := h.Discover(); b {

			// Try to find the rel path from CWD to RepoBase
			wd, err := os.Getwd()
			if err != nil {
				fmt.Printf("ERROR: Can not determine cwd! Failing!\n")
				os.Exit(1)
			}
			//fmt.Printf("DISCCOVER: WD=%q REPO=%q\n", wd, repodir)
			if repodir != wd && strings.HasSuffix(repodir, wd) {
				// This is a terrible hack.  We're basically guessing
				// at the filesystem layout.  That said, it works on macOS.
				// TODO(tlim): Abstract this out into a separate function
				// so we can do integration tests on it (to know if it fails on
				// a particular operating system.)
				repodir = wd
			}
			r, err := filepath.Rel(wd, repodir)
			if err != nil {
				// Wait, we're not relative to each other? Give up and
				// just return the abs repodir.
				return h, repodir
			}
			return h, r
		}
	}
	// This can't happen. If it does, we'll panic and that's ok.
	return nil, ""
}

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
