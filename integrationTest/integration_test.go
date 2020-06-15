package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/StackExchange/blackbox/v2/pkg/bblog"
	_ "github.com/StackExchange/blackbox/v2/pkg/bblog"
	_ "github.com/StackExchange/blackbox/v2/pkg/vcs/_all"
)

var vcsToTest = flag.String("testvcs", "GIT", "VCS to test")

//var crypterToTest = flag.String("crypter", "GnuPG", "crypter to test")

func init() {
	testing.Init()
	flag.Parse()

	op, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	originPath = op
}

func compile(t *testing.T) {
	if PathToBlackBox() != "" {
		// It's been compiled already.
		return
	}
	// Make sure we have the latest binary
	fmt.Println("========== Compiling")
	cmd := exec.Command("go", "build", "-o", "../bbintegration", "../cmd/blackbox")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf("setup_compile: %v", err)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	SetPathToBlackBox(filepath.Join(cwd, "../bbintegration"))
}

func setup(t *testing.T) {
	logDebug := bblog.GetDebug(*verbose)

	logDebug.Printf("flag.testvcs is %v", *vcsToTest)
	vh := getVcs(t, *vcsToTest)
	logDebug.Printf("Using BLACKBOX_VCS=%v", vh.Name())
	os.Setenv("BLACKBOX_VCS", vh.Name())

}

func TestInit(t *testing.T) {
	compile(t)
	makeHomeDir(t, "init")

	// Only zero or one args are permitted.
	invalidArgs(t, "init", "one", "two")
	invalidArgs(t, "init", "one", "two", "three")

	runBB(t, "init", "yes")
	assertFileEmpty(t, ".blackbox/blackbox-admins.txt")
	assertFileEmpty(t, ".blackbox/blackbox-files.txt")
	assertFilePerms(t, ".blackbox/blackbox-admins.txt", 0o640)
	assertFilePerms(t, ".blackbox/blackbox-files.txt", 0o640)
}

func TestList(t *testing.T) {
	compile(t)
	makeHomeDir(t, "init")

	runBB(t, "init", "yes")
	createDummyFilesAdmin(t)
	checkOutput(t, "admin", "list", "000-admin-list.txt")
	checkOutput(t, "file", "list", "000-file-list.txt")

	invalidArgs(t, "file", "list", "extra")
	invalidArgs(t, "admin", "list", "extra")
}

func TestStatus(t *testing.T) {
	compile(t)
	makeHomeDir(t, "init")

	runBB(t, "init", "yes")
	createFilesStatus(t)
	checkOutput(t, "status", "000-status.txt")
}

func TestShred(t *testing.T) {
	compile(t)
	makeHomeDir(t, "shred")
	runBB(t, "init", "yes")

	makeFile(t, "shredme.txt", "File with SHREDME in it.")
	assertFileExists(t, "shredme.txt")
	runBB(t, "shred", "shredme.txt")
	assertFileMissing(t, "shredme.txt")
}

func TestStatus_notreg(t *testing.T) {
	compile(t)
	makeHomeDir(t, "init")

	runBB(t, "init", "yes")
	createFilesStatus(t)
	checkOutput(t, "status", "status-ENCRYPTED.txt", "blah.txt",
		"status-noreg.txt")
}

// TestBasicCommands tests of the basic functions, using a fake homedir and repo.
// The files are full of garbage, not real encrypted data.
func TestBasic(t *testing.T) {
	// These are basic tests that work on a fake repo.
	// The repo has mostly real data, except any .gpg file
	// is just garbage.
	compile(t)
	setup(t)
	makeHomeDir(t, "Basic")

	runBB(t, "testing_init") // Runs "git init" or equiv
	assertFileExists(t, ".git")
	runBB(t, "init", "yes") // Creates .blackbox or equiv

	phase("Alice creates a repo.  Creates secret.txt.")
	makeFile(t, "secret.txt", "this is my secret")

	phase("Alice creates a GPG key")
	gpgdir := makeAdmin(t, "alice", "Alice Example", "alice@example.com")
	become(t, "alice")

	phase("Alice enrolls as an admin")
	runBB(t, "admin", "add", "alice@example.com", gpgdir)

	// encrypt
	phase("Alice registers foo.txt")
	plaintextFoo := "I am the foo.txt file!\n"
	makeFile(t, "foo.txt", plaintextFoo)
	runBB(t, "file", "add", "--shred", "foo.txt")
	//runBB(t, "encrypt", "--shred", "foo.txt")
	// We shred the plaintext so that we are sure that when Decrypt runs,
	// we can verify the contents wasn't just sitting there all the time.
	assertFileMissing(t, "foo.txt")
	assertFileExists(t, "foo.txt.gpg")

	phase("Alice decrypts foo.txt")
	// decrypt
	runBB(t, "decrypt", "foo.txt")
	assertFileExists(t, "foo.txt")
	assertFileExists(t, "foo.txt.gpg")
	assertFileContents(t, "foo.txt", plaintextFoo)

	// encrypts (without shredding)
	phase("Alice encrypts foo.txt (again)")
	runBB(t, "encrypt", "foo.txt")
	assertFileExists(t, "foo.txt")
	assertFileExists(t, "foo.txt.gpg")
	assertFileContents(t, "foo.txt", plaintextFoo)

	// reencrypt

	// cat

	// diff

	// shred

}

// func TestAliceAndBob(t *testing.T) {
// 	setupUser(t, "alice", "a")
// 	setupUser(t, "bob", "b")
// 	runBB(t, "init")
// 	runBB(t, "admin", "add", "alice@")

// FYI: test "admins add" with multiple people.

// Edit requires a name, and doesn't work with --all.
//	invalidArgs(t, "edit")
//	invalidArgs(t, "edit", "--all")

// }
