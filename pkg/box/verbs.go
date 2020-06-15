package box

// This file implements the business logic related to a black box.

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
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
	err := bx.getAdmins()
	if err != nil {
		return err
	}

	for _, v := range bx.Admins {
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

	if bulkpause {
		gpgAgentNotice()
	}

	groupchange := false
	gid := -1
	if setgroup != "" {
		gid, err = parseGroup(setgroup)
		if err != nil {
			return fmt.Errorf("Invalid group name or gid: %w", err)
		}
	}

	if len(names) == 0 {
		names = bx.Files
	}
	for _, name := range names {
		fmt.Printf("========== DECRYPTING %q\n", name)
		if !bx.FilesSet[name] {
			bx.logErr.Printf("Skipping %q: File not registered with Blackbox", name)
		}
		if (!overwrite) && bbutil.FileExistsOrProblem(name) {
			bx.logErr.Printf("Skipping %q: Will not overwrite existing file", name)
			continue
		}

		// TODO(tlim) v1 detects zero-length files and removes them, even
		// if overwrite is disabled. I don't think anyone has ever used that
		// feature. That said, we could immplement that here.

		// TODO(tlim) v1 takes the md5 has of the plaintext before it decrypts,
		// then compares the new plaintext's md5. It prints "EXTRACTED" if
		// there is a change.

		err := bx.Crypter.Decrypt(name, bx.Umask, overwrite)
		if err != nil {
			bx.logErr.Printf("%q: %v", name, err)
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
func (bx *Box) Encrypt(names []string, umask int, shred bool) error {
	var err error

	err = bx.getAdmins()
	if err != nil {
		return err
	}

	err = bx.getFiles()
	if err != nil {
		return err
	}
	if len(names) == 0 {
		names = bx.Files
	}

	var suggest []string
	for _, name := range names {
		fmt.Printf("========== ENCRYPTING %q\n", name)
		if !bx.FilesSet[name] {
			bx.logErr.Printf("Skipping %q: File not registered with Blackbox", name)
		}
		err := bx.Crypter.Encrypt(name, bx.Umask, bx.Admins)
		if err != nil {
			bx.logErr.Printf("Failed to encrypt %q: %v", name, err)
			continue
		}
		suggest = append(suggest, fmt.Sprintf("Updated: %q", name))
		if shred {
			bx.Shred(name)
		}
	}

	bx.Vcs.SuggestTracking(bx.RepoBaseDir,
		strings.Join(names, "\n")+"\n",
		names...,
	)

	return nil
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

	err := bx.getFiles()
	if err != nil {
		bx.logErr.Printf("Info: %v", err)
	}

	err = bx.getAdmins()
	if err != nil {
		bx.logErr.Printf("Info: %v", err)
	}

	//fmt.Printf("bx.Admins=%q\n", bx.Admins)
	//fmt.Printf("bx.Files=%q\n", bx.Files)

	fmt.Println("BLACKBOX:")
	fmt.Printf("      ConfigDir: %q\n", bx.ConfigDir)
	fmt.Printf("    RepoBaseDir: %q\n", bx.RepoBaseDir)
	fmt.Printf("         Admins: count=%v\n", len(bx.Admins))
	fmt.Printf("          Files: count=%v\n", len(bx.Files))
	fmt.Printf("            Vcs: %v\n", bx.Vcs)
	fmt.Printf("        VcsName: %q\n", bx.Vcs.Name())
	fmt.Printf("        Crypter: %v\n", bx.Crypter)
	fmt.Printf("    CrypterName: %q\n", bx.Crypter.Name())

	return nil
}

// Init initializes a repo.
func (bx *Box) Init(yes, vcsname string) error {
	//fmt.Printf("VCS root is: %q\n", bx.RepoBaseDir)

	//fmt.Printf("team is: %q\n", bx.Team)
	//fmt.Printf("configdir will be: %q\n", bx.ConfigDir)

	if yes != "yes" {
		fmt.Printf("Enable blackbox for this %v repo? (yes/no)", bx.Vcs.Name())
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		b, _ := strconv.ParseBool(input.Text())
		if !b {
			fmt.Println("Ok. Maybe some other time.")
			return nil
		}
	}

	err := os.Mkdir(bx.ConfigDir, 0o750)
	if err != nil {
		return err
	}

	bbadmins := filepath.Join(bx.ConfigDir, "blackbox-admins.txt")
	bbutil.Touch(bbadmins)
	bbadminsRel := filepath.Join(bx.ConfigDirRel, "blackbox-admins.txt")
	bx.Vcs.SetFileTypeUnix(bx.RepoBaseDir, bbadminsRel)

	bbfiles := filepath.Join(bx.ConfigDir, "blackbox-files.txt")
	bbutil.Touch(bbfiles)
	bbfilesRel := filepath.Join(bx.ConfigDirRel, "blackbox-files.txt")
	bx.Vcs.SetFileTypeUnix(bx.RepoBaseDir, bbfilesRel)

	bx.Vcs.IgnoreAnywhere(bx.RepoBaseDir,
		"pubring.gpg~",
		"pubring.kbx~",
		"secring.gpg",
	)

	bx.Vcs.SuggestTracking(bx.RepoBaseDir, "INITIALIZE BLACKBOX",
		bbadminsRel, bbfilesRel,
	)

	return nil
}

// Reencrypt decrypts and reencrypts files.
func (bx *Box) Reencrypt(names []string) error {
	return fmt.Errorf("NOT IMPLEMENTED: Reencrypt")
}

// Shred shreds files.
func (bx *Box) Shred(names ...string) error {
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
		var stat string
		var err error
		if _, ok := bx.FilesSet[name]; ok {
			stat, err = FileStatus(name)
		} else {
			stat, err = "NOTREG", nil
		}
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

// TestingInitRepo initializes a repo.
// Uses bx.Vcs to create ".git" or whatever.
// Uses bx.Vcs to discover what was created, testing its work.
func (bx *Box) TestingInitRepo() error {

	if bx.Vcs == nil {
		fmt.Println("bx.Vcs is nil")
		fmt.Printf("BLACKBOX_VCS=%q\n", os.Getenv("BLACKBOX_VCS"))
		os.Exit(1)
	}
	fmt.Printf("ABOUT TO CALL TestingInitRepo\n")
	fmt.Printf("vcs = %v\n", bx.Vcs.Name())
	err := bx.Vcs.TestingInitRepo()
	fmt.Printf("RETURNED from TestingInitRepo: %v\n", err)
	fmt.Println(os.Getwd())
	if err != nil {
		return fmt.Errorf("TestingInitRepo returned: %w", err)
	}
	if !bx.Vcs.Discover("") {
		return fmt.Errorf("TestingInitRepo failed Discovery")
	}
	return nil
}
