Set up automated users or "role accounts"
=========================================

TODO(tlim): I think this is overly complex. With GnuPG 2.2 and later,
you can use `--password '' --quick-generate-key userid` and you are
done. No need for subkeys.  Maybe rework this?

With role accounts, you have an automated system that needs to be able
to decrypt secrets without a password. This means the security of your
repo is based on how locked down the automation system is. This
is risky, so be careful.


i.e. This is how a Puppet Master can have access to the unencrypted data.

FYI: Your repo may use `keyrings/live` instead of `.blackbox`. See "Where is the configuration stored?"

An automated user (a "role account") is one that that must be able to decrypt without a passphrase. In general you'll want to do this for the user that pulls the files from the repo to the master. This may be automated with Jenkins CI or other CI system.

GPG keys have to have a passphrase. However, passphrases are optional on subkeys. Therefore, we will create a key with a passphrase then create a subkey without a passphrase. Since the subkey is very powerful, it should be created on a very secure machine.

There's another catch. The role account probably can't check files into Git/Mercurial. It probably only has read-only access to the repo. That's a good security policy. This means that the role account can't be used to upload the subkey public bits into the repo.

Therefore, we will create the key/subkey on a secure machine as yourself. From there we can commit the public portions into the repo. Also from this account we will export the parts that the role account needs, copy them to where the role account can access them, and import them as the role account.

ProTip: If asked to generate entropy, consider running this on the same machine in another window: `sudo dd if=/dev/sda of=/dev/null`

For the rest of this doc, you'll need to make the following substitutions:

- ROLEUSER: svc_deployacct or whatever your role account's name is.
- NEWMASTER: the machine this role account exists on.
- SECUREHOST: The machine you use to create the keys.

NOTE: This should be more automated/scripted. Patches welcome.

On SECUREHOST, create the puppet master's keys:

```
$ mkdir /tmp/NEWMASTER
$ cd /tmp/NEWMASTER
$ gpg --homedir . --gen-key
Your selection?
   (1) RSA and RSA (default)
What keysize do you want? (2048) DEFAULT
Key is valid for? (0) DEFAULT

# Real name: Puppet CI Deploy Account
# Email address: svc_deployacct@hostname.domain.name
```

NOTE: Rather than a real email address, use the username@FQDN of the host the key will be used on. If you use this role account on many machines, each should have its own key. By using the FQDN of the host, you will be able to know which key is which. In this doc, we'll refer to username@FQDN as $KEYNAME

Save the passphrase somewhere safe!

Create a sub-key that has no password:

```
$ gpg --homedir . --edit-key svc_deployacct
gpg> addkey
(enter passphrase)
  Please select what kind of key you want:
   (3) DSA (sign only)
   (4) RSA (sign only)
   (5) Elgamal (encrypt only)
   (6) RSA (encrypt only)
Your selection? 6
What keysize do you want? (2048)
Key is valid for? (0)
Command> key 2
(the new subkey has a "*" next to it)
Command> passwd
(enter the main key's passphrase)
(enter an empty passphrase for the subkey... confirm you want to do this)
Command> save
```

Now securely export this directory to NEWMASTER:

```
gpg --homedir . --export -a svc_sadeploy >/tmp/NEWMASTER/pubkey.txt
tar cvf /tmp/keys.tar .
rsync -avP /tmp/keys.tar NEWMASTER:/tmp/.
```

On NEWMASTER, receive the new GnuPG config:

```
sudo -u svc_deployacct bash
mkdir -m 0700 -p ~/.gnupg
cd ~/.gnupg && tar xpvf /tmp/keys.tar
```

<!---
Back on SECUREHOST, import the pubkey into the repository.

```
$ cd .blackbox
$ gpg --homedir . --import /tmp/NEWMASTER/pubkey.txt
```
-->

Back on SECUREHOST, add the new email address to .blackbox/blackbox-admins.txt:

```
cd /path/to/the/repo
blackbox_addadmin $KEYNAME /tmp/NEWMASTER
```

Verify that secring.gpg is a zero-length file. If it isn't, you have somehow added a private key to the keyring. Start over.

```
cd .blackbox
ls -l secring.gpg
```

Commit the recent changes:

```
cd .blackbox
git commit -m"Adding key for KEYNAME" pubring.gpg trustdb.gpg blackbox-admins.txt
```

Regenerate all encrypted files with the new key:

```
blackbox_update_all_files
git status
git commit -m"updated encryption" -a
git push
```

On NEWMASTER, import the keys and decrypt the files:

```
sudo -u svc_sadeploy bash   # Become the role account.
gpg --import /etc/puppet/.blackbox/pubring.gpg
export PATH=$PATH:/path/to/blackbox/bin
blackbox_postdeploy
sudo -u puppet cat /etc/puppet/hieradata/blackbox.yaml # or any encrypted file.
```

ProTip: If you get "gpg: decryption failed: No secret key" then you forgot to re-encrypt blackbox.yaml with the new key.

On SECUREHOST, securely delete your files:

```
cd /tmp/NEWMASTER
# On machines with the "shred" command:
shred -u /tmp/keys.tar
find . -type f -print0 | xargs -0 shred -u
# All else:
rm -rf /tmp/NEWMASTER
```

Also shred any other temporary files you may have made.



