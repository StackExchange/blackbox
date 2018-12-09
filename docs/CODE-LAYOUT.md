

# External access

blackbox.go
  - cmd\_\*() -- Validate command line and call op\_\*() functions to do the work.
    - This code deals with all command line flags and ENV variables.
    - Prints status messages, errors, warns, etc.
    - Calls bbutil.go op\_() operations exclusively.
    - Determines VCS in use, GPG version in use.

pkg/bbutil.go:
  - op\_\*() -- Perform operations (can be used from other Go code)
    - This code is ignorant of flags and ENV variables.
    - This code is silent. Never prints to stdio/stderr. Returns errors for parent to print.
    - Decrypt(filename)
    - CopyPermissions(src, dst)
    - EncryptedFilename(plainfilename string) (encryptedfilename string)
    - UnencryptedFilename(plainfilename string) (encryptedfilename string)

file-level access

pkg/admin:
  - main.go -- generic admin manager
    - ListAdmins()
    - AddAdmins()
    - RemoveAdmins()
    - ListFiles()
    - FileStatus()
    - FileStatusAll()
    - IsOnFilelist()
    - IsNotOnFilelist()
  - plain.go -- read/write blackbox-admins.txt
    - listAdminsPlain()
    - addAdminsPlain()
    - removeAdminsPlain()
    - listFilesPlain()
    - fileStatusPlain()
    - fileStatusAllPlain()
  - FUTURE: a json equivalent of each plain function. Functions in main.go decide which to call.

crypto functions

models.go:
  interface for gpg.
    - Decrypt(encrypted, unencrypted) error
    - Encrypt(unencrypted, encrypted) error
    - GetPubKey(dirname, keyname) (pubkey string)
    - ImportPubKey(dirname, pubkey) touchedFiles []string
cryptcmd/bb\_gpg\_v1.go:
  - struct Gpgv1
cryptcmd/bb\_gpg\_v2.go:
  - struct Gpgv2

vcs/models.go
  interface for talk with VCS systems.
vcs/gpgv1/
vcs/gpgv2/
