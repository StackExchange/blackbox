package gnupg

import (
	"os"
	"os/exec"
	"syscall"

	"github.com/StackExchange/blackbox/v2/pkg/bbutil"
	"github.com/StackExchange/blackbox/v2/pkg/crypters"
)

var pluginName = "GnuPG"

func init() {
	crypters.Register(pluginName, 100, registerNew)
}

// CrypterHandle is the handle
type CrypterHandle struct {
	GPGCmd string // "gpg2" or "gpg"
}

func registerNew() (crypters.Crypter, error) {

	crypt := &CrypterHandle{}

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

// Decrypt decrypts a file, possibly overwriting the plaintext.
func (crypt CrypterHandle) Decrypt(name string, overwrite bool, umask int) error {

	if overwrite {
		_ = os.Remove(name)
	}

	oldumask := syscall.Umask(umask)
	err := bbutil.RunBash(crypt.GPGCmd,
		"--use-agent",
		"-q",
		"--decrypt",
		"-o", name,
		name+".gpg",
	)
	syscall.Umask(oldumask)
	return err
}
