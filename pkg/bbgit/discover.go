package bbgit

// func New() (*GitInfo, error) {
// 	ri := new(GitInfo)
// 	path, err := exec.LookPath("git")
// 	if err != nil {
// 		return nil, nil
// 	}
// 	baseDir, err := exec.Command(path, "rev-parse", "--show-toplevel").Output()
// 	if err != nil {
// 		return nil, errors.Wrap(err, "bbgit:")
// 	}
// 	ri.baseDir = strings.TrimSuffix(string(baseDir), "\n") // remove a single newline.
// 	return ri, nil
// }
