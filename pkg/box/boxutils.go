package box

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/StackExchange/blackbox/v2/pkg/makesafe"
)

// FileStatus returns the status of a file.
func FileStatus(name string) (string, error) {
	/*
		DECRYPTED: File is decrypted and ready to edit (unknown if it has been edited).
		ENCRYPTED: GPG file is newer than plaintext. Indicates recented edited then encrypted.
		SHREDDED: Plaintext is missing.
		GPGMISSING: The .gpg file is missing. Oops?
		PLAINERROR: Can't access the plaintext file to determine status.
		GPGERROR: Can't access .gpg file to determine status.
	*/

	p := name
	e := p + ".gpg"
	ps, perr := os.Stat(p)
	es, eerr := os.Stat(e)
	if perr == nil && eerr == nil {
		if ps.ModTime().Before(es.ModTime()) {
			return "ENCRYPTED", nil
		}
		return "DECRYPTED", nil
	}

	if os.IsNotExist(perr) && os.IsNotExist(eerr) {
		return "BOTHMISSING", nil
	}

	if eerr != nil {
		if os.IsNotExist(eerr) {
			return "GPGMISSING", nil
		}
		return "GPGERROR", eerr
	}

	if perr != nil {
		if os.IsNotExist(perr) {
			return "SHREDDED", nil
		}
	}
	return "PLAINERROR", perr
}

func anyGpg(names []string) error {
	for _, name := range names {
		if strings.HasSuffix(name, ".gpg") {
			return fmt.Errorf(
				"no not specify .gpg files. Specify %q not %q",
				strings.TrimSuffix(name, ".gpg"), name)
		}
	}
	return nil
}

// func isChanged(pname string) (bool, error) {
// 	// if .gpg exists but not plainfile: unchanged
// 	// if plaintext exists but not .gpg: changed
// 	// if plainfile < .gpg: unchanged
// 	// if plainfile > .gpg: don't know, need to try diff

// 	// Gather info about the files:

// 	pstat, perr := os.Stat(pname)
// 	if perr != nil && (!os.IsNotExist(perr)) {
// 		return false, fmt.Errorf("isChanged(%q) returned error: %w", pname, perr)
// 	}
// 	gname := pname + ".gpg"
// 	gstat, gerr := os.Stat(gname)
// 	if gerr != nil && (!os.IsNotExist(perr)) {
// 		return false, fmt.Errorf("isChanged(%q) returned error: %w", gname, gerr)
// 	}

// 	pexists := perr == nil
// 	gexists := gerr == nil

// 	// Use the above rules:

// 	// if .gpg exists but not plainfile: unchanged
// 	if gexists && !pexists {
// 		return false, nil
// 	}

// 	// if plaintext exists but not .gpg: changed
// 	if pexists && !gexists {
// 		return true, nil
// 	}

// 	// At this point we can conclude that both p and g exist.
// 	//	Can't hurt to test that assertion.
// 	if (!pexists) && (!gexists) {
// 		return false, fmt.Errorf("Assertion failed. p and g should exist: pn=%q", pname)
// 	}

// 	pmodtime := pstat.ModTime()
// 	gmodtime := gstat.ModTime()
// 	// if plainfile < .gpg: unchanged
// 	if pmodtime.Before(gmodtime) {
// 		return false, nil
// 	}
// 	// if plainfile > .gpg: don't know, need to try diff
// 	return false, fmt.Errorf("Can not know for sure. Try git diff?")
// }

func parseGroup(userinput string) (int, error) {
	if userinput == "" {
		return -1, fmt.Errorf("group spec is empty string")
	}

	// If it is a valid number, use it.
	i, err := strconv.Atoi(userinput)
	if err == nil {
		return i, nil
	}

	// If not a number, look it up by name.
	g, err := user.LookupGroup(userinput)
	if err == nil {
		i, err = strconv.Atoi(g.Gid)
		return i, nil
	}

	// Give up.
	return -1, err
}

// FindConfigDir tests various places until it finds the config dir.
// If we can't determine the relative path, "" is returned.
func FindConfigDir(reporoot, team string) (string, error) {

	candidates := []string{}
	if team != "" {
		candidates = append(candidates, ".blackbox-"+team)
	}
	candidates = append(candidates, ".blackbox")
	candidates = append(candidates, "keyrings/live")
	logDebug.Printf("DEBUG: candidates = %q\n", candidates)

	maxDirLevels := 30 // Prevent an infinite loop
	relpath := "."
	for i := 0; i < maxDirLevels; i++ {
		// Does relpath contain any of our directory names?
		for _, c := range candidates {
			t := filepath.Join(relpath, c)
			logDebug.Printf("Trying %q\n", t)
			fi, err := os.Stat(t)
			if err == nil && fi.IsDir() {
				return t, nil
			}
			if err == nil {
				return "", fmt.Errorf("path %q is not a directory: %w", t, err)
			}
			if !os.IsNotExist(err) {
				return "", fmt.Errorf("dirExists access error: %w", err)
			}
		}

		// If we are at the root, stop.
		if abs, _ := filepath.Abs(relpath); abs == "/" {
			break
		}
		// Try one directory up
		relpath = filepath.Join("..", relpath)
	}

	return "", fmt.Errorf("No .blackbox (or equiv) directory found")
}

func gpgAgentNotice() {
	// Is gpg-agent configured?
	if os.Getenv("GPG_AGENT_INFO") != "" {
		return
	}
	// Are we on macOS?
	if runtime.GOOS == "darwin" {
		// We assume the use of https://gpgtools.org, which
		// uses the keychain.
		return
	}

	// TODO(tlim): v1 verifies that "gpg-agent --version" outputs a version
	// string that is 2.1.0 or higher.  It seems that 1.x is incompatible.

	fmt.Println("WARNING: You probably want to run gpg-agent as")
	fmt.Println("you will be asked for your passphrase many times.")
	fmt.Println("Example: $ eval $(gpg-agent --daemon)")
	fmt.Print("Press CTRL-C now to stop. ENTER to continue: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}

func shouldWeOverwrite() {
	fmt.Println()
	fmt.Println("WARNING: This will overwrite any unencrypted files laying about.")
	fmt.Print("Press CTRL-C now to stop. ENTER to continue: ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
}

// PrettyCommitMessage generates a pretty commit message.
func PrettyCommitMessage(verb string, files []string) string {
	if len(files) == 0 {
		// This use-case should probably be an error.
		return verb + " (no files)"
	}
	rfiles := makesafe.RedactMany(files)
	m, truncated := makesafe.FirstFewFlag(rfiles)
	if truncated {
		return verb + ": " + m
	}
	return verb + ": " + m
}
