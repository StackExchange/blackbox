package box

// This file implements the business logic related to a black box.

import (
	"fmt"
	"os"
	"strings"

	"github.com/StackExchange/blackbox/v2/pkg/bbutil"
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
func (bx *Box) Decrypt(names []string, overwrite bool, bulkpause bool, setgroup string) error {
	var err error

	err = bx.getFiles()
	if err != nil {
		return err
	}

	groupchange := false
	gid := -1
	if setgroup != "" {
		gid, err = parseGroup(setgroup)
		if err != nil {
			return fmt.Errorf("Invalid group name or gid: %w", err)
		}
	}

	if bulkpause {
		gpgAgentNotice()
	}

	if len(names) == 0 {
		names = bx.Files
	}
	for _, name := range names {
		fmt.Printf("========== DECRYPTING %q\n", name)
		if !bx.FilesSet[name] {
			logErr.Printf("Skipping %q: File not registered with Blackbox", name)
		}
		if (!overwrite) && bbutil.FileExistsOrProblem(name) {
			logErr.Printf("Skipping %q: Will not overwrite existing file", name)
			continue
		}

		// TODO(tlim) v1 detects zero-length files and removes them, even
		// if overwrite is disabled. I don't think anyone has ever used that
		// feature. That said, we could immplement that here.

		err := bx.Crypter.Decrypt(name, overwrite, bx.Umask)
		if err != nil {
			logErr.Printf("%q: %v", name, err)
			continue
		}

		if groupchange {
			os.Chown(name, -1, gid)
		}
	}

	return nil
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
	err := bx.getFiles()
	if err != nil {
		return err
	}
	for _, v := range bx.Files {
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

	err = bx.getFiles()
	if err != nil {
		logErr.Printf("getFiles error: %v", err)
	}

	//fmt.Printf("bx.Admins=%q\n", bx.Admins)
	//fmt.Printf("bx.Files=%q\n", bx.Files)

	fmt.Println("BLACKBOX:")
	fmt.Printf("      ConfigDir: %q\n", bx.ConfigDir)
	fmt.Printf("    RepoBaseDir: %q\n", bx.RepoBaseDir)
	fmt.Printf("         Admins: count=%v\n", len(bx.Admins))
	fmt.Printf("          Files: count=%v\n", len(bx.Files))
	fmt.Printf("            Vcs: %v\n", bx.Vcs)
	fmt.Printf("        VcsName: %q\n", bx.VcsName)
	fmt.Printf("        Crypter: %v\n", bx.Crypter)
	fmt.Printf("    CrypterName: %q\n", bx.CrypterName)

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

// Status prints the status of files.
func (bx *Box) Status(names []string, nameOnly bool, match string) error {

	err := bx.getFiles()
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
