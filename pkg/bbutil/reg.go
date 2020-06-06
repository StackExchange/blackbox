package bbutil

//
// import (
// 	"io/ioutil"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"sort"
// 	"strings"
//
// )
//
// // RegFile is a description of a registered file.
// type RegFile struct {
// 	Name string
// }
//
// // RegisteredFiles returns a list of the registered files.
// func (bbu *RepoInfo) RegisteredFiles() ([]RegFile, error) {
// 	blackboxFiles := filepath.Join(bbu.BlackboxConfigDir, "blackbox-files.txt")
// 	d, err := ioutil.ReadFile(blackboxFiles)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "Could not read the list of registered files")
// 	}
//
// 	// remove a trailing \n.
// 	// NB(tlim): We can't remove all trailing whitespace because filenames may contain whitespace.
// 	s := strings.TrimSuffix(string(d), "\n") // remove a single newline.
//
// 	names := strings.Split(s, "\n")
// 	if !sort.StringsAreSorted(names) {
// 		log.Fatalf("Files list is corrupted. It is not sorted; %q", blackboxFiles)
// 	}
// 	r := make([]RegFile, len(names))
// 	for i, name := range names {
// 		r[i].Name = name
// 	}
//
// 	return r, nil
// }
//
// // FileStatus returns the status of a file.
// func FileStatus(basedir, file string) string {
// 	p := filepath.Join(basedir, file)
// 	e := p + ".gpg"
// 	ps, perr := os.Stat(p)
// 	es, eerr := os.Stat(e)
// 	if perr == nil && eerr == nil {
// 		if ps.ModTime().Before(es.ModTime()) {
// 			return "GPGNEWER"
// 		}
// 		return "EDITING"
// 	}
// 	if eerr == nil {
// 		return "ENCRYPTED"
// 	}
// 	if os.IsExist(perr) {
// 		return "ERROR_NOGPG"
// 	}
// 	return "ERROR_NOEXIST"
// }
