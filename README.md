BlackBox
========

Safely store secrets in Git/Hg.  These commands make it easy
for you to GPG encrypt specific files in a repo so they are
"encrypted at rest" in your repository. However, the scripts
make it easy to decrypt them when you need to view or edit them,
and decrypt them for for use in production. Originally written
for Puppet, now works with any Git or Mercurial repository.


Overview
========

The goal is to have secret bits (passwords, private keys, and such) in your VCS repo but encrypted so that it is safe.  On the puppet masters they sit on disk unencrypted but are readable (decrypted) for use by the Puppet Master (or whoever needs full access).

How does this work?
===================

**Private keys (and anything that is the entire file):**  

Files are kept in git/hg encrypted (foo.txt is stored as foo.txt.gpg).

After deploying an update to your Puppet Master, the master runs a script that decrypts them.  The sit unencrypted on the master, which should already be locked down.

**Passwords (and any short string):**
Passwords are kept in hieradata/blackbox.yaml.gpg, which is decrypted to become hieradata/blackbox.yaml.  This data can be read by hiera.  This file is encrypted/decrypted just like any other blackbox file.

**Key management:**
The Puppet Masters have GPG keys with no passphrase so that they can decrypt the file unattended.  That means having root access on a puppet master gives you the ability to find out all our secrets.  That is ok because if you have root access to the puppet master, you own the world anyway.

The secret files are encrypted such that any one key on a list of keys can decrypt them.  That is, when encrypting it is is "encrypted for multiple users".  Each person that should have acecss to the secrets should have a key and be on the key list.  There should also be a key for account that deploys new code to the Puppet master.

What does this look like to the typical sysadmin?
================================

*  If you need to, start the GPG Agent:

``eval $(gpg-agent --daemon)``

*  Decrypt the file so it is editable:

``bin/blackbox_edit_start FILENAME``

(You will need to enter your GPG passphrase.)

*  Edit FILENAME as you desire.

``vim FILENAME``

*  Re-encrypt the file:

``bin/blackbox_edit_end FILENAME``

*  Commit the changes.

```
git commit -a
# or
hg commit
```


This content is released under the MIT License.  See the LICENSE.txt file.

How to use the secrets with Puppet?
================================

### Small strings:

Small strings, such as passwords and API keys, are stored in a hiera yaml file.  You can access them using the hiera() function.

Puppet example for a single password:

```
$the_password = hiera('module::test_password', 'fail')
file {'/tmp/debug-blackbox.txt':
    content => $the_password,
    owner   => 'root',
    group   => 'root',
    mode    => '0600',
}
```

### Entire files:

Entire files, such as SSL certs and private keys, are treated just like files.

Puppet example for an encrypted file:

```
file { '/etc/my_little_secret.key':
    ensure  => 'file',
    owner   => 'root',
    group   => 'puppet',
    mode    => '0760',
    source  => "puppet:///modules/${module_name}/secret_file.key",
}
```


How to enroll a new file into the system?
============================

*  If you need to, start the GPG Agent:

``eval $(gpg-agent --daemon)``

* Add the file to the system:

```
bin/blackbox_register_new_file path/to/file.name.key
```

How to indoctrinate a new user into the system?
============================

``keyrings/live/blackbox-admins.txt`` is a file that
lists which users are able to decrypt files.
(More pedantically, it is a list of the GnuPG key
names that the file is encrypted for.)

To join the list of people that can edit the file requires three steps; You create a GPG key and add it to the key ring.  Then, someone that already has access adds you to the system. Lastly, you should test your access.

### Step 1: YOU create a GPG key pair on a secure machine and add to public keychain.

```
gpg --gen-key
```

Pick defaults for encryption settings, 0 expiration.  Pick a VERY GOOD passphrase.

```
blackbox_addadmin KEYNAME
```
...where "KEYNAME" is the email address listed in the gpg key you created previously. For example:
```
blackbox_addadmin tal@example.com
```

When the command completes successfully, instructions on how to
commit these changes will be output.  Run the command as give.
```
NEXT STEP: Check these into the repo.  Probably with a command like...
git commit -m'NEW ADMIN: tal@example.com' keyrings/live/pubring.gpg keyrings/live/trustdb.gpg keyrings/live/blackbox-admins.txt
```

Role accounts: If you are adding the pubring.gpg of a role account, you can specify the directory where the pubring.gpg file can be found as a 2nd parameter:
```
blackbox_addadmin puppetmaster@puppet-master-1.example.com /path/to/the/dir
```

### Step 2: SOMEONE ELSE adds you to the system.

Ask someone that already has access to re-encrypt the data files. This gives you access.  They simply decrypt and re-encrypt the data without making any changes:

```
gpg --import keyrings/live/pubring.gpg
blackbox_update_all_files
```

Push the re-encrypted files:

```
git commit -a
git push

or

hg commit
hg push
```

### Step 3: YOU test.

