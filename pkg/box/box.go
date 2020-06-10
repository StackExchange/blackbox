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
	//
	Team        string // Name of the team (i.e. .blackbox-$TEAM)
	RepoBaseDir string // Base directory of the repo.
	ConfigDir   string // Path of the .blackbox config directory (Rel to RepoBseDir)
	//
	Admins   []string        // If non-empty, the list of admins.
	Files    []string        // If non-empty, the list of files.
	FilesSet map[string]bool // If non-nil, a set of Files.
	//
	Vcs     vcs.Vcs          // Interface access to the VCS.
	Crypter crypters.Crypter // Inteface access to GPG.
	//
	Umask int // umask to set when decrypting
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

// NewUninitialized creates a box when nothing exists.
// Only for use with the "init" subcommand.
func NewUninitialized(configdir, team string) *Box {
	bx := &Box{}
	bx.Team = team
	bx.ConfigDir = GenerateConfigDir(configdir, team)
	return bx
}

// NewForTestingInit creates a box in a bare environment, with no
// autodiscovery of VCS.
// Useful only in integration tests.
func NewForTestingInit(vcsname string) *Box {
	bx := &Box{}

	// Set up the vcs
	var vh vcs.Vcs
	var err error
	for _, v := range vcs.Catalog {
		if strings.ToLower(v.Name) == strings.ToLower(vcsname) {
			vh, err = v.New()
			if err != nil {
				return nil // No idea how that would happen.
			}
		}
	}
	bx.Vcs = vh

	return bx
}

/*

test_init:

* Generates .git in the current directory
* We don't know or use configdir, team, etc.

init:
* Assumes the git repo was created already.
* Finds VCS (search up the tree for .git, .hg, none)
* If not found assume VCS=None, configdir=$pwd
* Set repobase based on where .git was found.
* If .blackbox/blackbox-team found, error.
* Create .blackbox or .blackbox-$team in configdir

post-init commands:
* Assumes the git repo was created already.
* Assumes .blackbox or .blackbox-$team exists.
* Finds VCS & rebobase (search up the tree for .git, .hg, none)
* If VCS=None, search for repobase by looking for .blackbox.
* If .blackbox not found, error (needs init)

*/

// NewFromFlags creates a box using items from flags.
func NewFromFlags(c *cli.Context) *Box {
	bx := &Box{}

	bx.Umask = c.Int("umask")
	bx.Team = c.String("team")

	// repoBaseDir, configDir, err := findBaseAndConfigDir()
	// if err != nil {
	// 	logErr.Println(err)
	// 	return bx
	// }
	// bx.RepoBaseDir = repoBaseDir
	// bx.ConfigDir = configDir

	// Discover which kind of VCS is in use.
	bx.Vcs = vcs.DetermineVcs(bx.RepoBaseDir)

	// Pick a crypto backend (GnuPG, go-openpgp, etc.)
	bx.Crypter = crypters.SearchByName(c.String("crypto"))
	if bx.Crypter == nil {
		fmt.Printf("ERROR!  No CRYPTER found! Please set --crypto correctly or use the default\n")
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

func (bx *Box) getAdmins() ([]string, error) {
	// Memoized
	if len(bx.Admins) != 0 {
		return bx.Admins, nil
	}

	// TODO(tlim): Try the json file.

	// Try the legacy file:
	fn := filepath.Join(bx.ConfigDir, "blackbox-admins.txt")
	logErr.Printf("Admins file: %q", fn)
	a, err := bbutil.ReadFileLines(fn)
	if err != nil {
		return nil, fmt.Errorf("getAdmins can't load admins (%q): %v", fn, err)
	}
	return a, nil
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
