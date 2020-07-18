package box

// This file implements the business logic related to a black box.
// These functions are usually called from cmd/blackbox/drive.go or
// external sytems that use box as a module.
import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/StackExchange/blackbox/v2/pkg/bbutil"
	"github.com/StackExchange/blackbox/v2/pkg/makesafe"
	"github.com/olekukonko/tablewriter"
)

// AdminAdd adds admins.
func (bx *Box) AdminAdd(nom string, sdir string) error {
	err := bx.getAdmins()
	if err != nil {
		return err
	}

	//fmt.Printf("ADMINS=%q\n", bx.Admins)

	// Check for duplicates.
	if i := sort.SearchStrings(bx.Admins, nom); i < len(bx.Admins) && bx.Admins[i] == nom {
		return fmt.Errorf("Admin %v already an admin", nom)
	}

	bx.logDebug.Printf("ADMIN ADD rbd=%q\n", bx.RepoBaseDir)
	changedFiles, err := bx.Crypter.AddNewKey(nom, bx.RepoBaseDir, sdir, bx.ConfigPath)
	if err != nil {
		return fmt.Errorf("AdminAdd failed AddNewKey: %v", err)
	}

	// TODO(tlim): Try the json file.

	// Try the legacy file:
	fn := filepath.Join(bx.ConfigPath, "blackbox-admins.txt")
	bx.logDebug.Printf("Admins file: %q", fn)
	err = bbutil.AddLinesToSortedFile(fn, nom)
	if err != nil {
		return fmt.Errorf("could not update file (%q,%q): %v", fn, nom, err)
	}
	changedFiles = append([]string{fn}, changedFiles...)

	bx.Vcs.NeedsCommit("NEW ADMIN: "+nom, bx.RepoBaseDir, changedFiles)
	return nil
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
func (bx *Box) Cat(names []string) error {
	if err := anyGpg(names); err != nil {
		return fmt.Errorf("cat: %w", err)
	}

	err := bx.getFiles()
	if err != nil {
		return err
	}

	for _, name := range names {
		var out []byte
		var err error
		if _, ok := bx.FilesSet[name]; ok {
			out, err = bx.Crypter.Cat(name)
		} else {
			out, err = ioutil.ReadFile(name)
		}
		if err != nil {
			bx.logErr.Printf("BX_CRY3\n")
			return fmt.Errorf("cat: %w", err)
		}
		fmt.Print(string(out))
	}
	return nil
}

// Decrypt decrypts a file.
func (bx *Box) Decrypt(names []string, overwrite bool, bulkpause bool, setgroup string) error {
	var err error

	if err := anyGpg(names); err != nil {
		return err
	}

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
		groupchange = true
	}
	bx.logDebug.Printf("DECRYPT GROUP %q %v,%v\n", setgroup, groupchange, gid)

	if len(names) == 0 {
		names = bx.Files
	}
	return decryptMany(bx, names, overwrite, groupchange, gid)
}

