package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

var verbose = flag.Bool("verbose", false, "reveal stderr")
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
	fmt.Println("Compiling")
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
	fmt.Printf("flag.testvcs is %v\n", *vcsToTest)
	vh := getVcs(t, *vcsToTest)
	fmt.Printf("Using BLACKBOX_VCS=%v\n", vh.Name())
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
	// Verify blackbox-files.txt is empty, permissions (0o640)
	// Verify blackbox-admins.txt is empty, permissions (0o640)
}

func TestBasicCommands(t *testing.T) {
	// These are basic tests that work on a fake repo.
	// The repo has mostly real data, except any .gpg file
	// is just garbage.
	compile(t)
	setup(t)
	createDummyRepo(t, *vcsToTest)

	// admin
	checkOutput(t, "admin", "list", "000-admin-list.txt")
	invalidArgs(t, "admin", "list", "--all")
	invalidArgs(t, "admin", "one")

	// file
	checkOutput(t, "file", "list", "000-file-list.txt")
	invalidArgs(t, "file", "list", "one")
	invalidArgs(t, "file", "list", "--all")

	// status

	// reencrypt

	// decrypt

	// encrypt

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
