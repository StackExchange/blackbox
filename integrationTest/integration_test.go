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
var longTests = flag.Bool("long", false, "Run long version of tests")

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
	if !*longTests {
		return
	}
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
	if !*longTests {
		return
	}
	compile(t)
	makeHomeDir(t, "init")

	runBB(t, "init", "yes")
	createDummyFilesAdmin(t)
	checkOutput("000-admin-list.txt", t, "admin", "list")
	checkOutput("000-file-list.txt", t, "file", "list")

	invalidArgs(t, "file", "list", "extra")
	invalidArgs(t, "admin", "list", "extra")
}

func TestStatus(t *testing.T) {
	if !*longTests {
		return
	}
	compile(t)
	makeHomeDir(t, "init")

	runBB(t, "init", "yes")
	createFilesStatus(t)
	checkOutput("000-status.txt", t, "status")
}

func TestShred(t *testing.T) {
	if !*longTests {
		return
	}
	compile(t)
	makeHomeDir(t, "shred")
	runBB(t, "init", "yes")

	makeFile(t, "shredme.txt", "File with SHREDME in it.\n")
	assertFileExists(t, "shredme.txt")
	runBB(t, "shred", "shredme.txt")
	assertFileMissing(t, "shredme.txt")
}

func TestStatus_notreg(t *testing.T) {
	if !*longTests {
		return
	}
	compile(t)
	makeHomeDir(t, "init")

	runBB(t, "init", "yes")
	createFilesStatus(t)
	checkOutput("status-noreg.txt", t, "status", "status-ENCRYPTED.txt", "blah.txt")
}

// TestHard tests the functions using a fake homedir and repo.
func TestHard(t *testing.T) {
	if !*longTests {
		return
	}
	// These are basic tests that work on a fake repo.
	// The repo has mostly real data, except any .gpg file
	// is just garbage.
	compile(t)
	setup(t)

	for _, cx := range []struct{ subname, prefix string }{
		//{subname: ".", prefix: "."},
		{subname: "mysub", prefix: ".."},
	} {
		subname := cx.subname
		prefix := cx.prefix
		_ = prefix

		phase("========== SUBDIR = " + subname + " ==========")

		makeHomeDir(t, "BasicAlice")

		plaintextFoo := "I am the foo.txt file!\n"
		plainAltered := "I am the altered file!\n"

		runBB(t, "testing_init") // Runs "git init" or equiv
		assertFileExists(t, ".git")
		runBB(t, "init", "yes") // Creates .blackbox or equiv

		if subname != "." {
			err := os.Mkdir(subname, 0770)
			if err != nil {
				t.Fatal(fmt.Errorf("hard-mk-home %q: %v", subname, err))
			}
		}
		olddir, err := os.Getwd()
		if err != nil {
			t.Fatal(err)
		}

		os.Chdir(subname)
		os.Chdir(olddir)

		phase("Alice creates a GPG key")
		gpgdir := makeAdmin(t, "alice", "Alice Example", "alice@example.com")
		become(t, "alice")

		phase("Alice enrolls as an admin")
		//os.Chdir(subname)
		runBB(t, "admin", "add", "alice@example.com", gpgdir)
		//os.Chdir(olddir)

		// encrypt
		phase("Alice registers foo.txt")
		makeFile(t, "foo.txt", plaintextFoo)
		//os.Chdir(subname)
		//runBB(t, "file", "add", "--shred", filepath.Join(prefix, "foo.txt"))
		runBB(t, "file", "add", "--shred", "foo.txt")
		//os.Chdir(olddir)
		// "file add" encrypts the file.
		// We shred the plaintext so that we are sure that when Decrypt runs,
		// we can verify the contents wasn't just sitting there all the time.
		assertFileMissing(t, "foo.txt")
		assertFileExists(t, "foo.txt.gpg")

		phase("Alice decrypts foo.txt")
		// decrypt
		//os.Chdir(subname)
		runBB(t, "decrypt", "foo.txt")
		//runBB(t, "decrypt", filepath.Join(prefix, "foo.txt"))
		//os.Chdir(olddir)
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
		phase("Alice reencrypts")
		checkOutput("basic-status.txt", t, "status")
		runBB(t, "reencrypt", "--overwrite", "foo.txt")

		// Test variations of cat

		// foo.txt=plain    result=plain
		phase("Alice cats plain:plain")
		makeFile(t, "foo.txt", plaintextFoo)
		assertFileExists(t, "foo.txt")
		runBB(t, "encrypt", "foo.txt")
		assertFileExists(t, "foo.txt")
		assertFileExists(t, "foo.txt.gpg")
		checkOutput("alice-cat-plain.txt", t, "cat", "foo.txt")
		assertFileExists(t, "foo.txt")
		assertFileExists(t, "foo.txt.gpg")

		// foo.txt=altered    result=plain
		phase("Alice cats altered:plain")
		makeFile(t, "foo.txt", plainAltered)
		assertFileExists(t, "foo.txt")
		assertFileExists(t, "foo.txt.gpg")
		checkOutput("alice-cat-plain.txt", t, "cat", "foo.txt")
		assertFileExists(t, "foo.txt")
		assertFileExists(t, "foo.txt.gpg")

		// foo.txt=missing  result=plain
		phase("Alice cats missing:plain")
		removeFile(t, "foo.txt")
		assertFileMissing(t, "foo.txt")
		assertFileMissing(t, "foo.txt")
		assertFileExists(t, "foo.txt.gpg")
		checkOutput("alice-cat-plain.txt", t, "cat", "foo.txt")
		assertFileMissing(t, "foo.txt")
		assertFileExists(t, "foo.txt.gpg")

		// Chapter 2: Bob
		// Alice adds Bob.
		// Bob encrypts a file.
		// Bob makes sure he can decrypt alice's file.
		// Bob removes Alice.
		// Alice verifies she CAN'T decrypt files.
		// Bob adds Alice back.
		// Alice verifies she CAN decrypt files.
		// Bob adds an encrypted file by mistake, "bb add" and fixes it.
		// Bob corrupts the blackbox-admins.txt file, verifies that commands fail.

	}

}

