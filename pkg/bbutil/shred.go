package bbutil

// Pick an appropriate secure erase command for this operating system
// or just delete the file with os.Remove().

// Code rewritten based https://codereview.stackexchange.com/questions/245072

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

var shredCmds = []struct {
	name, opts string
}{
	{"sdelete", "-a"},
	{"shred", "-u"},
	{"srm", "-f"},
	{"rm", "-Pf"},
}

func shredTemp(path, opts string) error {
	file, err := ioutil.TempFile("", "shredTemp.")
	if err != nil {
		return err
	}
	filename := file.Name()
	defer os.Remove(filename)
	defer file.Close()

	err = file.Close()
	if err != nil {
		return err
	}
	err = RunBash(path, opts, filename)
	if err != nil {
		return err
	}
	return nil
}

var shredPath, shredOpts = func() (string, string) {
	for _, cmd := range shredCmds {
		path, err := exec.LookPath(cmd.name)
		if err != nil {
			continue
		}
		err = shredTemp(path, cmd.opts)
		if err == nil {
			return path, cmd.opts
		}
	}
	return "", ""
}()

// ShredInfo reveals the shred command and flags (for "blackbox info")
func ShredInfo() string {
	return shredPath + " " + shredOpts
}

// shredFile shreds one file.
func shredFile(filename string) error {
	fi, err := os.Stat(filename)
	if err != nil {
		return err
	}
	if !fi.Mode().IsRegular() {
		err := fmt.Errorf("filename is not mode regular")
		return err
	}

	if shredPath == "" {
		// No secure erase command found.  Default to a normal file delete.
		// TODO(tlim): Print a warning? Have a flag that causes this to be an error?
		return os.Remove(filename)
	}

	err = RunBash(shredPath, shredOpts, filename)
	if err != nil {
		return err
	}
	return nil
}

// ShredFiles securely erases a list of files.
func ShredFiles(names []string) error {

	// TODO(tlim) DO the shredding in parallel like in v1.

	var eerr error
	for _, n := range names {
		_, err := os.Stat(n)
		if err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("======= already gone: %q\n", n)
				continue
			}
		}
		fmt.Printf("========== SHREDDING: %q\n", n)
		e := shredFile(n)
		if e != nil {
			eerr = e
			fmt.Printf("ERROR: %v\n", e)
		}
	}
	return eerr
}
