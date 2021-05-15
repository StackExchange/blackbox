// +build windows

package bbutil

// Umask is a no-op on Windows, and calls syscall.Umask on all other
// systems. On Windows it returns 0, which is a decoy.
func Umask(mask int) int {
	return 0o000
}
