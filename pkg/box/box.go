package box

// box implements the box model.

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	// "github.com/StackExchange/blackbox/vcs"
	"github.com/StackExchange/blackbox/vcs"
	"github.com/urfave/cli/v2"
)

// Box describes what we know about a box.
type Box struct {
	//
	RepoBaseDir string // Base directory of the repo.
	ConfigDir   string // Path to the .blackbox config directory.
	//
	Admins  []string // If non-empty, the list of admins.
	Files   []string // If non-empty, the list of files.
	Vcs     vcs.Vcs  // Handle for VCS access.
	VcsName string   // name of the VCS
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

// NewFromFlags creates a box using items from flags.
func NewFromFlags(c *cli.Context) *Box {
	bx := &Box{}

	repoBaseDir, configDir, err := findBaseAndConfigDir()
	if err != nil {
		logErr.Println(err)
		return bx
	}
	bx.RepoBaseDir = repoBaseDir
	bx.ConfigDir = configDir

	// Discover which kind of VCS is in use.
	var h vcs.Vcs
	for _, v := range vcs.Catalog {
		h, err = v.New()
		if err != nil {
			return nil
		}
		if h.Discover() {
			bx.Vcs = h
			bx.VcsName = v.Name
			break
		}
	}

	return bx
}

func dirExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err == nil {
		return stat.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func findBaseAndConfigDir() (repodir, configdir string, err error) {

	// If BLACKBOXDATA/BLACKBOX_CONFIGDIR is set, that is the config dir.
	d := os.Getenv("BLACKBOXDATA")
	c := os.Getenv("BLACKBOX_CONFIGDIR")
	r := os.Getenv("BLACKBOX_REPOBASEDIR")
	// If any of those are used, r must be set and one or both of d & c
	// must be set. d is used before c.
	if d != "" {
		logErr.Printf("BLACKBOXDATA deprecated. Please use BLACKBOX_CONFIGDIR")
	}
	if (d != "") || (c != "") || (r != "") {
		if (d != "") && (r != "") {
			return r, d, nil
		}
		if (c != "") && (r != "") {
			return r, c, nil
		}
		return c, r, fmt.Errorf("if BLACKBOX_REPOBASEDIR or BLACKBOX_REPOBASEDIR is used, BLACKBOX_REPOBASEDIR must be set")
	}

	// Otherwise, search up the tree for the config dir.

	candidates := []string{}
	if team := os.Getenv("BLACKBOX_TEAM"); team != "" {
		candidates = append([]string{".blackbox-" + team}, candidates...)
	}
	candidates = append(candidates, ".blackbox")
	candidates = append(candidates, "keyrings/live")

	// Prevent an infinite loop by only doing "cd .." this many times
	maxDirLevels := 100

	relpath := ""
	for i := 0; i < maxDirLevels; i++ {
		// Does relpath contain any of our directory names?
		for _, c := range candidates {
			t := filepath.Join(relpath, c)
			d, err := dirExists(t)
			if err != nil {
				return "", "", fmt.Errorf("dirExists(%q) failed: %v", t, err)
			}
			if d {
				return relpath, t, nil
			}
		}
		// If we are at the root, stop.
		if abs, _ := filepath.Abs(relpath); abs == "/" {
			break
		}
		// Try one directory up
		relpath = filepath.Join("..", relpath)
	}

	return "", "", fmt.Errorf("No .blackbox directory found in cwd or above")
}

func (bx *Box) getAdmins() ([]string, error) {
	if len(bx.Admins) != 0 {
		return bx.Admins, nil
	}

	// TODO(tlim): Try the json file.

	// Try the legacy file:
	fn := filepath.Join(bx.ConfigDir, "blackbox-admins.txt")
	b, err := ioutil.ReadFile(fn)
	c := strings.TrimSpace(string(b))
	if err == nil {
		bx.Admins = strings.Split(c, "\n")
		return bx.Admins, nil
	}
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("getAdmins can't open %q: %v", fn, err)
	}

	return nil, fmt.Errorf("getAdmins can't load admin list")
}

func (bx *Box) getFiles() ([]string, error) {
	if len(bx.Files) != 0 {
		return bx.Files, nil
	}

	// TODO(tlim): Try the json file.

	// Try the legacy file:
	fn := filepath.Join(bx.ConfigDir, "blackbox-files.txt")
	b, err := ioutil.ReadFile(fn)
	c := strings.TrimSpace(string(b))
	if err == nil {
		bx.Files = strings.Split(c, "\n")
		return bx.Files, nil
	}
	if !os.IsNotExist(err) {
		return nil, fmt.Errorf("getFiles can't open %q: %v", fn, err)
	}

	return nil, fmt.Errorf("getFiles can't load file list")
}
