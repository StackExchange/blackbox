package models

// Crypter is gpg binaries, go-opengpg, etc.
type Crypter interface {
	// Decrypt name+".gpg", possibly overwriting name.
	Decrypt(filename string, overwrite bool, umask int) error
}
