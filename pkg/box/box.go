package box

// box implements the box model.

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/StackExchange/blackbox/v2/pkg/bblog"
	"github.com/StackExchange/blackbox/v2/pkg/bbutil"
	"github.com/StackExchange/blackbox/v2/pkg/crypters"
	"github.com/StackExchange/blackbox/v2/pkg/vcs"
	"github.com/urfave/cli/v2"
)

var logErr *log.Logger
var logDebug *log.Logger

// Box describes what we know about a box.
type Box struct {
	// Paths:
	Team        string // Name of the team (i.e. .blackbox-$TEAM)
	RepoBaseDir string // Rel path to the VCS repo.
	ConfigPath  string // Abs or Rel path to the .blackbox (or whatever) directory.
	ConfigRO    bool   // True if we should not try to change files in ConfigPath.
	// Settings:
	Umask  int    // umask to set when decrypting
	Editor string // Editor to call
	Debug  bool   // Are we in debug logging mode?
	// Cache of data gathered from .blackbox:
	Admins   []string        // If non-empty, the list of admins.
	Files    []string        // If non-empty, the list of files.
	FilesSet map[string]bool // If non-nil, a set of Files.
	// Handles to interfaces:
	Vcs      vcs.Vcs          // Interface access to the VCS.
	Crypter  crypters.Crypter // Inteface access to GPG.
	logErr   *log.Logger
	logDebug *log.Logger
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

// NewFromFlags creates a box using items from flags.  Nearly all subcommands use this.
func NewFromFlags(c *cli.Context) *Box {

	// The goal of this is to create a fully-populated box (and box.Vcs)
	// so that all subcommands have all the fields and interfaces they need
	// to do their job.

	logErr = bblog.GetErr()
	logDebug = bblog.GetDebug(c.Bool("debug"))

	bx := &Box{
		Umask:    c.Int("umask"),
		Editor:   c.String("editor"),
		Team:     c.String("team"),
		logErr:   bblog.GetErr(),
		logDebug: bblog.GetDebug(c.Bool("debug")),
		Debug:    c.Bool("debug"),
	}

	// Discover which kind of VCS is in use, and the repo root.
	bx.Vcs, bx.RepoBaseDir = vcs.Discover()

	// Discover the crypto backend (GnuPG, go-openpgp, etc.)
	bx.Crypter = crypters.SearchByName(c.String("crypto"), c.Bool("debug"))
	if bx.Crypter == nil {
		fmt.Printf("ERROR!  No CRYPTER found! Please set --crypto correctly or use the damn default\n")
		os.Exit(1)
	}

	// Find the .blackbox (or equiv.) directory.
	var err error
	configFlag := c.String("config")
	if configFlag != "" {
		// Flag is set. Better make sure it is valid.
		if !filepath.IsAbs(configFlag) {
			fmt.Printf("config flag value is a relative path. Too risky. Exiting.\n")
			os.Exit(1)
			// NB(tlim): We could return filepath.Abs(config) or maybe it just
			// works as is. I don't know, and until we have a use case to prove
			// it out, it's best to just not implement this.
		}
		bx.ConfigPath = configFlag
		bx.ConfigRO = true // External configs treated as read-only.
		// TODO(tlim): We could get fancy here and set ConfigReadOnly=true only
		// if we are sure configFlag is not within bx.RepoBaseDir. Again, I'd
		// like to see a use-case before we implement this.
		return bx

	}
	// Normal path. Flag not set, so we discover the path.
	bx.ConfigPath, err = FindConfigDir(bx.RepoBaseDir, c.String("team"))
	if err != nil && c.Command.Name != "info" {
		fmt.Printf("Can't find .blackbox or equiv. Have you run init?\n")
		os.Exit(1)
	}
	return bx
}

// NewUninitialized creates a box in a pre-init situation.
func NewUninitialized(c *cli.Context) *Box {
	/*
		This is for "blackbox init" (used before ".blackbox*" exists)

		Init needs:       How we populate it:
		bx.Vcs:           Discovered by calling each plug-in until succeeds.
		bx.ConfigDir:     Generated algorithmically (it doesn't exist yet).
	*/
	bx := &Box{
		Umask:    c.Int("umask"),
		Editor:   c.String("editor"),
		Team:     c.String("team"),
		logErr:   bblog.GetErr(),
		logDebug: bblog.GetDebug(c.Bool("debug")),
		Debug:    c.Bool("debug"),
	}
	bx.Vcs, bx.RepoBaseDir = vcs.Discover()
	if c.String("configdir") == "" {
		rel := ".blackbox"
		if bx.Team != "" {
			rel = ".blackbox-" + bx.Team
		}
		bx.ConfigPath = filepath.Join(bx.RepoBaseDir, rel)
	} else {
		// Wait. The user is using the --config flag on a repo that
		// hasn't been created yet?  I hope this works!
		fmt.Printf("ERROR: You can not set --config when initializing a new repo.  Please run this command from within a repo, with no --config flag.  Or, file a bug explaining your use caseyour use-case. Exiting!\n")
		os.Exit(1)
		// TODO(tlim): We could get fancy here and query the Vcs to see if the
		// path would fall within the repo, figure out the relative path, and
		// use that value. (and error if configflag is not within the repo).
		// That would be error prone and would only help the zero users that
		// ever see the above error message.
	}
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

func (bx *Box) getAdmins() error {
	// Memoized
	if len(bx.Admins) != 0 {
		return nil
	}

	// TODO(tlim): Try the json file.

	// Try the legacy file:
	fn := filepath.Join(bx.ConfigPath, "blackbox-admins.txt")
	bx.logDebug.Printf("Admins file: %q", fn)
	a, err := bbutil.ReadFileLines(fn)
	if err != nil {
		return fmt.Errorf("getAdmins can't load %q: %v", fn, err)
	}
	if !sort.StringsAreSorted(a) {
		return fmt.Errorf("file corrupt. Lines not sorted: %v", fn)
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
	fn := filepath.Join(bx.ConfigPath, "blackbox-files.txt")
	bx.logDebug.Printf("Files file: %q", fn)
	a, err := bbutil.ReadFileLines(fn)
	if err != nil {
		return fmt.Errorf("getFiles can't load %q: %v", fn, err)
	}
	if !sort.StringsAreSorted(a) {
		return fmt.Errorf("file corrupt. Lines not sorted: %v", fn)
	}
	for _, n := range a {
		bx.Files = append(bx.Files, filepath.Join(bx.RepoBaseDir, n))
	}

	bx.FilesSet = make(map[string]bool, len(bx.Files))
	for _, s := range bx.Files {
		bx.FilesSet[s] = true
	}

	return nil
}
