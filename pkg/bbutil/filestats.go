package bbutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// DirExists returns true if directory exists.
func DirExists(path string) (bool, error) {
	stat, err := os.Stat(path)
	if err == nil {
		return stat.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// FileExistsOrProblem returns true if the file exists or if we can't determine its existance.
func FileExistsOrProblem(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return true
}

// Touch updates the timestamp of a file.
func Touch(name string) error {
	var err error
	_, err = os.Stat(name)
	if os.IsNotExist(err) {
		file, err := os.Create(name)
		if err != nil {
			return fmt.Errorf("TouchFile failed: %w", err)
		}
		file.Close()
	}

	currentTime := time.Now().Local()
	return os.Chtimes(name, currentTime, currentTime)
}

var shredPath, shredFlag string // Memoization cache

// shredCmd determines which command to use to securely erase a file. It returns
// the command to run and what flags to use with it. Determining the answer
// can be slow, therefore the answer is memoized and returned to future callers.
func shredCmd() (string, string) {
	// Use the memoized result.
	if shredPath != "" {
		return shredPath, shredFlag
	}

	var path string
	var err error
	if path, err = exec.LookPath("shred"); err == nil {
		shredPath, shredFlag = path, "-u"
	} else if path, err = exec.LookPath("srm"); err == nil {
		shredPath, shredFlag = path, "-f"
	} else if path, err = exec.LookPath("rm"); err == nil {
		shredPath, shredFlag = path, "-f"
		// Does this command support the "-P" flag?
		tmpfile, err := ioutil.TempFile("", "rmtest")
		defer os.Remove(tmpfile.Name()) // clean up
		err = RunBash("rm", "-P", tmpfile.Name())
		if err != nil {
			shredFlag = "-Pf"
		}
	}

	// Single exit, so we don't have to repeat the memoization code.
	return shredPath, shredFlag
}

// ShredFiles securely erases a list of files.
func ShredFiles(names []string) error {

	// TODO(tlim) DO the shredding in parallel like in v1.

	path, flag := shredCmd()
	var err error
	for _, n := range names {
		_, err := os.Stat(n)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("======= already gone: %q\n", n)
				continue
			}
		}
		fmt.Printf("========== SHREDDING: %q\n", n)
		e := RunBash(path, flag, n)
		if e != nil {
			err = e
			fmt.Printf("ERROR: %v\n", e)
		}
	}
	return err
}

// ReadFileLines is like ioutil.ReadFile() but returns an []string.
func ReadFileLines(filename string) ([]string, error) {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	s := string(b)
	s = strings.TrimSuffix(s, "\n")
	if s == "" {
		return []string{}, nil
	}
	l := strings.Split(s, "\n")
	return l, nil
}

// AddLinesToSortedFile adds a line to a sorted file.
func AddLinesToSortedFile(filename string, newlines ...string) error {
	lines, err := ReadFileLines(filename)
	//fmt.Printf("DEBUG: read=%q\n", lines)
	if err != nil {
		return fmt.Errorf("AddLinesToSortedFile can't read %q: %w", filename, err)
	}
	if !sort.StringsAreSorted(lines) {
		return fmt.Errorf("AddLinesToSortedFile: file wasn't sorted: %v", filename)
	}
	lines = append(lines, newlines...)
	sort.Strings(lines)
	contents := strings.Join(lines, "\n") + "\n"
	//fmt.Printf("DEBUG: write=%q\n", contents)
	err = ioutil.WriteFile(filename, []byte(contents), 0o660)
	if err != nil {
		return fmt.Errorf("AddLinesToSortedFile can't write %q: %w", filename, err)
	}
	return nil
}

// AddLinesToFile adds lines to the end of a file.
func AddLinesToFile(filename string, newlines ...string) error {
	lines, err := ReadFileLines(filename)
	if err != nil {
		return fmt.Errorf("AddLinesToFile can't read %q: %w", filename, err)
	}
	lines = append(lines, newlines...)
	contents := strings.Join(lines, "\n") + "\n"
	err = ioutil.WriteFile(filename, []byte(contents), 0o660)
	if err != nil {
		return fmt.Errorf("AddLinesToFile can't write %q: %w", filename, err)
	}
	return nil
}

// FindDirInParent looks for target in CWD, or .., or ../.., etc.
func FindDirInParent(target string) (string, error) {
	// Prevent an infinite loop by only doing "cd .." this many times
	maxDirLevels := 30
	relpath := "."
	for i := 0; i < maxDirLevels; i++ {
		// Does relpath contain our target?
		t := filepath.Join(relpath, target)
		//logDebug.Printf("Trying %q\n", t)
		_, err := os.Stat(t)
		if err == nil {
			return t, nil
		}
		if !os.IsNotExist(err) {
			return "", fmt.Errorf("stat failed FindDirInParent (%q): %w", t, err)
		}
		// Ok, it really wasn't found.

		// If we are at the root, stop.
		if abs, err := filepath.Abs(relpath); err == nil && abs == "/" {
			break
		}
		// Try one directory up
		relpath = filepath.Join("..", relpath)
	}
	return "", fmt.Errorf("Not found")
}
