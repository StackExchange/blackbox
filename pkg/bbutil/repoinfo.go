package bbutil

// import (
// 	"os"
// 	"path/filepath"
//
// 	"github.com/StackExchange/blackbox/v2/pkg/bbgit"
// 	"github.com/StackExchange/blackbox/v2/pkg/bbnone"
// )
//
// // Vcser is the interface that defines a plug-in VCS system.
// type Vcser interface {
// 	Name() string        // Returns the name of this VCS type.
// 	RepoBaseDir() string // Returns the full path leading to this repo.
// }
//
// // RepoInfo stores info about the current repository.
// type RepoInfo struct {
// 	Vcs Vcser
// 	// BaseDir specifies the path (from "/") to the base of the VCS repo.
// 	RepoBaseDir       string // REPOBASE
// 	BlackboxConfigDir string // BLACKBOXDATA
// 	// KEYRINGDIR="$REPOBASE/$BLACKBOXDATA"
// 	// BB_ADMINS_FILE="blackbox-admins.txt"
// 	// BB_ADMINS="${KEYRINGDIR}/${BB_ADMINS_FILE}"
// 	// SECRING="${KEYRINGDIR}/secring.gpg"
// }
//
// // New is a factory.
// func New() (*RepoInfo, error) {
// 	repo := &RepoInfo{}
//
// 	vcs, err := vcsType()
// 	if err != nil {
// 		return nil, err
// 	}
// 	repo.Vcs = vcs
//
// 	// What is the base directory of the repo?
// 	base := os.Getenv("BLACKBOX_REPOBASE")
// 	if base == "" {
// 		base = repo.Vcs.RepoBaseDir()
// 	}
// 	repo.RepoBaseDir = base
//
// 	// Where are the blackbox config files?
// 	repo.BlackboxConfigDir, err = findConfigDir(repo.RepoBaseDir)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return repo, nil
// }
//
// // vcsType discovers the VCS type based on the BB_VCSTYPE env variable or by probing.
// func vcsType() (Vcser, error) {
// 	switch vcsName := os.Getenv("BB_VCSTYPE"); vcsName {
// 	case "git":
// 		return bbgit.New()
// 	// case "hg":
// 	// case "svn":
// 	default:
// 		break
// 	}
// 	answer, err := bbgit.New()
// 	if err == nil {
// 		return answer, nil
// 	}
// 	// TODO(tlim): Should we print err?
// 	// answer = bbhg.New()
// 	// if answer != nil {
// 	// 	return answer
// 	// }
// 	// answer = bbsvn.New()
// 	// if answer != nil {
// 	// 	return answer
// 	// }
// 	return bbnone.New()
// }
//
// // If BLACKBOXDATA is not set, search list this of directory paths.
// var configDirCandidates = []string{
// 	"keyrings/live",
// 	".blackbox", // Last item is the default.
// }
//
// // findConfigDir returns the configuration directory. It first checks the
// // BLACKBOXDATA env variable, then a list of candidates, lastly returning the
// // last candidate as the default.
// func findConfigDir(repoBase string) (string, error) {
// 	if dir := os.Getenv("BLACKBOXDATA"); dir != "" {
// 		//fmt.Fprintln(os.Stderr, "USING BBENV", dir)
// 		return filepath.Join(repoBase, dir), nil
// 	}
// 	var p string
// 	for _, c := range configDirCandidates {
// 		p = filepath.Join(repoBase, c)
// 		//fmt.Fprintf(os.Stderr, "Trying %q\n", p)
// 		if st, err := os.Stat(p); err == nil {
// 			mode := st.Mode()
// 			if mode.IsDir() {
// 				// FIXME(tlim): We are assuming that "not found"
// 				// and "i/o error" are both reasons to skip the candidates.
// 				// Maybe we should see what kind of error it is an output
// 				// some diagnostics if the problem is more than just "no found"?
// 				//fmt.Fprintf(os.Stderr, "RETURNING %q\n", p)
// 				return p, nil
// 			}
// 		}
// 	}
// 	return p, nil
// }
