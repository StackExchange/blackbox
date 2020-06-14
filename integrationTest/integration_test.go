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

	op, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	originPath = op
}

func TestInit(t *testing.T) {
	compile(t)

	// Only zero or one args are permitted.
	invalidArgs(t, "init", "one", "two")
	invalidArgs(t, "init", "one", "two", "three")

	runBB(t, "init", "yes")
	assertFileEmpty(t, ".blackbox/blackbox-admins.txt")
	assertFileEmpty(t, ".blackbox/blackbox-files.txt")
	assertFilePerms(t, ".blackbox/blackbox-admins.txt", 0o640)
	assertFilePerms(t, ".blackbox/blackbox-files.txt", 0o640)
}

func TestBasicCommands(t *testing.T) {
	// These are basic tests that work on a fake repo.
	// The repo has mostly real data, except any .gpg file
	// is just garbage.
	compile(t)
	setup(t)
	createDummyRepo(t, *vcsToTest)

	phase("Alice creates a repo.  Creates secret.txt.")
	makeFile(t, "secret.txt", "this is my secret")

	phase("Alice creates a GPG key...")
	makeAdmin(t, "alice", "Alice Example", "alice@example.com")
	become(t, "alice")

	runBB(t, "admin", "add", "alice@example.com")

}

func TestAlice(t *testing.T) {
	// Create an empty repo with a user named Alice who
	// performs many operations.  All files are valid.
	compile(t)
	setup(t)
	createDummyRepo(t, *vcsToTest)
	createDummyFilesAdmin(t)
	// encrypt
	runBB(t, "encrypt", "foo.txt")
	assertFileMissing(t, "foo.txt")
	assertFileExists(t, "foo.txt.gpg")

	// decrypt
	runBB(t, "decrypt", "foo.txt")
	assertFileExists(t, "foo.txt")
	assertFileExists(t, "foo.txt.gpg")

	// reencrypt

	// edit
	invalidArgs(t, "edit")
	invalidArgs(t, "edit", "--all")

	// cat

	// diff

	// shred

}

// func TestAliceAndBob(t *testing.T) {
// 	setupUser(t, "alice", "a")
// 	setupUser(t, "bob", "b")
// 	runBB(t, "init")
// 	runBB(t, "admin", "add", "alice@")
// }
