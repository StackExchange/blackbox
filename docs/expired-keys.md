Replacing expired keys
======================

If someone's key has already expired, blackbox will stop
encrypting.  You see this error:

```
$ blackbox_edit_end modified_file.txt
--> Error: can't re-encrypt because a key has expired.
```

FYI: Your repo may use `keyrings/live` instead of `.blackbox`. See "Where is the configuration stored?"

You can also detect keys that are about to expire by issuing this command and manually reviewing the "expired:" dates:

    gpg --homedir=.blackbox  --list-keys

or... list UIDs that will expire within 1 month from today: (Warning: this also lists keys without an expiration date)

    gpg --homedir=.blackbox --list-keys  --with-colons --fixed-list-mode  | grep ^uid | awk -F: '$6 < '$(( $(date +%s) + 2592000))

Here's how to replace the key:

- Step 1. Administrator removes expired user:

Warning: This process will erase any unencrypted files that you were in the process of editing. Copy them elsewhere and restore the changes when done.

```
blackbox_removeadmin expired_user@example.com
# This next command overwrites any changed unencrypted files. See warning above.
blackbox_update_all_files
git commit -m "Re-encrypt all files"
gpg --homedir=.blackbox --delete-key expired_user@example.com
git commit -m 'Cleaned expired_user@example.com from keyring'  .blackbox/*
git push
```

- Step 2. Expired user adds an updated key:

```
git pull
blackbox_addadmin updated_user@example.com
git commit -m'NEW ADMIN: updated_user@example.com .blackbox/pubring.gpg .blackbox/trustdb.gpg .blackbox/blackbox-admins.txt
git push
```

- Step 3. Administrator re-encrypts all files with the updated key of the expired user:

```
git pull
gpg --import .blackbox/pubring.gpg
blackbox_update_all_files
git commit -m "Re-encrypt all files"
git push
```

- Step 4: Clean up:

Any files that were temporarily copied in the first step so as to not be overwritten can now be copied back and re-encrypted with the `blackbox_edit_end` command.

(Thanks to @chishaku for finding a solution to this problem!)

