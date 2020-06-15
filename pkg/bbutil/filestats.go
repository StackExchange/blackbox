package bbutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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

// ShredFiles securely erases a list of files.
func ShredFiles(names []string) error {
	var path, flag string
	var err error
	if path, err = exec.LookPath("shred"); err == nil {
		flag = "-u"
	} else if path, err = exec.LookPath("srm"); err == nil {
		flag = "-f"
	} else if path, err = exec.LookPath("rm"); err == nil {
		flag = "-f"
		// FIXME(tlim): Test if "rm -P $tempfile" returns a error.
		// If it doesn't, flag = "-Pf"
	}

	// TODO(tlim) DO the shredding in parallel like in v1.

	for _, n := range names {
		fmt.Printf("SHREDDING (%q, %q): %q\n", path, flag, n)
		e := RunBash(path, flag, n)
		if e != nil {
			err = e
			fmt.Printf("ERROR: %v", e)
		} else {
			fmt.Println()
		}
	}
	fmt.Println("DONE.")
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
