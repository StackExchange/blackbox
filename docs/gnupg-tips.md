GnuPG tips
==========

# Common error messages

* Message: `gpg: filename: skipped: No public key`
* Solution: Usually this means there is an item in
  `.blackbox/blackbox-admins.txt` that is not the name of the key.
  Either something invalid was inserted (like a filename instead of a
  username) or a user has left the organization and their key was
  removed from the keychain, but their name wasn't removed from the
  blackbox-admins.txt file.

* Message: `gpg: decryption failed: No secret key`
* Solution: Usually means you forgot to re-encrypt the file with the new key.

* Message: `Error: can't re-encrypt because a key has expired.`
* Solution: A user's key has expired and can't be used to encrypt any more. Follow the [Replace expired keys](expired-keys.md) page.

FYI: Your repo may use `keyrings/live` instead of `.blackbox`. See "Where is the configuration stored?"

# GnuPG problems

Blackbox is just a front-end to GPG. If you get into a problem with a
key or file, you'll usually have better luck asking for advice on
the gnupg users mailing list TODO: Get link to this list


The author of Blackbox is not a GnuPG expert. He wrote Blackbox
because it was better than trying to remember GPG's horrible flag
names.
