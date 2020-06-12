package gnupg

import (
	"log"
	"os/exec"
	"syscall"

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

	oldumask := syscall.Umask(umask)
	err := bbutil.RunBash(crypt.GPGCmd, a...)
	syscall.Umask(oldumask)
	return err
}

// Encrypt name, overwriting name+".gpg"
func (crypt CrypterHandle) Encrypt(filename string, umask int, receivers []string) error {
	a := []string{
		"--use-agent",
		"--yes",
		"--trust-model=always",
		"--encrypt",
		"-o", filename + ".gpg",
	}
	for _, f := range receivers {
		a = append(a, "-r", f)
	}
	a = append(a, filename)

	oldumask := syscall.Umask(umask)
	crypt.logDebug.Printf("Args = %q", a)
	err := bbutil.RunBash(crypt.GPGCmd, a...)
	syscall.Umask(oldumask)

	return err
}
