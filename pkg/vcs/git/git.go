package git

import (
	"fmt"
	"path/filepath"

	"github.com/StackExchange/blackbox/v2/pkg/bbutil"
	"github.com/StackExchange/blackbox/v2/pkg/vcs"
)

var pluginName = "GIT"

func init() {
	vcs.Register(pluginName, 100, newGit)
}

// VcsHandle is the handle
type VcsHandle struct {
}

func newGit() (vcs.Vcs, error) {
	return &VcsHandle{}, nil
}

// Name returns my name.
func (v VcsHandle) Name() string {
	return pluginName
}

// Discover returns false.
func (v VcsHandle) Discover(repobasedir string) bool {
	n := filepath.Join(repobasedir, ".git")
	found, err := bbutil.DirExists(n)
	if err != nil {
		return false
	}
	return found
}

// SetFileTypeUnix informs the VCS that files should maintain unix-style line endings.
func (v VcsHandle) SetFileTypeUnix(repobasedir string, files ...string) error {
	// Add to the .gitattributes in the same directory as the file.
	for _, file := range files {
		d, n := filepath.Split(file)
		err := bbutil.Touch(filepath.Join(repobasedir, d, ".gitattributes"))
		if err != nil {
			return err
		}
		err = bbutil.AddLinesToFile(filepath.Join(repobasedir, d, ".gitattributes"),
			fmt.Sprintf("%q text eol=lf", n))
		if err != nil {
			return err
		}
	}
	return nil
}

// IgnoreAnywhere tells the VCS to ignore these files anywhere rin the repo.
func (v VcsHandle) IgnoreAnywhere(repobasedir string, files ...string) error {
	// Add to the .gitignore file in the repobasedir.
	ignore := filepath.Join(repobasedir, ".gitignore")
	err := bbutil.Touch(ignore)
	if err != nil {
		return err
	}
	return bbutil.AddLinesToFile(ignore, files...)
}

// SuggestTracking tells the VCS to suggest the user commit these files.
func (v VcsHandle) SuggestTracking(repobasedir string, message string, files ...string) error {
	fmt.Print(`
NEXT STEP: You need to manually check these in:
     git commit -m'INITIALIZE BLACKBOX'`)
	for _, file := range files {
		fmt.Print(fmt.Sprintf(" %q", file))
	}
	fmt.Println()
	return nil
}

//echo "========== Encrypting: $unencrypted" >&2
//$GPG --use-agent --yes --trust-model=always --encrypt -o "$encrypted"  $(awk '{ print "-r" $1 }' < "$BB_ADMINS") "$unencrypted" >&2
//echo '========== Encrypting: DONE' >&2

// The following are "secret" functions only used by the integration testing system.

// TestingInitRepo initializes a repo.

func (v VcsHandle) TestingInitRepo() error {
	bbutil.RunBash("git", "init")
	return nil
}
