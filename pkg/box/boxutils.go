package box

import (
	"bufio"
	"fmt"
	"os"
	"os/user"
	"strconv"
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
	if err != nil {
		i, err = strconv.Atoi(g.Gid)
		return i, nil
	}

	// Give up.
	return -1, err
}

func gpgAgentNotice() {
	// Is gpg-agent configured?
	if os.Getenv("GPG_AGENT_INFO") != "" {
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
