User Guide
==========

# Overview

Suppose you have a VCS repository (i.e. a Git or Mercurial repo) and
certain files contain secrets such as passwords or SSL private keys.
Often people just store such files "and hope that nobody finds them in
the repo". That's not safe.  Hope is not a strategy.

With BlackBox, those files are stored encrypted using GPG. Access to
the repo without also having the right GPG keys makes those files as worthless
as random bits. As long as you keep your GPG keys safe, you don't
have to worry about storing your VCS repo on an untrusted server or
letting anyone clone the repo.

Heck, even if you trust your server, now you don't have to trust the
people that do backups of that server!

Each person ("admin") of the system can decrypt all the files using
their GPG key, which has its own passphrase.  The authorized GPG keys
can decrypt any file.  This is better than systems that use one
GPG key (and passphrase) that must be shared among a group of people.
It is much better than having one passphrase for each file (I don't
think anyone actually does that).

Since any admin's GPG key can decrypt the files, if one person leaves
the company, you don't have to communicate a new passphrase to everyone.
Simply disable the one key that should no longer have access.
The process for doing this is as easy as running 2 commands (1 to
disable their key, 1 to re-encrypt all files.)  Obviously if they kept
a copy of the repo (and their own passphrase) before leaving the
company, they have access to the secrets. However, you should rotate
those secrets anyway. ("rotate secrets" means changing the passwords,
regenerating TLS certs, and so on).

# Sample session:

First we are going to list the files currently in the blackbox. In
this case, it is an SSH private key.

```
$ blackbox file list
modules/log_management/files/id_rsa
```

Excellent! Our coworkers have already registered a file with the
system.  Let's decrypt it, edit it, and re-encrypt it.

```
$ blackbox decrypt modules/log_management/files/id_rsa
========== DECRYPTING "modules/log_management/files/id_rsa"
$ vi modules/log_management/files/id_rsa
```

That was easy so far!

When we encrypt it, Blackbox will not commit the changes, but it
will give a hint that you should. It spells out the exact command you
need to type and even proposes a commit message.

```
$ blackbox encrypt modules/log_management/files/id_rsa
========== ENCRYPTING "modules/log_management/files/id_rsa"

NEXT STEP: You need to manually check these in:
     git commit -m"ENCRYPTED modules/log_management/files/id_rsa" modules/log_management/files/id_rsa.gpg
```

You can also use `blackbox edit <filename>` to decrypt a file, edit it
(it will call `$EDITOR`) and re-encrypt it.


Now let's register a new file with the blackbox system.
`data/pass.yaml` is a small file that stores a very important
password.  In this example, we had just stored the unecrypted
password in our repo. That's bad.  Let's encrypt it.

```
$ blackbox file add data/pass.yaml
========== SHREDDING ("/bin/rm", "-f"): "data/pass.yaml"

NEXT STEP: You need to manually check these in:
     git commit -m"NEW FILES: data/pass.yaml" .gitignore keyrings/live/blackbox-files.txt modules/stacklb/pass.yaml modules/stacklb/pass.yaml.gpg
```

Before we commit the change, let's do a `git status` to see what else
has changed.

```
$ git status
On branch master
Changes to be committed:
  (use "git restore --staged <file>..." to unstage)
	modified:   .gitignore
	modified:   keyrings/live/blackbox-files.txt
	deleted:    modules/stacklb/pass.yaml
	new file:   modules/stacklb/pass.yaml.gpg

```

Notice that a number of files were modified:

* `.gitignore`: This file is updated to include the plaintext
  filename, so that you don't accidentally add it to the repo in the
  future.
* `.blackbox/blackbox-files.txt`: The list of files that are registered with the system.
* `data/pass.yaml`: The file we encrypted is deleted from the repo.
* `data/pass.yaml.gpg`: The encrypted file is added to the repo.

Even though pass.yaml was deleted from the repo, it is still in the
repo's history. Anyone with an old copy of the repo, or a new copy
that knows how to view the repo's history, can see the secret
password.  For that reason, you should change the password and
re-encrypt the file.  This is an important point.  Blackbox is not
magic and it doesn't have a "Men In Black"-style neuralizer that
can make people forget the past.  If someone leaves a project, you
have to change the old passwords, etc.

Those are the basics.  Your next step might be:

* TODO: How to enable Blackbox for a repo.
* TODO: How to add yourself as an admin to a repo.
* TODO: Complete list of [all blackbox commands](all-commands)
