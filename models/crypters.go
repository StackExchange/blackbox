package models

// Crypter is gpg binaries, go-opengpg, etc.
type Crypter interface {
	// Name returns the plug-in's canonical name.
	Name() string
	// Decrypt name+".gpg", possibly overwriting name.
	Decrypt(filename string, umask int, overwrite bool) error
	// Encrypt name, overwriting name+".gpg"
	Encrypt(filename string, umask int, receivers []string) (string, error)
	// Cat outputs a file, unencrypting if needed.
	Cat(filename string) ([]byte, error)
	// AddNewKey extracts keyname from sourcedir's GnuPG chain to destdir keychain.
	AddNewKey(keyname, repobasename, sourcedir, destdir string) ([]string, error)
}
