package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/andreyvit/diff"
)

func assertFileMissing(t *testing.T, name string) {
	t.Helper()
	_, err := os.Stat(name)
	if err != nil && os.IsNotExist(err) {
		return
	}
	if err == nil {
		t.Fatalf("assertFileMissing failed: %v exists", name)
	}
	t.Fatalf("assertFileMissing: %q: %v", name, err)
}

func assertFileExists(t *testing.T, name string) {
	t.Helper()
	_, err := os.Stat(name)
	if err == nil {
		return
	}
	if os.IsNotExist(err) {
		t.Fatalf("assertFileExists failed: %v not exist", name)
	}
	t.Fatalf("assertFileExists: file can't be accessed: %v: %v", name, err)
}

func assertFileEmpty(t *testing.T, name string) {
	t.Helper()
	c, err := ioutil.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}
	if len(c) != 0 {
		t.Fatalf("got=%v want=%v: %v", len(c), 0, name)
	}
}

func assertFileContents(t *testing.T, name string, contents string) {
	t.Helper()
	c, err := ioutil.ReadFile(name)
	if err != nil {
		t.Fatal(err)
	}

	if w, g := contents, string(c); w != g {
		t.Errorf("assertFileContents(%q) mismatch (-got +want):\n%s",
			name, diff.LineDiff(g, w))
	}
}

func assertFilePerms(t *testing.T, name string, perms os.FileMode) {
	t.Helper()
	s, err := os.Stat(name)
	if err != nil {
		t.Fatal(err)
	}
	if s.Mode() != perms {
		t.Fatalf("got=%#o want=%#o: %v", s.Mode(), perms, name)
	}
}
