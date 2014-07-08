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

*  Decrypt the file so it is editable:

``bin/blackbox_edit_start.sh FILENAME``

(You will need to enter your GPG passphrase.)

*  Edit FILENAME as you desire.

``vim FILENAME``

*  Re-encrypt the file:

``bin/blackbox_edit_end.sh FILENAME``


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
bin/blackbox_register_new_file.sh path/to/file.name.key
```

How do to indoctrinate a new user into the system?
============================

``keyrings/live/blackbox-admins.txt`` is a file that
lists which users are able to decrypt files.
(More pedantically, it is a list of the GnuPG key
names that the file is encrypted for.)

To join the list of people that can edit the file requires three steps; You create a GPG key and add it to the key ring.  Then, someone that already has access adds you to the system. Lastly, you should test your access.

### Step 1: YOU create a GPG key pair on a secure machine and add to public keychain.


```
KEYNAME=$USER@$DOMAINNAME
gpg --gen-key
```

Pick defaults for encryption settings, 0 expiration.  Pick a VERY GOOD passphrase.

When GPG is generating entropy, consider running this on the machine in another window:

```
dd if=/dev/sda of=/dev/null
```

Add your public key to the public key-ring.

```
gpg --export -a $KEYNAME >~/.gnupg/pubkey.txt
wc -l ~/.gnupg/pubkey.txt
```

The output of "wc" should be non-zero (usually it is 30 or more)

Add your keyname to the list of keys:

```
cd keyrings/live
gpg --homedir=. --import ~/.gnupg/pubkey.txt
echo $KEYNAME >>blackbox-admins.txt
sort  -fdu -o blackbox-admins.txt <(echo $KEYNAME) blackbox-admins.txt
```

Check all these updates into the VCS:

```
git commit -m"Adding my gpg key" pubring.gpg trustdb.gpg blackbox-admins.txt

or

hg commit -m"Adding my gpg key" pubring.gpg trustdb.gpg blackbox-admins.txt
```


### Step 2: SOMEONE ELSE adds you to the system.

Ask someone that already has access to re-encrypt the data files. This gives you access.  They simply decrypt and re-encrypt the data without making any changes:

```
gpg --import keyrings/live/pubring.gpg
bin/blackbox_update_all_files.sh
```

Push the re-encrypted files:

```
git push

or

hg push
```

### Step 3: YOU test.

Make sure you can decrypt a file.  (NOTE: It is a good idea to keep a dummy file in VCS just for new people to practice on.)


Setting up the Puppet Master:
===========================

Whatever user that pushes code updates to the Puppet master must (1) have a GPG key with no pass phrase, (2) run the ``bin/blackbox_postinstall.sh`` script after new code is pushed.

(docs coming soon.)

Setting up hiera:
=================

(docs coming soon)
