package git

import (
	"fmt"
	"path/filepath"
	"strings"

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
	commitTitle  string
	toCommit     *commitlater.List // List of future commits
	commitHeader bool              // Has the "NEXT STEPS" header been printed?
}

func newGit() (vcs.Vcs, error) {
	l := &commitlater.List{}
	return &VcsHandle{toCommit: l}, nil
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
	seen := make(map[string]bool)

	// Add to the .gitattributes in the same directory as the file.
	for _, file := range files {
		d, n := filepath.Split(file)
		af := filepath.Join(repobasedir, d, ".gitattributes")
		err := bbutil.Touch(af)
		if err != nil {
			return err
		}
		err = bbutil.AddLinesToFile(af, fmt.Sprintf("%q text eol=lf", n))
		if err != nil {
			return err
		}
		seen[af] = true
	}

	var keys []string
	for k := range seen {
		keys = append(keys, k)
	}

	v.NeedsCommit(
		"set gitattr=UNIX "+tainedname.RedactList(files),
		repobasedir,
		keys,
	)

	return nil
}

// IgnoreAnywhere tells the VCS to ignore these files anywhere rin the repo.
func (v VcsHandle) IgnoreAnywhere(repobasedir string, files []string) error {
	// Add to the .gitignore file in the repobasedir.
	ignore := filepath.Join(repobasedir, ".gitignore")
	err := bbutil.Touch(ignore)
	if err != nil {
		return err
	}

	err = bbutil.AddLinesToFile(ignore, files...)
	if err != nil {
		return err
	}

	v.NeedsCommit(
		"gitignore "+tainedname.RedactList(files),
		repobasedir,
		[]string{".gitignore"},
	)
	return nil
}

func gitSafeFilename(name string) string {
	// TODO(tlim): Add unit tests.
	// TODO(tlim): Confirm that *?[] escaping works.
	if name == "" {
		return "ERROR"
	}
	var b strings.Builder
	b.Grow(len(name) + 2)
	for _, r := range name {
		if r == ' ' || r == '*' || r == '?' || r == '[' || r == ']' {
			b.WriteRune('\\')
			b.WriteRune(r)
		} else {
			b.WriteRune(r)
		}
	}
	if name[0] == '!' || name[0] == '#' {
		return `\` + b.String()
	}
	return b.String()
}

// IgnoreFiles tells the VCS to ignore these files, specified relative to RepoBaseDir.
func (v VcsHandle) IgnoreFiles(repobasedir string, files []string) error {

	var lines []string
	for _, f := range files {
		lines = append(lines, "/"+gitSafeFilename(f))
	}

	// Add to the .gitignore file in the repobasedir.
	ignore := filepath.Join(repobasedir, ".gitignore")
	err := bbutil.Touch(ignore)
	if err != nil {
		return err
	}
	err = bbutil.AddLinesToFile(ignore, lines...)
	if err != nil {
		return err
	}

	v.NeedsCommit(
		"gitignore "+tainedname.RedactList(files),
		repobasedir,
		[]string{".gitignore"},
	)
	return nil
}

// Add makes a file visible to the VCS (like "git add").
func (v VcsHandle) Add(repobasedir string, files []string) error {

	if len(files) == 0 {
		return nil
	}

	// TODO(tlim): Make sure that files are within repobasedir.

	var gpgnames []string
	for _, n := range files {
		gpgnames = append(gpgnames, n+".gpg")
	}
	return bbutil.RunBash("git", append([]string{"add"}, gpgnames...)...)
}

// CommitTitle indicates what the next commit title will be.
// This is used if a group of commits are merged into one.
func (v *VcsHandle) CommitTitle(title string) {
	v.commitTitle = title
}

// NeedsCommit queues up commits for later execution.
func (v *VcsHandle) NeedsCommit(message string, repobasedir string, names []string) {
	v.toCommit.Add(message, repobasedir, names)
}

// DebugCommits dumps the list of future commits.
func (v VcsHandle) DebugCommits() commitlater.List {
	return *v.toCommit
}

// FlushCommits informs the VCS to do queued up commits.
func (v VcsHandle) FlushCommits() error {
	return v.toCommit.Flush(
		v.commitTitle,
		func(files []string) error { return bbutil.RunBash("git", append([]string{"add"}, files...)...) },
		v.suggestCommit,
	)
	// TODO(tlim): Some day we can add a command line flag that indicates that commits are
	// to be done for real, not just suggested to the user.  At that point, this function
	// can call v.toCommit.Flush() with a function that actually does the commits insteada
	// of suggesting them.  Flag could be called --commit=auto vs --commit=suggest.
}

// suggestCommit tells the user what commits are needed.
func (v *VcsHandle) suggestCommit(messages []string, repobasedir string, files []string) error {
	if !v.commitHeader {
		fmt.Printf("NEXT STEP: You need to manually check these in:\n")
	}
	v.commitHeader = true

	fmt.Print(`     git commit -m'`, strings.Join(messages, `' -m'`)+`'`)
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