// TestEvilFilenames verifies commands work with "difficult" file names
func TestEvilFilenames(t *testing.T) {
	if !*longTests {
		return
	}
	compile(t)
	setup(t)
	makeHomeDir(t, "Mallory")

	runBB(t, "testing_init") // Runs "git init" or equiv
	assertFileExists(t, ".git")
	runBB(t, "init", "yes") // Creates .blackbox or equiv

	phase("Malory creates a GPG key")
	gpgdir := makeAdmin(t, "mallory", "Mallory Evil", "mallory@example.com")
	become(t, "mallory")

	phase("Mallory enrolls as an admin")
	runBB(t, "admin", "add", "mallory@example.com", gpgdir)

	_ = os.MkdirAll("my/path/to", 0o770)
	_ = os.Mkdir("other", 0o770)

	for i, name := range []string{
		"!important!.txt",
		"#andpounds.txt",
		"stars*bars?.txt",
		"space space.txt",
		"tab\ttab.txt",
		"ret\rret.txt",
		"smile😁eyes",
		"¡que!",
		"thé",
		"pound£",
		"*.go",
		"rm -f erase ; echo done",
		`smile☺`,
		`dub𝓦`,
		"my/path/to/relsecrets.txt",
		//"my/../my/path/../path/to/myother.txt",  // Not permitted yet
		//"other/../my//path/../path/to/otherother.txt", // Not permitted yet
		//"new\nnew.txt",  // \n not permitted
		//"two\n",  // \n not permitted (yet)
		//"four\U0010FFFF",  // Illegal byte sequence. git won't accept.
	} {
		phase(fmt.Sprintf("Mallory tries %02d: %q", i, name))
		contents := "the name of this file is the talking heads... i mean, " + name
		makeFile(t, name, contents)
		assertFileExists(t, name)
		assertFileMissing(t, name+".gpg")
		assertFileContents(t, name, contents)

		runBB(t, "file", "add", name)
		assertFileMissing(t, name)
		assertFileExists(t, name+".gpg")

		runBB(t, "decrypt", name)
		assertFileExists(t, name)
		assertFileExists(t, name+".gpg")
		assertFileContents(t, name, contents)

		runBB(t, "encrypt", name)
		assertFileExists(t, name)
		assertFileExists(t, name+".gpg")
		assertFileContents(t, name, contents)

		runBB(t, "shred", name)
		assertFileMissing(t, name)
		assertFileExists(t, name+".gpg")
	}
}

// More tests to implement.
// 1. Verify that the --gid works (blackbox decrypt --gid)
