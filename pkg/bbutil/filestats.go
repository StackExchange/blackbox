package bbutil

import (
	"fmt"
	"io/ioutil"
	"os"
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

// FileExistsOrProblem returns true if the file exists or if we can't determine its existence.
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
