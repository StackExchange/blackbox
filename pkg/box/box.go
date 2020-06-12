package box

// box implements the box model.

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/StackExchange/blackbox/v2/pkg/bbutil"
	"github.com/StackExchange/blackbox/v2/pkg/crypters"
	"github.com/StackExchange/blackbox/v2/pkg/vcs"
	"github.com/urfave/cli/v2"
)

// Box describes what we know about a box.
type Box struct {
	// Paths:
	Team         string // Name of the team (i.e. .blackbox-$TEAM) TODO(tlim): Can this be deleted?
	RepoBaseDir  string // Abs path to the VCS repo.
	ConfigDir    string // Abs path to the .blackbox (or whatever) directory.
	ConfigDirRel string // Path to the .blackbox (or whatever) directory relative to RepoBaseDir
	// Settings:
	Umask int // umask to set when decrypting
	// Cache of data gathered from .blackbox:
	Admins   []string        // If non-empty, the list of admins.
	Files    []string        // If non-empty, the list of files.
	FilesSet map[string]bool // If non-nil, a set of Files.
	// Handles to interfaces:
	Vcs     vcs.Vcs          // Interface access to the VCS.
	Crypter crypters.Crypter // Inteface access to GPG.
}

// StatusMode is a type of query.
type StatusMode int

const (
	// Itemized is blah
	Itemized StatusMode = iota // Individual files by name
	// All files is blah
	All
	// Unchanged is blah
	Unchanged
	// Changed is blah
	Changed
)

var logErr *log.Logger

func init() {
	logErr = log.New(os.Stderr, "", 0)
}

// NewFromFlags creates a box using items from flags.  Nearly all subcommands use this.
func NewFromFlags(c *cli.Context) *Box {
	/*
				 Nearly all subcommands use this.  It is used with a VCS repo
				 that has blackbox already initialized.

				 Commands need:    How we populate it:
			     bx.Vcs:           Discovered by calling each plug-in until succeeds.
		       bx.ConfigDir:     Is discovered.
		       bx.RepoBaseDir:   Is discovered.
	*/
	bx := &Box{
		Umask: c.Int("umask"),
		Team:  c.String("team"),
	}

	// Discover which kind of VCS is in use.
	bx.Vcs = vcs.Discover(bx.RepoBaseDir)

	// Pick a crypto backend (GnuPG, go-openpgp, etc.)
	bx.Crypter = crypters.SearchByName(c.String("crypto"))
	if bx.Crypter == nil {
		fmt.Printf("ERROR!  No CRYPTER found! Please set --crypto correctly or use the damn default\n")
		os.Exit(1)
	}

	// Are we using .blackbox or what?
	var err error
	bx.ConfigDir, err = FindConfigDir(c.String("config"), c.String("team"))
	if err != nil {
		return nil
	}

	return bx
}

// NewUninitialized creates a box in a pre-init situation.
func NewUninitialized(configdir, team string) *Box {
	/*
		   This is for "blackbox init" (used before ".blackbox*" exists)

			 Init needs:       How we populate it:
			   bx.Vcs:           Discovered by calling each plug-in until succeeds.
			   bx.ConfigDir:     Generated algorithmically (it doesn't exist yet).
				 bx.RepoBaseDir:   Generated algorithmically (it doesn't exist yet).
	*/
	bx := &Box{
		Team: team,
	}
	bx.Vcs = vcs.Discover(bx.RepoBaseDir)
	bx.ConfigDir = GenerateConfigDir(configdir, team)
	return bx
}

// NewForTestingInit creates a box in a bare environment.
func NewForTestingInit(vcsname string) *Box {
	/*

		This is for "blackbox test_init" (secret command used in integration tests; when nothing exists)

		TestingInitRepo only uses bx.Vcs, so that's all we set.

		Populates bx.Vcs by finding the provider named vcsname.
	*/
	bx := &Box{}

	// Find the
	var vh vcs.Vcs
	var err error
	vcsname = strings.ToLower(vcsname)
	for _, v := range vcs.Catalog {
		if strings.ToLower(v.Name) == vcsname {
			vh, err = v.New()
			if err != nil {
				return nil // No idea how that would happen.
			}
		}
	}
	bx.Vcs = vh

	return bx
}

// func findBaseAndConfigDir() (repodir, configdir string, err error) {

// 	// Otherwise, search up the tree for the config dir.

// 	candidates := []string{}
// 	if team := os.Getenv("BLACKBOX_TEAM"); team != "" {
// 		candidates = append([]string{".blackbox-" + team}, candidates...)
// 	}
// 	candidates = append(candidates, ".blackbox")
// 	candidates = append(candidates, "keyrings/live")

// 	// Prevent an infinite loop by only doing "cd .." this many times
// 	maxDirLevels := 100

// 	relpath := ""
// 	for i := 0; i < maxDirLevels; i++ {
// 		// Does relpath contain any of our directory names?
// 		for _, c := range candidates {
// 			t := filepath.Join(relpath, c)
// 			d, err := bbutil.DirExists(t)
// 			if err != nil {
// 				return "", "", fmt.Errorf("dirExists(%q) failed: %v", t, err)
// 			}
// 			if d {
// 				return relpath, t, nil
// 			}
// 		}
// 		// If we are at the root, stop.
// 		if abs, _ := filepath.Abs(relpath); abs == "/" {
// 			break
// 		}
// 		// Try one directory up
// 		relpath = filepath.Join("..", relpath)
// 	}

// 	return "", "", fmt.Errorf("No .blackbox directory found in cwd or above")
// }

func (bx *Box) getAdmins() error {
	// Memoized
	if len(bx.Admins) != 0 {
		return nil
	}

	// TODO(tlim): Try the json file.

	// Try the legacy file:
	fn := filepath.Join(bx.ConfigDir, "blackbox-admins.txt")
	logErr.Printf("Admins file: %q", fn)
	a, err := bbutil.ReadFileLines(fn)
	if err != nil {
		return fmt.Errorf("getAdmins can't load admins (%q): %v", fn, err)
	}
	bx.Admins = a
	return nil
}

// getFiles populates Files and FileMap.
func (bx *Box) getFiles() error {
	if len(bx.Files) != 0 {
		return nil
	}

	// TODO(tlim): Try the json file.

	// Try the legacy file:
	fn := filepath.Join(bx.ConfigDir, "blackbox-files.txt")
	b, err := ioutil.ReadFile(fn)
	if err != nil {
		return fmt.Errorf("getFiles can't read %q: %v", fn, err)
	}

	c := strings.TrimSpace(string(b))

	bx.Files = strings.Split(c, "\n")
	bx.FilesSet = make(map[string]bool)
	for _, s := range bx.Files {
		bx.FilesSet[s] = true
	}

	return nil
}
