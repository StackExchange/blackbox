package bbutil

import (
	"sort"
)

// FileIterator return a list of files to process.
func (bbu *RepoInfo) FileIterator(allFiles bool, fnames []string) ([]string, []bool, error) {
	regfiles, err := bbu.RegisteredFiles()
	if err != nil {
		return nil, nil, err
	}

	allnames := make([]string, len(regfiles))
	for i, r := range regfiles {
		allnames[i] = r.Name
	}

	if allFiles {
		isvalid := make([]bool, len(allnames))
		for n := range allnames {
			isvalid[n] = true
		}
		return allnames, isvalid, nil
	}

	retnames := make([]string, len(fnames))
	isvalid := make([]bool, len(fnames))
	for n, fn := range fnames {
		retnames[n] = fn
		i := sort.SearchStrings(allnames, fn)
		isvalid[n] = i < len(allnames) && allnames[i] == fn
	}

	return retnames, isvalid, nil
}
