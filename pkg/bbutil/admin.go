package bbutil

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

// Administrator is a description of the admininstrators.
type Administrator struct {
	Name string
}

// Administrators returns the list of administrators.
func (bbu *RepoInfo) Administrators() ([]Administrator, error) {
	adminFilename := filepath.Join(bbu.BlackboxConfigDir, "blackbox-admins.txt")
	d, err := ioutil.ReadFile(adminFilename)
	if err != nil {
		return nil, errors.Wrap(err, "Could not read the list of administrators")
	}

	// remove a trailing \n.
	s := strings.TrimSuffix(string(d), "\n") // remove a single newline.
	names := strings.Split(s, "\n")
	if !sort.StringsAreSorted(names) {
		log.Fatalf("Admin list is corrupted. It is not sorted; %q", adminFilename)
	}
	r := make([]Administrator, len(names))
	for i, name := range names {
		r[i].Name = name
	}

	return r, nil
}
