User Management
===============


# Who are the current admins?

```
blackbox admin list
```


# Add a new user (admin)

FYI: Your repo may use `keyrings/live` instead of `.blackbox`. See "Where is the configuration stored?"

`.blackbox/blackbox-admins.txt` is a file that lists which users are able to decrypt files. (More pedantically, it is a list of the GnuPG key names that the file is encrypted for.)

To join the list of people that can edit the file requires three steps; You create a GPG key and add it to the key ring. Then, someone that already has access adds you to the system. Lastly, you should test your access.

## Step 1: NEWPERSON creates a GPG key pair on a secure machine and add to public keychain.

If you don't already have a GPG key, here's how to generate one:

```
gpg --gen-key
```

WARNING: New versions of GPG generate keys which are not understood by
old versions of GPG.  If you generate a key with a new version of GPG,
this will cause problems for users of older versions of GPG.
Therefore it is recommended that you either assure that everyone using
Blackbox have the exact same version of GPG, or generate GPG keys
using a version of GPG as old as the oldest version of GPG used by
everyone using Blackbox.

Pick defaults for encryption settings, 0 expiration. Pick a VERY GOOD
passphrase. Store a backup of the private key someplace secure. For
example, keep the backup copy on a USB drive that is locked in safe.
Or, at least put it on a secure machine with little or no internet
access, full-disk-encryption, etc. Your employer probably has rules
about how to store such things.

FYI: If generating the key is slow, this is usually because the system
isn't generating enough entropy.  Tip: Open another window on that
machine and run this command: `ls -R /`

Now that you have a GPG key, add yourself as an admin:

```
blackbox admin add KEYNAME
```

...where "KEYNAME" is the email address listed in the gpg key you created previously. For example:

```
blackbox admin add tal@example.com
```

When the command completes successfully, instructions on how to commit these changes will be output. Run the command as given to commit the changes. It will look like this:

```
git commit -m'NEW ADMIN: tal@example.com' .blackbox/pubring.gpg .blackbox/trustdb.gpg .blackbox/blackbox-admins.txt
```


Then push it to the repo:

```
git push

or

ht push

(or whatever is appropriate)
```

NOTE: Creating a Role Account? If you are adding the pubring.gpg of a role account, you can specify the directory where the pubring.gpg file can be found as a 2nd parameter: `blackbox admin add puppetmaster@puppet-master-1.example.com /path/to/the/dir`

## Step 2: AN EXISTING ADMIN accepts you into the system.

Ask someone that already has access to re-encrypt the data files. This
gives you access. They simply decrypt and re-encrypt the data without
making any changes.

Pre-check: Verify the new keys look good.

```
git pull    # Or whatever is required for your system
gpg --homedir=.blackbox --list-keys
```

For example, examine the key name (email address) to make sure it conforms to corporate standards.

Import the keychain into your personal keychain and reencrypt:

```
gpg --import .blackbox/pubring.gpg
blackbox reencrypt --all shred
```

Push the re-encrypted files:

```
git commit -a
git push

or

hg commit
hg push
```

### Step 3: NEWPERSON tests.

Make sure you can decrypt a file. (Suggestion: Keep a dummy file in
VCS just for new people to practice on.)


# Remove a user

Simply run `blackbox admin remove` with their keyname then re-encrypt:

Example:

```
blackbox admin remove olduser@example.com
blackbox reencrypt --all shred
```

When the command completes, you will be given a reminder to check in the change and push it.

Note that their keys will still be in the key ring, but they will go unused. If you'd like to clean up the keyring, use the normal GPG commands and check in the file.

FYI: Your repo may use `keyrings/live` instead of `.blackbox`. See "Where is the configuration stored?"

```
gpg --homedir=.blackbox --list-keys
gpg --homedir=.blackbox --delete-key olduser@example.com
git commit -m'Cleaned olduser@example.com from keyring'  .blackbox/*
```

FYI: Your repo may use `keyrings/live` instead of `.blackbox`. See "Where is the configuration stored?"

The key ring only has public keys. There are no secret keys to delete.

Remember that this person did have access to all the secrets at one time. They could have made a copy. Therefore, to be completely secure, you should change all passwords, generate new SSL keys, and so on just like when anyone that had privileged access leaves an organization.

