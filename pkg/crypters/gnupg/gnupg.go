package gnupg

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/StackExchange/blackbox/v2/pkg/bblog"
	"github.com/StackExchange/blackbox/v2/pkg/bbutil"
	"github.com/StackExchange/blackbox/v2/pkg/crypters"
)

var pluginName = "GnuPG"

func init() {
	crypters.Register(pluginName, 100, registerNew)
}

// CrypterHandle is the handle
type CrypterHandle struct {
	GPGCmd   string // "gpg2" or "gpg"
	logErr   *log.Logger
	logDebug *log.Logger
}

func registerNew(debug bool) (crypters.Crypter, error) {

	crypt := &CrypterHandle{
		logErr:   bblog.GetErr(),
		logDebug: bblog.GetDebug(debug),
	}

	// Which binary to use?
	path, err := exec.LookPath("gpg2")
	if err != nil {
		path, err = exec.LookPath("gpg")
		if err != nil {
			path = "gpg2"
		}
	}
	crypt.GPGCmd = path

	return crypt, nil
}

// Name returns my name.
func (crypt CrypterHandle) Name() string {
	return pluginName
}

// Decrypt name+".gpg", possibly overwriting name.
func (crypt CrypterHandle) Decrypt(filename string, umask int, overwrite bool) error {

	a := []string{
		"--use-agent",
		"-q",
		"--decrypt",
		"-o", filename,
	}
	if overwrite {
		a = append(a, "--yes")
	}
	a = append(a, filename+".gpg")

	oldumask := bbutil.Umask(umask)
	err := bbutil.RunBash(crypt.GPGCmd, a...)
	bbutil.Umask(oldumask)
	return err
}

// Cat returns the plaintext or, if it is missing, the decrypted cyphertext.
func (crypt CrypterHandle) Cat(filename string) ([]byte, error) {

	a := []string{
		"--use-agent",
		"-q",
		"--decrypt",
	}

	// TODO(tlim): This assumes the entire gpg file fits in memory. If
	// this becomes a problem, re-implement this using exec Cmd.StdinPipe()
	// and feed the input in chunks.
	in, err := ioutil.ReadFile(filename + ".gpg")
	if err != nil {

		if os.IsNotExist(err) {
			// Encrypted file doesn't exit? Return the plaintext.
			return ioutil.ReadFile(filename)
		}

		return nil, err
	}

	return bbutil.RunBashInputOutput(in, crypt.GPGCmd, a...)
}

// Encrypt name, overwriting name+".gpg"
func (crypt CrypterHandle) Encrypt(filename string, umask int, receivers []string) (string, error) {
	var err error

	crypt.logDebug.Printf("Encrypt(%q, %d, %q)", filename, umask, receivers)
	encrypted := filename + ".gpg"
	a := []string{
		"--use-agent",
		"--yes",
		"--trust-model=always",
		"--encrypt",
		"-o", encrypted,
	}
	for _, f := range receivers {
		a = append(a, "-r", f)
	}
	a = append(a, "--encrypt")
	a = append(a, filename)
	//err = bbutil.RunBash("ls", "-la")

	oldumask := bbutil.Umask(umask)
	crypt.logDebug.Printf("Args = %q", a)
	err = bbutil.RunBash(crypt.GPGCmd, a...)
	bbutil.Umask(oldumask)

	return encrypted, err
}

// AddNewKey extracts keyname from sourcedir's GnuPG chain to destdir keychain.
// It returns a list of files that may have changed.
func (crypt CrypterHandle) AddNewKey(keyname, repobasedir, sourcedir, destdir string) ([]string, error) {

	// $GPG --homedir="$2" --export -a "$KEYNAME" >"$pubkeyfile"
	args := []string{
		"--export",
		"-a",
	}
	if sourcedir != "" {
		args = append(args, "--homedir", sourcedir)
	}
	args = append(args, keyname)
	crypt.logDebug.Printf("ADDNEWKEY: Extracting key=%v: gpg, %v\n", keyname, args)
	pubkey, err := bbutil.RunBashOutput("gpg", args...)
	if err != nil {
		return nil, err
	}
	if len(pubkey) == 0 {
		return nil, fmt.Errorf("Nothing found when %q exported from %q", keyname, sourcedir)
	}

	// $GPG --no-permission-warning --homedir="$KEYRINGDIR" --import "$pubkeyfile"
	args = []string{
		"--no-permission-warning",
		"--homedir", destdir,
		"--import",
	}
	crypt.logDebug.Printf("ADDNEWKEY: Importing: gpg %v\n", args)
	// fmt.Printf("DEBUG: crypter ADD %q", args)
	err = bbutil.RunBashInput(pubkey, "gpg", args...)
	if err != nil {
		return nil, fmt.Errorf("AddNewKey failed: %w", err)
	}

	// Suggest: ${pubring_path} trustdb.gpg  blackbox-admins.txt
	var changed []string

	// Prefix each file with the relative path to it.
	prefix, err := filepath.Rel(repobasedir, destdir)
	if err != nil {
		//fmt.Printf("FAIL (%v) (%v) (%v)\n", repobasedir, destdir, err)
		prefix = destdir
	}
	for _, file := range []string{"pubring.gpg", "pubring.kbx", "trustdb.gpg"} {
		path := filepath.Join(destdir, file)
		if bbutil.FileExistsOrProblem(path) {
			changed = append(changed, filepath.Join(prefix, file))
		}
	}
	return changed, nil
}