func decryptMany(bx *Box, names []string, overwrite bool, groupchange bool, gid int) error {

	// TODO(tlim): If we want to decrypt them in parallel, go has a helper function
	// called "sync.WaitGroup()"" which would be useful here.  We would probably
	// want to add a flag on the command line (stored in a field such as bx.ParallelMax)
	// that limits the amount of parallelism. The default for the flag should
	// probably be runtime.NumCPU().

	for _, name := range names {
		fmt.Printf("========== DECRYPTING %q\n", name)
		if !bx.FilesSet[name] {
			bx.logErr.Printf("Skipping %q: File not registered with Blackbox", name)
			continue
		}
		if (!overwrite) && bbutil.FileExistsOrProblem(name) {
			bx.logErr.Printf("Skipping %q: Will not overwrite existing file", name)
			continue
		}

		// TODO(tlim) v1 detects zero-length files and removes them, even
		// if overwrite is disabled. I don't think anyone has ever used that
		// feature. That said, if we want to do that, we would implement it here.

		// TODO(tlim) v1 takes the md5 hash of the plaintext before it decrypts,
		// then compares the new plaintext's md5. It prints "EXTRACTED" if
		// there is a change.

		err := bx.Crypter.Decrypt(name, bx.Umask, overwrite)
		if err != nil {
			bx.logErr.Printf("%q: %v", name, err)
			continue
		}

		// FIXME(tlim): Clone the file perms from the .gpg file to the plaintext file.

		if groupchange {
			// FIXME(tlim): Also "chmod g+r" the file.
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
func (bx *Box) Edit(names []string) error {

	if err := anyGpg(names); err != nil {
		return err
	}

	err := bx.getFiles()
	if err != nil {
		return err
	}

	for _, name := range names {
		if _, ok := bx.FilesSet[name]; ok {
			if !bbutil.FileExistsOrProblem(name) {
				err := bx.Crypter.Decrypt(name, bx.Umask, false)
				if err != nil {
					return fmt.Errorf("edit failed %q: %w", name, err)
				}
			}
		}
		err := bbutil.RunBash(bx.Editor, name)
		if err != nil {
			return err
		}
	}
	return nil
}

// Encrypt encrypts a file.
func (bx *Box) Encrypt(names []string, shred bool) error {
	var err error

	if err = anyGpg(names); err != nil {
		return err
	}

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

	enames, err := encryptMany(bx, names, shred)

	bx.Vcs.NeedsCommit(
		PrettyCommitMessage("ENCRYPTED", names),
		bx.RepoBaseDir,
		enames,
	)

	return err
}

func encryptMany(bx *Box, names []string, shred bool) ([]string, error) {
	var enames []string
	for _, name := range names {
		fmt.Printf("========== ENCRYPTING %q\n", name)
		if !bx.FilesSet[name] {
			bx.logErr.Printf("Skipping %q: File not registered with Blackbox", name)
			continue
		}
		if !bbutil.FileExistsOrProblem(name) {
			bx.logErr.Printf("Skipping. Plaintext does not exist: %q", name)
			continue
		}
		ename, err := bx.Crypter.Encrypt(name, bx.Umask, bx.Admins)
		if err != nil {
			bx.logErr.Printf("Failed to encrypt %q: %v", name, err)
			continue
		}
		enames = append(enames, ename)
		if shred {
			bx.Shred([]string{name})
		}
	}

	return enames, nil
}

// FileAdd enrolls files.
func (bx *Box) FileAdd(names []string, shred bool) error {
	bx.logDebug.Printf("FileAdd(shred=%v, %v)", shred, names)

	// Check for dups.
	// Encrypt them all.
	// If that succeeds, add to the blackbox-files.txt file.
	// (optionally) shred the plaintext.

	// FIXME(tlim): Check if the plaintext is in GIT.  If it is,
	// remove it from Git and print a warning that they should
	// eliminate the history or rotate any secrets.

	if err := anyGpg(names); err != nil {
		return err
	}

	err := bx.getAdmins()
	if err != nil {
		return err
	}
	err = bx.getFiles()
	if err != nil {
		return err
	}
	if err := anyGpg(names); err != nil {
		return err
	}

	// Check for newlines
	for _, n := range names {
		if strings.ContainsAny(n, "\n") {
			return fmt.Errorf("file %q contains a newlineregistered", n)
		}
	}

	// Check for duplicates.
	for _, n := range names {
		if i := sort.SearchStrings(bx.Files, n); i < len(bx.Files) && bx.Files[i] == n {
			return fmt.Errorf("file %q already registered", n)
		}
	}

	// Encrypt
	var needsCommit []string
	for _, name := range names {
		s, err := bx.Crypter.Encrypt(name, bx.Umask, bx.Admins)
		if err != nil {
			return fmt.Errorf("AdminAdd failed AddNewKey: %v", err)
		}
		needsCommit = append(needsCommit, s)
	}

	// TODO(tlim): Try the json file.

	// Try the legacy file:
	fn := filepath.Join(bx.ConfigPath, "blackbox-files.txt")
	bx.logDebug.Printf("Files file: %q", fn)
	err = bbutil.AddLinesToSortedFile(fn, names...)
	if err != nil {
		return fmt.Errorf("could not update file (%q,%q): %v", fn, names, err)
	}

	err = bx.Shred(names)
	if err != nil {
		bx.logErr.Printf("Error while shredding: %v", err)
	}

	bx.Vcs.CommitTitle("BLACKBOX ADD FILE: " + makesafe.FirstFew(makesafe.ShellMany(names)))

	bx.Vcs.IgnoreFiles(bx.RepoBaseDir, names)

	bx.Vcs.NeedsCommit(
		PrettyCommitMessage("blackbox-files.txt add", names),
		bx.RepoBaseDir,
		append([]string{filepath.Join(bx.ConfigPath, "blackbox-files.txt")}, needsCommit...),
	)
	return nil
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
		bx.logErr.Printf("Info getFiles: %v", err)
	}

	err = bx.getAdmins()
	if err != nil {
		bx.logErr.Printf("Info getAdmins: %v", err)
	}

	fmt.Println("BLACKBOX:")
	fmt.Printf("          Debug: %v\n", bx.Debug)
	fmt.Printf("           Team: %q\n", bx.Team)
	fmt.Printf("    RepoBaseDir: %q\n", bx.RepoBaseDir)
	fmt.Printf("     ConfigPath: %q\n", bx.ConfigPath)
	fmt.Printf("          Umask: %04o\n", bx.Umask)
	fmt.Printf("         Editor: %v\n", bx.Editor)
	fmt.Printf("       Shredder: %v\n", bbutil.ShredInfo())
	fmt.Printf("         Admins: count=%v\n", len(bx.Admins))
	fmt.Printf("          Files: count=%v\n", len(bx.Files))
	fmt.Printf("       FilesSet: count=%v\n", len(bx.FilesSet))
	fmt.Printf("            Vcs: %v\n", bx.Vcs)
	fmt.Printf("        VcsName: %q\n", bx.Vcs.Name())
	fmt.Printf("        Crypter: %v\n", bx.Crypter)
	fmt.Printf("    CrypterName: %q\n", bx.Crypter.Name())

	return nil
}

// Init initializes a repo.
func (bx *Box) Init(yes, vcsname string) error {
	fmt.Printf("VCS root is: %q\n", bx.RepoBaseDir)

	fmt.Printf("team is: %q\n", bx.Team)
	fmt.Printf("configdir will be: %q\n", bx.ConfigPath)

	if yes != "yes" {
		fmt.Printf("Enable blackbox for this %v repo? (yes/no)? ", bx.Vcs.Name())
		input := bufio.NewScanner(os.Stdin)
		input.Scan()
		ans := input.Text()
		b, err := strconv.ParseBool(ans)
		if err != nil {
			b = false
			if len(ans) > 0 {
				if ans[0] == 'y' || ans[0] == 'Y' {
					b = true
				}
			}
		}
		if !b {
			fmt.Println("Ok. Maybe some other time.")
			return nil
		}
	}

	err := os.Mkdir(bx.ConfigPath, 0o750)
	if err != nil {
		return err
	}

	ba := filepath.Join(bx.ConfigPath, "blackbox-admins.txt")
	bf := filepath.Join(bx.ConfigPath, "blackbox-files.txt")
	bbutil.Touch(ba)
	bbutil.Touch(bf)
	bx.Vcs.SetFileTypeUnix(bx.RepoBaseDir, ba, bf)

	bx.Vcs.IgnoreAnywhere(bx.RepoBaseDir, []string{
		"pubring.gpg~",
		"pubring.kbx~",
		"secring.gpg",
	})

	fs := []string{ba, bf}
	bx.Vcs.NeedsCommit(
		"NEW: "+strings.Join(makesafe.RedactMany(fs), " "),
		bx.RepoBaseDir,
		fs,
	)

	bx.Vcs.CommitTitle("INITIALIZE BLACKBOX")
	return nil
}

// Reencrypt decrypts and reencrypts files.
func (bx *Box) Reencrypt(names []string, overwrite bool, bulkpause bool) error {

	allFiles := false

	if err := anyGpg(names); err != nil {
		return err
	}
	if err := bx.getAdmins(); err != nil {
		return err
	}
	if err := bx.getFiles(); err != nil {
		return err
	}
	if len(names) == 0 {
		names = bx.Files
		allFiles = true
	}

	if bulkpause {
		gpgAgentNotice()
	}

	fmt.Println("========== blackbox administrators are:")
	bx.AdminList()
	fmt.Println("========== (the above people will be able to access the file)")

	if overwrite {
		bbutil.ShredFiles(names)
	} else {
		warned := false
		for _, n := range names {
			if bbutil.FileExistsOrProblem(n) {
				if !warned {
					fmt.Printf("========== Shred these files?\n")
					warned = true
				}
				fmt.Println("SHRED?", n)
			}
		}
		if warned {
			shouldWeOverwrite()
		}
	}

	// Decrypt
	if err := decryptMany(bx, names, overwrite, false, 0); err != nil {
		return fmt.Errorf("reencrypt failed decrypt: %w", err)
	}
	enames, err := encryptMany(bx, names, false)
	if err != nil {
		return fmt.Errorf("reencrypt failed encrypt: %w", err)
	}
	if err := bbutil.ShredFiles(names); err != nil {
		return fmt.Errorf("reencrypt failed shred: %w", err)
	}

	if allFiles {
		// If the "--all" flag was used, don't try to list all the files.
		bx.Vcs.NeedsCommit(
			"REENCRYPT all files",
			bx.RepoBaseDir,
			enames,
		)
	} else {
		bx.Vcs.NeedsCommit(
			PrettyCommitMessage("REENCRYPT", names),
			bx.RepoBaseDir,
			enames,
		)

	}

	return nil
}

// Shred shreds files.
func (bx *Box) Shred(names []string) error {

	if err := anyGpg(names); err != nil {
		return err
	}

	err := bx.getFiles()
	// Calling getFiles() has the benefit of making sure we are in a repo.
	if err != nil {
		return err
	}

	if len(names) == 0 {
		names = bx.Files
	}

	return bbutil.ShredFiles(names)
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
		return nil
	}

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
	if b, _ := bx.Vcs.Discover(); !b {
		return fmt.Errorf("TestingInitRepo failed Discovery")
	}
	return nil
}
