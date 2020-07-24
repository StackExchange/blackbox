package gnupg

/*

# How does Blackbox manage key rings?

Blackbox uses the user's .gnupg directory for most actions, such as decrypting data.
Decrypting requires the user's private key, which is stored by the user in their
home directory (and up to them to store safely).
Black box does not store the user's private key in the repo.

When encrypting data, black needs the public key of all the admins, not just the users.
To assure that the user's `.gnupg` has all these public keys, prior to
encrypting data the public keys are imported from .blackbox, which stores
a keychain that stores the public (not private!) keys of all the admins.

FYI: v1 does this import before decrypting, because I didn't know any better.

# Binary compatibility:

When writing v1, we didn't realize that the pubkey.gpg file is a binary format
that is not intended to be portable. In fact, it is intentionally not portable.
This means that all admins must use the exact same version of GnuPG
or the files (pubring.gpg or pubring.kbx) may get corrupted.

In v2, we store the public keys in the portable ascii format
in a file called `.blackbox/public-keys-db.asc`.
It will also update the binary files if they exist.
If `.blackbox/public-keys-db.asc` doesn't exist, it will be created.

Eventually we will stop updating the binary files.

# Importing public keys to the user

How to import the public keys to the user's GPG system:

If pubkeyring-ascii.txt exists:
	gpg --import pubkeyring-ascii.asc
Else if pubring.kbx
	gpg --import pubring.kbx
Else if pubring.gpg
	gpg --import pubring.gpg

This is what v1 does:
  #if gpg2 is installed next to gpg like on ubuntu 16
  if [[ "$GPG" != "gpg2" ]]; then
    $GPG --export --no-default-keyring --keyring "$(get_pubring_path)" >"$keyringasc"
    $GPG --import "$keyringasc" 2>&1 | egrep -v 'not changed$' >&2
  Else
    $GPG --keyring "$(get_pubring_path)" --export | $GPG --import
  fi

# How to add a key to the keyring?

Old, binary format:
    # Get the key they want to add:
        FOO is a user-specified directory, otherwise $HOME/.gnupg:
	    $GPG --homedir="FOO" --export -a "$KEYNAME" >TEMPFILE
	# Import into the binary files:
	    KEYRINGDIR is .blackbox
        $GPG --no-permission-warning --homedir="$KEYRINGDIR" --import TEMPFILE
	# Git add any of these files if they exist:
	    pubring.gpg pubring.kbx trustdb.gpg blackbox-admins.txt
	# Tell the user to git commit them.

New, ascii format:
	# Get the key to be added.  Write to a TEMPFILE
        FOO is a user-specified directory, otherwise $HOME/.gnupg:
	    $GPG --homedir="FOO" --export -a "$KEYNAME" >TEMPFILE
	# Make a tempdir called TEMPDIR
	# Import the pubkeyring-ascii.txt to TEMPDIR's keyring. (Skip if file not found)
	# Import the temp1 data to TEMPDIR
	# Export the TEMPDIR to create a new .blackbox/pubkeyring-ascii.txt
	    PATH_TO_BINARY is the path to .blackbox/pubring.gpg; if that's not found then pubring.kbx
        $GPG --keyring PATH_TO_BINARY --export -a --output .blackbox/pubkeyring-ascii.txt
	# Git add .blackbox/pubkeyring-ascii.txt and .blackbox/blackbox-admins.txt
	# Tell the user to git commit them.
	# Delete TEMPDIR

# How to remove a key from the keyring?

Old, binary format:
    # Remove key from the binary file
    $GPG --no-permission-warning --homedir="$KEYRINGDIR" --batch --yes --delete-key "$KEYNAME" || true
	# Git add any of these files if they exist:
	    pubring.gpg pubring.kbx trustdb.gpg blackbox-admins.txt
	# Tell the user to git commit them.

New, ascii format:
	# Make a tempdir called TEMPDIR
	# Import the pubkeyring-ascii.txt to TEMPDIR's keyring. (Skip if file not found)
    # Remove key from the ring file
    $GPG --no-permission-warning --homedir="$KEYRINGDIR" --batch --yes --delete-key "$KEYNAME" || true
	# Export the TEMPDIR to create a new .blackbox/pubkeyring-ascii.txt
	    PATH_TO_BINARY is the path to .blackbox/pubring.gpg; if that's not found then pubring.kbx
        $GPG --keyring PATH_TO_BINARY --export -a --output .blackbox/pubkeyring-ascii.txt
	# Git add .blackbox/pubkeyring-ascii.txt and .blackbox/blackbox-admins.txt
	# Update the .blackbox copy of pubring.gpg, pubring.kbx, or trustdb.gpg (if they exist)
	#     with copies from TEMPDIR (if they exist).  Git add any files that are updated.
	# Tell the user to git commit them.
	# Delete TEMPDIR

*/

//func prepareUserKeychain() error {
//	return nil
//}
