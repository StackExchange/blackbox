package git

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/StackExchange/blackbox/v2/pkg/bbutil"
	"github.com/StackExchange/blackbox/v2/pkg/commitlater"
	"github.com/StackExchange/blackbox/v2/pkg/makesafe"
	"github.com/StackExchange/blackbox/v2/pkg/vcs"
)

var pluginName = "GIT"

func init() {
	vcs.Register(pluginName, 100, newGit)
}

// VcsHandle is the handle
type VcsHandle struct {
	commitTitle         string
	commitHeaderPrinted bool              // Has the "NEXT STEPS" header been printed?
	toCommit            *commitlater.List // List of future commits
}

func newGit() (vcs.Vcs, error) {
	l := &commitlater.List{}
	return &VcsHandle{toCommit: l}, nil
}

// Name returns my name.
func (v VcsHandle) Name() string {
	return pluginName
}

func ultimate(s string) int { return len(s) - 1 }

// Discover returns true if we are a repo of this type; along with the Abs path to the repo root (or "" if we don't know).
func (v VcsHandle) Discover() (bool, string) {
	out, err := bbutil.RunBashOutputSilent("git", "rev-parse", "--show-toplevel")
	if err != nil {
		return false, ""
	}
	if out == "" {
		fmt.Printf("WARNING: git rev-parse --show-toplevel has NO output??.  Seems broken.")
		return false, ""
	}
	if out[ultimate(out)] == '\n' {
		out = out[0:ultimate(out)]
	}
	return err == nil, out
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

	var changedfiles []string
	for k := range seen {
		changedfiles = append(changedfiles, k)
	}

	v.NeedsCommit(
		"set gitattr=UNIX "+strings.Join(makesafe.RedactMany(files), " "),
		repobasedir,
		changedfiles,
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
		"gitignore "+strings.Join(makesafe.RedactMany(files), " "),
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
		"gitignore "+strings.Join(makesafe.RedactMany(files), " "),
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
		func(files []string) error {
			return bbutil.RunBash("git", append([]string{"add"}, files...)...)
		},
		v.suggestCommit,
	)
	// TODO(tlim): Some day we can add a command line flag that indicates that commits are
	// to be done for real, not just suggested to the user.  At that point, this function
	// can call v.toCommit.Flush() with a function that actually does the commits instead
	// of suggesting them.  Flag could be called --commit=auto vs --commit=suggest.
}

// suggestCommit tells the user what commits are needed.
func (v *VcsHandle) suggestCommit(messages []string, repobasedir string, files []string) error {
	if !v.commitHeaderPrinted {
		fmt.Printf("NEXT STEP: You need to manually check these in:\n")
	}
	v.commitHeaderPrinted = true

	fmt.Print(`     git commit -m'`, strings.Join(messages, `' -m'`)+`'`)
	fmt.Print(" ")
	fmt.Print(strings.Join(makesafe.ShellMany(files), " "))
	fmt.Println()
	return nil
}

// The following are "secret" functions only used by the integration testing system.

// TestingInitRepo initializes a repo.
func (v VcsHandle) TestingInitRepo() error {
	return bbutil.RunBash("git", "init")

}