Make sure you can decrypt a file.  (Suggestion: Keep a dummy file in VCS just for new people to practice on.)

First Time Setup
===========================

Overview:

To add "blackbox" to a git repo, you'll need to do the following:

  1. Create some directories
  2. For each user, have them create a GPG key and add it to the key ring.
  3. For any automated user (one that must be able to decrypt without a passphrase), create a GPG key and create a subkey with an empty passphrase.
  4. Add

###  Create some directories

You'll want to include blackbox's binaries in your PATH:
```
export PATH=$PATH:/the/path/to/blackbox/bin
```

In the git repo you plan on using blackbox, add these two lines to .gitignore

```
pubring.gpg~
secring.gpg
```

Create this directory.  It is where the pubkeys will be stored:
```
mkdir -p keyrings/live
```

And commit the change:

```
git add keyrings
git add .gitignore
git commit -m'Update .gitignore' .gitignore keyrings
```


### For each user, have them create a GPG key and add it to the key ring.

Follow the instructions for "How to indoctrinate a new user into
the system?".  For the first user, you only have to do Step 1.

Once that is done, is a good idea to test the system by making sure
a file can be added to the system (see "How to enroll a new file
into the system?"), and a different user can decrypt the file.

Make a new file and register it:

```
rm -f foo.txt.gpg foo.txt
echo This is a test. >foo.txt
blackbox_register_new_file foo.txt
```

Decrypt it:

```
blackbox_edit_start foo.txt.gpg 
cat foo.txt
echo This is the new file contents. >foo.txt
```

Re-encrypt it:
```
blackbox_edit_end foo.txt.gpg 
ls -l foo.txt*
```

Push these changes to the repo.  Make sure another user can
check out and change the contents of the file.


### For any automated user create a key and subkey.

An automated user (a "role account") is one that that must be able
to decrypt without a passphrase.  In general you'll want to do this
for the user that pulls the files from the repo to the master.  This
may be automated with Jenkins CI or other CI system.

GPG keys have to have a passphrase. However, passphrases are optional
on subkeys. Therefore, we will create a key with a passphrase then
create a subkey without a passphrase.
Since the subkey is very powerful, it should be created on a very
secure machine.

There's another catch.  The role account probably can't check files
into Git.  It probably only has read-only access to the repo. That's
a good security policy.  This means that the role account can't
be used to upload the subkey public bits into the repo.

Therefore, we will create the key/subkey on a secure machine
as yourself.  From there we can commit the public portions into
the repo.  Also from this account we will export the parts
that the role account needs, copy them to where the role account
can access them, and import them as the role account.

ProTip: If asked to generate entropy, consider running this on the same machine in another window:`sudo dd if=/dev/sda of=/dev/null`

For the rest of this doc, you'll need to make the following substittions:

  - ROLEUSER: svc_deployacct or whatever your role account's name is.
  - NEWMASTER: the machine this role account exists on.
  - SECUREHOST: The machine you use to create the keys. 

NOTE: This should be more automated.  Patches welcome.

On SECUREHOST, create thew puppet master's keys:

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

NOTE: Rather than a real email address, use the username@FQDN of
the host the key will be used on.  If you use this role account on
many machines, each should have its own key.  By using the FQDN of
the host, you will be able to know which key is which.
In this doc, we'll refer to username@FQDN as $KEYNAME

Save the passphrase somewhere safe!

ProTip: If asked to generate entropy, consider running this on the same machine in another window:`sudo dd if=/dev/sda of=/dev/null`

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
$ gpg --homedir . --export -a svc_sadeploy >/tmp/NEWMACHINE/pubkey.txt
$ tar cvf /tmp/keys.tar .
$ rsync -avP /tmp/keys.tar NEWMASTER:/tmp/.
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
$ cd keyrings/live
$ gpg --homedir . --import /tmp/NEWMACHINE/pubkey.txt
```
-->

Back on SECUREHOST, add the new email address to keyrings/live/blackbox-admins.txt:

```
cd /path/to/the/repo
blackbox_addadmin $KEYNAME
```

Verify that secring.gpg is a zero-length file. If it isn't, you have
somehow added a private key to the keyring.  Start over.

```
$ cd keyrings/live
$ ls -l secring.gpg
```

Commit the recent changes:

```
$ cd keyrings/live
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
gpg --import /etc/puppet/keyrings/live/pubring.gpg
export PATH=$PATH:/path/to/blackbox/bin
blackbox_postinstall
sudo -u puppet cat /etc/puppet/hieradata/blackbox.yaml # or any encrypted file.
```

ProTip: If you get "gpg: decryption failed: No secret key" then you forgot to re-encrypt blackbox.yaml with the new key.

On SECUREHOST, securerly delete your files:

```
cd /tmp/NEWMASTER
# On machines with the "shred" command:
shred -u /tmp/keys.tar
find . -type f -print0 | xargs -0 shred -u
# All else:
rm -rf /tmp/NEWMASTER
```

Also shred any other temporary files you may have made.
