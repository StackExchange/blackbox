package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var sha = flag.String("sha", "", "SHA of current commit")

var goos = flag.String("os", "", "OS to build (linux, windows, or darwin) Defaults to all.")

func main() {
	flag.Parse()
	flags := fmt.Sprintf(`-s -w -X main.SHA="%s" -X main.BuildTime=%d`, getVersion(), time.Now().Unix())
	pkg := "github.com/StackExchange/blackbox/v2/cmd/blackbox"

	build := func(out, goos string) {
		log.Printf("Building %s", out)
		cmd := exec.Command("go", "build", "-o", out, "-ldflags", flags, pkg)
		os.Setenv("GOOS", goos)
		os.Setenv("GO111MODULE", "on")
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, env := range []struct {
		binary, goos string
	}{
		{"blackbox-Linux", "linux"},
		{"blackbox.exe", "windows"},
		{"blackbox-Darwin", "darwin"},
	} {
		if *goos == "" || *goos == env.goos {
			build(env.binary, env.goos)
		}
	}
}

func getVersion() string {
	if *sha != "" {
		return *sha
	}
	// check teamcity build version
	if v := os.Getenv("BUILD_VCS_NUMBER"); v != "" {
		return v
	}
	// check git
	cmd := exec.Command("git", "rev-parse", "HEAD")
	v, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	ver := strings.TrimSpace(string(v))
	// see if dirty
	cmd = exec.Command("git", "diff-index", "--quiet", "HEAD", "--")
	err = cmd.Run()
	// exit status 1 indicates dirty tree
	if err != nil {
		if err.Error() == "exit status 1" {
			ver += "[dirty]"
		} else {
			log.Printf("!%s!", err.Error())
		}
	}
	return ver
}
