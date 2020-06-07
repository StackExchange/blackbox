package box

// This file implements the business logic related to a black box.

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

// AdminAdd adds admins.
func (bx *Box) AdminAdd([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: AdminAdd")
}

// AdminList lists the admin id's.
func (bx *Box) AdminList() error {

	admins, err := bx.getAdmins()
	if err != nil {
		return err
	}

	for _, v := range admins {
		fmt.Println(v)
	}
	return nil
}

// AdminRemove removes an id from the admin list.
func (bx *Box) AdminRemove([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: AdminRemove")
}

// Cat outputs a file, unencrypting if needed.
func (bx *Box) Cat([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Cat")
}

// Decrypt decrypts a file.
func (bx *Box) Decrypt(names []string, overwrite bool, bulk bool, setgroup string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Decrypt")
}

// Diff ...
func (bx *Box) Diff([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Diff")
}

// Edit unencrypts, calls editor, calls encrypt.
func (bx *Box) Edit([]string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Edit")
}

// Encrypt encrypts a file.
func (bx *Box) Encrypt(names []string, bulk bool, setgroup string, overwrite bool) error {
	return fmt.Errorf("NOT IMPLEMENTED: Encrypt")
}

// FileAdd enrolls files.
func (bx *Box) FileAdd(names []string, overwrite bool) error {
	return fmt.Errorf("NOT IMPLEMENTED: FileAdd")
}

// FileList lists the files.
func (bx *Box) FileList() error {
	files, err := bx.getFiles()
	if err != nil {
		return err
	}
	for _, v := range files {
		fmt.Println(v)
	}
	return nil
}

// FileRemove de-enrolls files.
func (bx *Box) FileRemove(names []string) error {
	return fmt.Errorf("NOT IMPLEMENTED: FileRemove")
}

// Info prints debugging info.
func (bx *Box) Info() error {

	_, err := bx.getAdmins()
	if err != nil {
		logErr.Printf("getAdmins error: %v", err)
	}

	_, err = bx.getFiles()
	if err != nil {
		logErr.Printf("getFiles error: %v", err)
	}

	fmt.Println("BLACKBOX:")
	fmt.Printf("bx.ConfigDir=%q\n", bx.ConfigDir)
	//fmt.Printf("bx.Admins=%q\n", bx.Admins)
	fmt.Printf("len(bx.Admins)=%v\n", len(bx.Admins))
	//fmt.Printf("bx.Files=%q\n", bx.Files)
	fmt.Printf("len(bx.Files)=%v\n", len(bx.Files))
	fmt.Printf("bx.Vcs=%v\n", bx.Vcs)
	fmt.Printf("bx.VcsName=%q\n", bx.VcsName)

	return nil
}

// Init initializes a repo.
func (bx *Box) Init() error {
	return fmt.Errorf("NOT IMPLEMENTED: Init")
}

// Reencrypt decrypts and reencrypts files.
func (bx *Box) Reencrypt(names []string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Reencrypt")
}

// Shred shreds files.
func (bx *Box) Shred(names []string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Shred")
}

// func isChanged(pname string) (bool, error) {
// 	// if .gpg exists but not plainfile: unchanged
// 	// if plaintext exists but not .gpg: changed
// 	// if plainfile < .gpg: unchanged
// 	// if plainfile > .gpg: don't know, need to try diff

// 	// Gather info about the files:

// 	pstat, perr := os.Stat(pname)
// 	if perr != nil && (!os.IsNotExist(perr)) {
// 		return false, fmt.Errorf("isChanged(%q) returned error: %w", pname, perr)
// 	}
// 	gname := pname + ".gpg"
// 	gstat, gerr := os.Stat(gname)
// 	if gerr != nil && (!os.IsNotExist(perr)) {
// 		return false, fmt.Errorf("isChanged(%q) returned error: %w", gname, gerr)
// 	}

// 	pexists := perr == nil
// 	gexists := gerr == nil

// 	// Use the above rules:

// 	// if .gpg exists but not plainfile: unchanged
// 	if gexists && !pexists {
// 		return false, nil
// 	}

// 	// if plaintext exists but not .gpg: changed
// 	if pexists && !gexists {
// 		return true, nil
// 	}

// 	// At this point we can conclude that both p and g exist.
// 	//	Can't hurt to test that assertion.
// 	if (!pexists) && (!gexists) {
// 		return false, fmt.Errorf("Assertion failed. p and g should exist: pn=%q", pname)
// 	}

// 	pmodtime := pstat.ModTime()
// 	gmodtime := gstat.ModTime()
// 	// if plainfile < .gpg: unchanged
// 	if pmodtime.Before(gmodtime) {
// 		return false, nil
// 	}
// 	// if plainfile > .gpg: don't know, need to try diff
// 	return false, fmt.Errorf("Can not know for sure. Try git diff?")
// }

// FileStatus returns the status of a file.
func FileStatus(name string) (string, error) {
	/*
		DECRYPTED: File is decrypted and ready to edit (unknown if it has been edited).
		ENCRYPTED: GPG file is newer than plaintext. Indicates recented edited then encrypted.
		SHREDDED: Plaintext is missing.
		GPGMISSING: The .gpg file is missing. Oops?
		PLAINERROR: Can't access the plaintext file to determine status.
		GPGERROR: Can't access .gpg file to determine status.
	*/

	p := name
	e := p + ".gpg"
	ps, perr := os.Stat(p)
	es, eerr := os.Stat(e)
	if perr == nil && eerr == nil {
		if ps.ModTime().Before(es.ModTime()) {
			return "ENCRYPTED", nil
		}
		return "DECRYPTED", nil
	}

	if eerr != nil {
		if os.IsNotExist(eerr) {
			return "GPGMISSING", nil
		}
		return "GPGERROR", eerr
	}

	if perr != nil {
		if os.IsNotExist(perr) {
			return "SHREDDED", nil
		}
	}
	return "PLAINERROR", perr
}

// Status prints the status of files.
func (bx *Box) Status(names []string, nameOnly bool, match string) error {

	_, err := bx.getFiles()
	if err != nil {
		return err
	}

	var flist []string
	if len(names) == 0 {
		flist = bx.Files
	} else {
		flist = names
	}

	var data [][]string
	var onlylist []string
	thirdColumn := false
	var tcData bool

	for _, name := range flist {
		stat, err := FileStatus(name)
		if (match == "") || (stat == match) {
			if err == nil {
				data = append(data, []string{stat, name})
				onlylist = append(onlylist, name)
			} else {
				thirdColumn = tcData
				data = append(data, []string{stat, name, fmt.Sprintf("%v", err)})
				onlylist = append(onlylist, fmt.Sprintf("%v: %v", name, err))
			}
		}
	}

	if nameOnly {
		fmt.Println(strings.Join(onlylist, "\n"))
	} else {
		table := tablewriter.NewWriter(os.Stdout)
		table.SetAutoWrapText(false)
		if thirdColumn {
			table.SetHeader([]string{"Status", "Name", "Error"})
		} else {
			table.SetHeader([]string{"Status", "Name"})
		}
		for _, v := range data {
			table.Append(v)
		}
		table.Render() // Send output
	}

	return nil
}
