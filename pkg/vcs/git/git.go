package git

import (
	"fmt"
	"path/filepath"

	"github.com/StackExchange/blackbox/v2/pkg/bbutil"
	"github.com/StackExchange/blackbox/v2/pkg/commitlater"
	"github.com/StackExchange/blackbox/v2/pkg/tainedname"
	"github.com/StackExchange/blackbox/v2/pkg/vcs"
)

var pluginName = "GIT"

func init() {
	vcs.Register(pluginName, 100, newGit)
}

// VcsHandle is the handle
type VcsHandle struct {
	toCommit commitlater.List
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

// Add makes a file visible to the VCS (like "git add").
func (v VcsHandle) Add(repobasedir string, files []string) error {

	// TODO(tlim): Make sure that files are within repobasedir.

	var gpgnames []string
	for _, n := range files {
		gpgnames = append(gpgnames, n+".gpg")
	}
	return bbutil.RunBash("git", append([]string{"add"}, gpgnames...)...)
}

// NeedsCommit queues up commits for later execution.
func (v VcsHandle) NeedsCommit(message string, repobasedir string, names []string) {
	v.toCommit.Add(message, repobasedir, names)
}

// FlushCommits informs the VCS to do queued up commits.
func (v VcsHandle) FlushCommits() error {
	return v.toCommit.Flush(v.suggestCommit)
	// TODO(tlim): Some day we can add a command line flag that indicates that commits are
	// to be done for real, not just suggested to the user.  At that point, this function
	// can call v.toCommit.Flush() with a function that actually does the commits insteada
	// of suggesting them.  Flag could be called --commit=auto vs --commit=suggest.
}

// suggestCommit tells the user what commits are needed.
func (v VcsHandle) suggestCommit(repobasedir string, message string, files []string) error {
	fmt.Printf(`
NEXT STEP: You need to manually check these in:
     git commit -m%q`, message)
	for _, file := range files {
		fmt.Print(" " + tainedname.New(file).String())
	}
	fmt.Println()
	return nil
}

// The following are "secret" functions only used by the integration testing system.

// TestingInitRepo initializes a repo.
func (v VcsHandle) TestingInitRepo() error {
	fmt.Println("RUNNING GIT INIT")
	return bbutil.RunBash("git", "init")

}
