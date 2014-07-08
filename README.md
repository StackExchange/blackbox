BlackBox
========

Safely store secrets in Git/Hg for use by Puppet.




Overview
========

The goal is to have secret bits (passwords, private keys, and such) in your VCS repo but encrypted so that
it is safe.  On the puppet masters they sit on disk unencrypted but only readable by Puppet Master.

How does this work?
===================

**Private keys (and anything that is the entire file):**  

Files are kept in git/hg encrypted (foo.txt is stored as foo.txt.gpg).

After deploying an update to your Puppet Master, the master runs a script that decrypts them.  The sit unencrypted on the master, which should already be locked down.

**Passwords (and any short string):**
Passwords are kept in hieradata/blackbox.yaml.gpg, which is decrypted to become hieradata/blackbox.yaml.  This data can be read by hiera.  This file is encrypted/decrypted just like any other blackbox file.

**Key management:**
The Puppet Masters have GPG keys with no passphrase so that they can decrypt the file unattended.  That means having root access on a puppet master gives you the ability to find out all our secrets.  That's ok because if you have root access to the puppet master, you own the world anyway.

The secret files are encrypted such that any one key on a list of keys can decrypt them.  That is, when encrypting it is is "encrypted for multiple users".  Each person that should have acecss to the secrets should have a key and be on the key list.  There should also be a key for account that deploys new code to the Puppet master.

What does this look like to the typical sysadmin?
================================

*  If you need to, start the GPG Agent:

``eval $(gpg-agent --daemon)``

*  Decrypt so you can edit:

``bin/blackbox_edit_start.sh FILENAME``

This decrypts the data. (You will need to enter your GPG passphrase.)

*  Edit FILENAME as you desire.

``vim FILENAME``

*  Re-encrypt the file.

``bin/blackbox_edit_end.sh FILENAME``

Encrypts the data.

*  Commit the changes.

``git commit -a``
or
``hg commit``


This content is released under the MIT License.  See the LICENSE.txt file.

