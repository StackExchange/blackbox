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
	compile(t)
	fmt.Printf("flag.testvcs is %v\n", *vcsToTest)
	vh := getVcs(t, *vcsToTest)
	fmt.Printf("Using BLACKBOX_FLAG_VCS=%v\n", vh.Name())
	os.Setenv("BLACKBOX_FLAG_VCS", vh.Name())

}

func TestInitInvalidArgs(t *testing.T) {
	compile(t)
	// Only zero or one args are permitted.
	invalidArgs(t, "init", "one", "two")
	invalidArgs(t, "init", "one", "two", "three")
}

func TestBasicCommands(t *testing.T) {
	setup(t)
	createDummyRepo(t, *vcsToTest)

	// admin
	checkOutput(t, "file", "list",
		"000-admin-list.txt",
	)
	invalidArgs(t, "admin", "list", "--all")
	invalidArgs(t, "admin", "one")

	// file
	checkOutput(t, "000-file-list.txt", "file", "list")
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
