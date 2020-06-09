package models

// Crypter is gpg binaries, go-opengpg, etc.
type Crypter interface {
	// Name returns the plug-in's canonical name.
	Name() string
	// Decrypt name+".gpg", possibly overwriting name.
	Decrypt(filename string, overwrite bool, umask int) error
}
