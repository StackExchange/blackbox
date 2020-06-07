package gnupg

import (
	"fmt"

	"github.com/StackExchange/blackbox/v2/pkg/crypters"
)

func init() {
	crypters.Register("GnuPG", 100, registerNew)
}

// CrypterHandle is the handle
type CrypterHandle struct {
}

func registerNew() (crypters.Crypter, error) {
	return &CrypterHandle{}, nil
}

// Decrypt decrypts a file, possibly overwriting the plaintext.
func (crypt CrypterHandle) Decrypt(name string, overwrite bool) error {
	fmt.Printf("WOULD decrypt %v (overwrite=%v)\n", name, overwrite)
	return nil
}
