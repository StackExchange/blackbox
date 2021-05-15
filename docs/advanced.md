Advanced Techniques
===================


# Using Blackbox without a repo

If the files are copied out of a repo they can still be decrypted and
edited. Obviously edits, changes to keys, and such will be lost if
they are made outside the repo. Also note that commands are most
likely to only work if run from the base directory (i.e. the parent to
the .blackbox directory).

Without a repo, all commands must be run from the same directory
as the ".blackbox" directory.  It might work otherwise but no
promises.


# Mixing gpg 1.x/2.0 and 2.2

WARNING: Each version of GnuPG uses a different, and incompatible,
binary format to store the keychain.  When Blackbox was originally
created, I didn't know this.  Things are mostly upwards compatible.
That said, if you have some admins with GnuPG 1.x and others with GnuPG 2.2,
you may corrupt the keychain.

A future version will store the keychain in an GnuPG-approved
version-neutral format.


# Having gpg and gpg2 on the same machine

NOTE: This is not implemented at this time. TODO(tlim) Use GPG to find
the binary.

In some situations, team members or automated roles need to install gpg
2.x alongside the system gpg version 1.x to catch up with the team's gpg
version. On Ubuntu 16, you can ```apt-get install gnupg2``` which
installs the binary gpg2. If you want to use this gpg2 binary, run every
blackbox command with GPG=gpg2.

For example:

```
GPG=gpg2 blackbox_postdeploy
```

