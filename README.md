BlackBox
========

Safely store secrets in a VCS repo (i.e. Git, Mercurial, Subversion or Perforce). These
commands make it easy
for you to Gnu Privacy Guard (GPG) encrypt specific files in a repo so they are
"encrypted at rest" in your repository. However, the scripts
make it easy to decrypt them when you need to view or edit them,
and decrypt them for use in production. Originally written
for Puppet, BlackBox now works with any Git or Mercurial repository.

A slide presentation about an older release [is on SlideShare](http://www.slideshare.net/TomLimoncelli/the-blackbox-project-sfae).

Table of Contents
===============

* [Overview](#overview)
* [Why is this important?](#why-is-this-important)
* [Installation Instructions](#installation-instructions)
* [Commands](#commands)
* [Compatibility](#compatibility)
* [How is the encryption done?](#how-is-the-encryption-done)
* [What does this look like to the typical user?](#what-does-this-look-like-to-the-typical-user)
* [How to use the secrets with Puppet?](#how-to-use-the-secrets-with-puppet)
* [How to enroll a new file into the system?](#how-to-enroll-a-new-file-into-the-system)
* [How to remove a file from the system?](#how-to-remove-a-file-from-the-system)
* [How to indoctrinate a new user into the system?](#how-to-indoctrinate-a-new-user-into-the-system)
* [How to remove a user from the system?](#how-to-remove-a-user-from-the-system)
* [Enabling Blackbox For a Repo](#enabling-blackbox-for-a-repo)
* [Set up automated users or "role accounts"](#set-up-automated-users-or-role-accounts)
* [Some common errors](#some-common-errors)
* [Using Blackbox without a repo](#using-blackbox-without-a-repo)
* [How to submit bugs or ask questions?](#how-to-submit-bugs-or-ask-questions)
* [Developer Info](#developer-info)
* [Alternatives](#alternatives)
* [License](#license)


Overview
========

Suppose you have a VCS repository (i.e. a Git or Mercurial repo)
and certain files contain secrets such as passwords or SSL private
keys.  Often people just store such files "and hope that nobody finds
them in the repo".  That's not safe.

With BlackBox, those files are stored encrypted using GPG. Access to
the VCS repo without also having the right GPG keys
makes it worthless to have the files.  As long as you keep your GPG keys
safe, you don't have to worry about storing your VCS repo on an untrusted
server.  Heck, even if you trust your server, now you don't have to trust
the people that do backups of that server, or the people that handle the
backup tapes!

Rather than one GPG passphrase for all the files, each person with access
has their own GPG keys in the system.  Any file can be decrypted by
anyone with their GPG key.  This way, if one person leaves the company,
you don't have to communicate a new password to everyone with access.
Simply disable the one key that should no longer have access.
The process for doing this is as easy as running 2 commands (1 to
disable their key, 1 to re-encrypt all files.)

Automated processes often need access to all the decrypted files.
This is easy too.  For example, suppose Git is being used for Puppet
files.  The master needs access to the decrypted version of all the
files.  Simply set up a GPG key for the Puppet master (or the role
account that pushes new files to the Puppet master) and have that
user run `blackbox_postdeploy` after any files are updated.

Getting started is easy.  Just `cd` into a Git, Mercurial, Subversion or
Perforce repository and run `blackbox_initialize`.  After that, if a file is to
be encrypted, run `blackbox_register_new_file` and you are done.  Add
and remove keys with `blackbox_addadmin` and `blackbox_removeadmin`.
To view and/or edit a file, run `blackbox_edit`; this will decrypt the
file and open with whatever is specified by your $EDITOR environment
variable.  When you close the editor the file will automatically be
encrypted again and the temporary plaintext file will be shredded.  If
you need to leave the file decrypted while you update you can use the
`blackbox_edit_start` to decrypt the file and  `blackbox_edit_end` when
you want to "put it back in the box."


Why is this important?
============================

OBVIOUSLY we don't want secret things like SSL private keys
and passwords to be leaked.

NOT SO OBVIOUSLY when we store "secrets" in a VCS repo like Git or
Mercurial, suddenly we are less able to share our code with other
people.  Communication between subteams of an organization is hurt.
You can't collaborate as well.  Either you find yourself emailing
individual files around (yuck!), making a special repo with just
the files needed by your collaborators (yuck!!), or just deciding that
collaboration isn't worth all that effort (yuck!!!).

The ability to be open and transparent about our code, with the
exception of a few specific files, is key to the kind of
collaboration that DevOps and modern IT practitioniers
need to do.

Installation Instructions:
==========================

  * *The MacPorts Way*: `sudo port install vcs_blackbox`
  * *The RPM way*: Check out the repo and make an RPM via `make packages-rpm`; now you can distribute the RPM via local methods.
  * *The Debian/Ubuntu way*: Check out the repo and install [fpm](https://github.com/jordansissel/fpm). Now you can make a DEB `make packages-deb` that can be distributed via local methods.
  * *The hard way*: Copy all the files in "bin" to your "bin".
  * *The manual way*: `make manual-install` to install. `make manual-uninstall` to uninstall.
  * *The Antigen Way*: Add `antigen bundle StackExchange/blackbox` to your .zshrc
  * *The Zgen Way*: Add `zgen load StackExchange/blackbox` to your .zshrc where you're loading your other plugins.


Commands:
============================


| Name: | Description: |
| --- | --- |
| `blackbox_edit` | Decrypt, run $EDITOR, re-encrypt a file |
| `blackbox_edit_start` | Decrypt a file so it can be updated |
| `blackbox_edit_end` | Encrypt a file after blackbox_edit_start was used |
| `blackbox_cat` | Decrypt and view the contents of a file |
| `blackbox_diff` | Diff decrypted files against their original crypted version |
| `blackbox_initialize` | Enable blackbox for a GIT or HG repo |
| `blackbox_register_new_file` | Encrypt a file for the first time |
| `blackbox_deregister_file` | Remove a file from blackbox |
| `blackbox_list_files` | List the files maintained by blackbox |
| `blackbox_decrypt_all_files` | Decrypt all managed files (INTERACTIVE) |
| `blackbox_postdeploy` | Decrypt all managed files (batch) |
| `blackbox_addadmin` | Add someone to the list of people that can encrypt/decrypt secrets |
| `blackbox_removeadmin` | Remove someone from the list of people that can encrypt/decrypt secrets |
| `blackbox_shred_all_files` | Safely delete any decrypted files |
| `blackbox_update_all_files` | Decrypt then re-encrypt all files. Useful after keys are changed |
| `blackbox_whatsnew` | show what has changed in the last commit for a given file |


Compatibility:
============================

Blackbox automatically determines which VCS you are using
and does the right thing.  It has a plug-in architecture
to make it easy to extend to work with other systems.
It has been tested to work with many operating systems.

* Version Control systems
  * `git` -- The Git
  * `hg` -- Mercurial
  * `svn` -- SubVersion (Thanks, Ben Drasin!)
  * `p4` -- Perforce
  * none -- The files can be decrypted outside of a repo if the keyrings directory is intact
* Operating system
  * CentOS / RedHat
  * MacOS X
  * Cygwin (Thanks, Ben Drasin!)

To add or fix support for a VCS system, look for code at the end
of `bin/_blackbox_common.sh`

To add or fix support for a new operating system, look for the case
statements in `bin/_blackbox_common.sh` and `bin/_stack_lib.sh` and
maybe `tools/confidence_test.sh`

Note: Cywin support requires the following packages:

  * Normal operation:
    * gnupg
    * git or mercurial or subversion or perforce (as appropriate)
  * Development (if you will be adding code and want to run the confidence test)
    * procps
    * make
    * git (the confidence test currently only tests git)


How is the encryption done?
============================

GPG has many different ways to encrypt a file.  BlackBox uses
the mode that lets you specify a list of keys that can decrypt
the messsage.

If you have 5 people ("admins") that should be able to access
the secrets, each creates a GPG key and adds their public key
to the keychain.  The GPG command used to encrypt the file lists
all 5 key names, and therefore any 1 key can decrypt the file.

To remove someone's access, remove that admin's key name (i.e. email
address) from the list of admins and re-encrypt all the files.
They can still read the .gpg file (assuming they have access
to the repository) but they can't decrypt it any more.

*What if they kept a copy of the old repo before you removed
access?*  Yes, they can decrypt old versions of the file. This
is why when an admin leaves the team, you should change all
your passwords, SSL certs, and so on.  You should have been
doing that before BlackBox, right?

*Why don't you use symmetric keys?*  In other words, why mess
with all this GPG key stuff and instead why don't we just encrypt
all the files with a single passphrase.  Yes, GPG supports that,
but then we are managing a shared password, which is fraught with problems.
If someone "leaves the team" we would have to communicate to everyone
a new password. Now we just have to remove their key.  This scales
better.

*How do automated processes decrypt without asking for a password?*
GPG requires a passphrase on a private key.  However, it permits
the creation of subkeys that have no passphrase.  For automated
processes, create a subkey that is only stored on the machine
that needs to decrypt the files. For example, at Stack Exchange,
when our Continuous Integration (CI) system pushes
a code change to our Puppet masters, they run `blackbox_postdeploy`
to decrypt all the files.  The user that runs this code has a subkey
that doesn't require a passphrase. Since we have many masters,
each has its own key.  And, yes, this means our Puppet Masters
have to be very secure.  However, they were already secure because,
like, dude... if you can break into someone's puppet master you own
their network.

*If you use Puppet, why didn't you just use hiera-eyaml?*
There are 4 reasons:

  1. This works with any Git or Mercurial repo, even if you aren't using Puppet.
  2. hiera-eyaml decrypts "on demand" which means your Puppet Master now uses a lot of CPU to decrypt keys every time it is contacted. It slows down your master, which, in my case, is already slow enough.
  3. This works with binary files, without having to ASCIIify them and paste them into a YAML file.  Have you tried to do this with a cert that is 10K long and changes every few weeks?  Ick.
  4. hiera-eyaml didn't exist when I wrote this.


What does this look like to the typical user?
================================

*  If you need to, start the GPG Agent: `eval $(gpg-agent --daemon)`
*  Decrypt the file so it is editable: `blackbox_edit_start FILENAME`
*  (You will need to enter your GPG passphrase.)
*  Edit FILENAME as you desire: `vim FILENAME`
*  Re-encrypt the file: `blackbox_edit_end FILENAME`
*  Commit the changes.  `git commit -a` or `hg commit`

Wait... it can be even easier than that!
Run `blackbox_edit FILENAME`, and it'll decrypt the file
in a temp file and call `$EDITOR` on it, re-encrypting again after the editor
is closed.


How to use the secrets with Puppet?
================================

### Entire files:

Entire files, such as SSL certs and private keys, are treated just like
regular files.  You decrypt them any time you push a new release
to the puppet master.

Puppet example for an encrypted file: `secret_file.key.gpg`

```
file { '/etc/my_little_secret.key':
    ensure  => 'file',
    owner   => 'root',
    group   => 'puppet',
    mode    => '0760',
    source  => "puppet:///modules/${module_name}/secret_file.key",
}
```


### Small strings:

Small strings, such as passwords and API keys, are stored in a hiera
yaml file, which you encrypt with `blackbox_register_new_file`. For
example, we use a file called `blackbox.yaml`.  You can access them
using the hiera() function.

*Setup:* Configure `hiera.yaml` by adding "blackbox" to the search hierarchy:

```
:hierarchy:
  - ...
  - blackbox
  - ...
```

In blackbox.yaml specify:

```
---
module::test_password: "my secret password"
```

In your Puppet Code, access the password as you would any hiera data:

```
$the_password = hiera('module::test_password', 'fail')

file {'/tmp/debug-blackbox.txt':
    content => $the_password,
    owner   => 'root',
    group   => 'root',
    mode    => '0600',
}
```

The variable `$the_password` will contain "my secret password" and
can be used anywhere strings are used.


How to enroll a new file into the system?
============================

*  If you need to, start the GPG Agent: `eval $(gpg-agent --daemon)`
*  Add the file to the system:
```
blackbox_register_new_file path/to/file.name.key
```

Multiple file names can be specified on the command line:

Example 1: Register 2 files:

    blackbox_register_new_file file1.txt file2.txt

Example 2: Register all the files in `$DIR`:

    find $DIR -type f -not -name '*.gpg' -print0 | xargs -0 blackbox_register_new_file


How to remove a file from the system?
============================

This happens quite rarely, but we've got it covered:
```
blackbox_deregister_file path/to/file.name.key
```

How to indoctrinate a new user into the system?
============================

``keyrings/live/blackbox-admins.txt`` is a file that
lists which users are able to decrypt files.
(More pedantically, it is a list of the GnuPG key
names that the file is encrypted for.)

To join the list of people that can edit the file requires three steps; You create a GPG key and add it to the key ring.  Then, someone that already has access adds you to the system. Lastly, you should test your access.

### Step 1: YOU create a GPG key pair on a secure machine and add to public keychain.

If you don't already have a GPG key, here's how to generate one:

```
gpg --gen-key
```

Pick defaults for encryption settings, 0 expiration.  Pick a VERY GOOD passphrase.  Store the private key securely.  Tip: Store it on a secure machine, or one with little or no internet access, with full-disk-encryption, etc.  Your employer problably has rules about how to store such things.

Now that you have a GPG key, add yourself as an admin:

```
blackbox_addadmin KEYNAME
```

...where "KEYNAME" is the email address listed in the gpg key you created previously. For example:

```
blackbox_addadmin tal@example.com
```

When the command completes successfully, instructions on how to commit these changes will be output.  Run the command as given to commit the changes.  It will look like this:

```
git commit -m'NEW ADMIN: tal@example.com' keyrings/live/pubring.gpg keyrings/live/trustdb.gpg keyrings/live/blackbox-admins.txt
```

Then push it to the repo:

```
git push

or

ht push

(or whatever is appropriate)
```

NOTE: Creating a Role Account?  If you are adding the pubring.gpg of a role account, you can specify the directory where the pubring.gpg file can be found as a 2nd parameter: `blackbox_addadmin puppetmaster@puppet-master-1.example.com /path/to/the/dir`


### Step 2: SOMEONE ELSE adds you to the system.

Ask someone that already has access to re-encrypt the data files. This gives you access.  They simply decrypt and re-encrypt the data without making any changes.

Pre-check: Verify the new keys look good.

```
$ gpg --homedir=keyrings/live --list-keys 
```

For example, examine the key name (email address) to make sure
it conforms to corporate standards.

Import the keychain into your personal keychain and reencrypt:

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

How to remove a user from the system?
============================

Simply run `blackbox_removeadmin` with their keyname then re-encrypt:

Example:

```
blackbox_removeadmin olduser@example.com
blackbox_update_all_files
```

When the command completes, you will be given a reminder to check in the change and push it.

Note that their keys will still be in the key ring, but they will
go unused.  If you'd like to clean up the keyring, use the normal
GPG commands and check in the file.

```
gpg --homedir=keyrings/live --list-keys
gpg --homedir=keyrings/live --delete-key olduser@example.com
git commit -m'Cleaned olduser@example.com from keyring'  keyrings/live/*
```

The key ring only has public keys. There are no secret keys to delete.

Remember that this person did have access to all the secrets at one
time. They could have made a copy.  Therefore, to be completely
secure, you should change all passwords, generate new SSL keys, and
so on just like when anyone that had privileged access leaves an
organization.


Enabling Blackbox For a Repo
===========================

Overview:

To add "blackbox" to a git or mercurial repo, you'll need to do the following:

  1. Run the initialize script.  This adds a few files to your repo in a directory called "keyrings".
  2. For the first user, create a GPG key and add it to the key ring.
  3. Encrypt the files you want to be "secret".
  4. For any automated user (one that must be able to decrypt without a passphrase), create a GPG key and create a subkey with an empty passphrase.

###  Run the initialize script.

You'll want to include blackbox's "bin" directory in your PATH:

```
export PATH=$PATH:/the/path/to/blackbox/bin
blackbox_initialize
```

If you're using antigen, adding `antigen bundle StackExchange/blackbox` to
your .zshrc will download this repository and add it to your $PATH.

### For the first user, create a GPG key and add it to the key ring.

Follow the instructions for "How to indoctrinate a new user into
the system?".  Only do Step 1.

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


Set up automated users or "role accounts"
=========================================

i.e. This is how a Puppet Master can have access to the unencrypted data.


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
into Git/Mercurial.  It probably only has read-only access to the repo. That's
a good security policy.  This means that the role account can't
be used to upload the subkey public bits into the repo.

Therefore, we will create the key/subkey on a secure machine
as yourself.  From there we can commit the public portions into
the repo.  Also from this account we will export the parts
that the role account needs, copy them to where the role account
can access them, and import them as the role account.

ProTip: If asked to generate entropy, consider running this on the
same machine in another window: `sudo dd if=/dev/sda of=/dev/null`

For the rest of this doc, you'll need to make the following substitutions:

  - ROLEUSER: svc_deployacct or whatever your role account's name is.
  - NEWMASTER: the machine this role account exists on.
  - SECUREHOST: The machine you use to create the keys.

NOTE: This should be more automated/scripted.  Patches welcome.

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

NOTE: Rather than a real email address, use the username@FQDN of
the host the key will be used on.  If you use this role account on
many machines, each should have its own key.  By using the FQDN of
the host, you will be able to know which key is which.
In this doc, we'll refer to username@FQDN as $KEYNAME

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
$ gpg --homedir . --export -a svc_sadeploy >/tmp/NEWMASTER/pubkey.txt
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
$ gpg --homedir . --import /tmp/NEWMASTER/pubkey.txt
```
-->

Back on SECUREHOST, add the new email address to keyrings/live/blackbox-admins.txt:

```
cd /path/to/the/repo
blackbox_addadmin $KEYNAME /tmp/NEWMASTER
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
blackbox_postdeploy
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


Some common errors:
=========================================

`gpg: filename: skipped: No public key` -- Usually this means there
is an item in `keyrings/live/blackbox-admins.txt` that is not the
name of the key.  Either something invalid was inserted (like a
filename instead of a username) or a user has left the organization
and their key was removed from the keychain, but their name wasn't
removed from the blackbox-admins.txt file.

`gpg: decryption failed: No secret key` -- Usually means you forgot
to re-encrypt the file with the new key.


Using Blackbox without a repo
===========================

If the files are copied out of a repo they can still be decrypted
and edited.  Obviously edits, changes to keys, and such will be lost
if they are made outside the repo. Also note that commands are most
likely to only work if run from the base directory (i.e. the parent to
the keyrings directory).

The following commands have been tested outside a repo:

  * `blackbox_postdeploy`
  * `blackbox_edit_start`
  * `blackbox_edit_end`


How to submit bugs or ask questions?
============

We welcome questions, bug reports and feedback!

  * https://github.com/StackExchange/blackbox/issues

Developer Info
============

Code submissions are gladly welcomed!  The code is
fairly easy to read.


Get the code:

```
git clone git@github.com:StackExchange/blackbox.git
```

Test your changes:

```
make confidence
```

This runs through a number of system tests. It
creates a repo, encrypts files, decrypts files, and so on.
You can run these tests to verify that the changes you made
didn't break anything.  You can also use these tests to
verify that the system works with a new operating system.

Please submit tests with code changes:

The best way to change Blackbox is via Test Driven Development.
First add a test to `tools/confidence.sh`. This test should
fail, and demonstrate the need for the change you are about to
make. Then fix the bug or add the feature you want.  When
you are done, `make confidence` should pass all tests.
The PR you submit should include your code as well as the new
test.  This way the confidence tests accumulate as the system
grows as we know future changes don't break old features.

Note: The tests currently assume "git" and have been tested
only on CentOS, Mac OS X, and Cygwin.  Patches welcome!


Alternatives
============

Here are other open source packages that do something similar to Blackbox. If you like them better than Blackbox, please use them.

  * git-crypt:  https://www.agwa.name/projects/git-crypt/
  * Pass: http://www.zx2c4.com/projects/password-store/
  * Transcrypt: https://github.com/elasticdog/transcrypt

git-crypt has the best git integration.  Once set up it is nearly transparent to the users.  However it only works with git.


License
=======
This content is released under the MIT License.  See the LICENSE.txt file.
