package bbgit

import (
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// GitInfo contains Git-specific info about this repository.
type GitInfo struct {
	baseDir string
}

// New is a factory; returns error if this is not a Git repo.
func New() (*GitInfo, error) {
	ri := new(GitInfo)
	path, err := exec.LookPath("git")
	if err != nil {
		return nil, nil
	}
	baseDir, err := exec.Command(path, "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return nil, errors.Wrap(err, "bbgit:")
	}
	ri.baseDir = strings.TrimSuffix(string(baseDir), "\n") // remove a single newline.
	return ri, nil
}

// Name returns the name of this type of repo.
func (repo *GitInfo) Name() string {
	return "git"
}

// RepoBaseDir returns
func (repo *GitInfo) RepoBaseDir() string {
	return repo.baseDir
}
