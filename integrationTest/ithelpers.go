package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/StackExchange/blackbox/v2/pkg/vcs"
	_ "github.com/StackExchange/blackbox/v2/pkg/vcs/_all"
	"github.com/andreyvit/diff"
)

func getVcs(t *testing.T, name string) vcs.Vcs {
	t.Helper()
	// Set up the vcs
	for _, v := range vcs.Catalog {
		fmt.Printf("Testing vcs: %v == %v\n", name, v.Name)
		if strings.ToLower(v.Name) == strings.ToLower(name) {
			h, err := v.New()
			if err != nil {
				return nil // No idea how that would happen.
			}
			return h
		}
		fmt.Print("Nope.\n")

	}
	return nil
}

// TestBasicCommands's helpers

func createDummyRepo(t *testing.T, vcsname string) {
	t.Helper()
	fmt.Printf("createDummyRepo()\n")

	dir, err := ioutil.TempDir("", "repo")
	if err != nil {
		t.Fatalf("createDummyRepo: Could not make tempdir: %v", err)
	}
	defer os.RemoveAll(dir) // clean up

	os.Chdir(dir)

	runBB(t, "testing_init") // Runs "git init" and then vcs.Discover()
	runBB(t, "init", "yes")
	addLineSorted(t, ".blackbox/blackbox-admins.txt", "user1@example.com")
	addLineSorted(t, ".blackbox/blackbox-admins.txt", "user2@example.com")
	addLineSorted(t, ".blackbox/blackbox-files.txt", "foo.txt")
	addLineSorted(t, ".blackbox/blackbox-files.txt", "bar.txt")
	makeFile(t, "foo.txt", "I am the foo.txt file!")
	makeFile(t, "bar.txt", "I am the foo.txt file!")
	makeFile(t, "foo.txt.gpg", "V nz gur sbb.gkg svyr!")
	makeFile(t, "bar.txt.gpg", "V nz gur one.gkg svyr!")
}

func addLineSorted(t *testing.T, name string, newlines ...string) {
	t.Helper()

	contents, err := ioutil.ReadFile(name)
	if err != nil {
		t.Fatalf("addLinesSorted can't read %q: %v", name, err)
	}
	lines := strings.Split(string(contents), "\n")
	lines = append(lines, newlines...)
	sort.Strings(lines)
	err = ioutil.WriteFile(name, []byte(strings.Join(lines, "\n")), 0o666)
	if err != nil {
		t.Fatalf("addLinesSorted can't write %q: %v", name, err)
	}
}

func makeFile(t *testing.T, name string, lines ...string) {
	t.Helper()

	err := ioutil.WriteFile(name, []byte(strings.Join(lines, "\n")), 0o666)
	if err != nil {
		t.Fatalf("makeFile can't create %q: %v", name, err)
	}
}

// checkOutput runs blackbox with args, the last arg is the filename
// of the expected output. Error if output is not expected.
func checkOutput(t *testing.T, args ...string) {
	t.Helper()

	// Pop off the last arg. Use it as the filename for expected output.
	n := len(args) - 1
	name := args[n]
	args = args[:n]

	want, err := ioutil.ReadFile(filepath.Join("test_data", name))
	if err != nil {
		t.Fatalf("checkOutput can't read %v: %v", name, err)
	}

	cmd := exec.Command("blackbox", args...)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr
	got, err := cmd.Output()
	if err != nil {
		t.Fatal(fmt.Errorf("checkOutput(%q): %w", args, err))
	}

	if w, g := string(want), string(got); w != g {
		t.Errorf("checkOutput(%q) mismatch (-got +want):\n%s",
			args, diff.LineDiff(g, w))
	}

}

func invalidArgs(t *testing.T, args ...string) {
	t.Helper()

	fmt.Printf("invalidArgs(%q): \n", args)
	cmd := exec.Command(PathToBlackBox(), args...)
	cmd.Stdin = nil
	if *verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	if err == nil {
		fmt.Println("BAD")
		t.Fatal(fmt.Errorf("invalidArgs(%q): wanted failure but got success", args))
	}
	if *verbose {
		fmt.Printf("GOOD (expected): err=%q\n", err)
	} else {
		fmt.Println("GOOD (expected)")
	}
}

// TestAliceAndBob's helpers.

func setupUser(t *testing.T, user, passphrase string) {
	t.Helper()
	fmt.Printf("DEBUG: setupUser %q %q\n", user, passphrase)
}

var pathToBlackBox string

// PathToBlackBox returns the path to the executable we compile for integration testing.
func PathToBlackBox() string { return pathToBlackBox }

// SetPathToBlackBox sets the path.
func SetPathToBlackBox(n string) { pathToBlackBox = n }

func runBB(t *testing.T, args ...string) {
	t.Helper()

	fmt.Printf("runBB(%q)\n", args)
	cmd := exec.Command(PathToBlackBox(), args...)
	cmd.Stdin = nil
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatal(fmt.Errorf("runBB(%q): %w", args, err))
	}
}
