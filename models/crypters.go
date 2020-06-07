package models

// CryptoSystem is gpg binaries, go-opengpg, etc.
type Crypter interface {
	// Decrypt name+".gpg", possibly overwriting name.
	Decrypt(filename string, overwrite bool) error
}
