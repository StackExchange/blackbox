BlackBox v2
===========

WARNING: v2 is still experimental.  It is in the same git repo as v1
because the filenames do not overlap.  Please do not mix the two.  v1
is in `bin`.  v2 is in `cmd/blackbox` and `binv2`.

Blackbox is an open source tool that enables you to safe store sensitive information in
Git (or other) repos by encrypting them with GPG.  Only the encrypted
version of the file is available.  You can be free to provide access
to the repo, as but only people with the right GPG keys can access the
encrypted data.

Things you should **never** store in a repo without encryption:

* TLS (SSL) certificates
* Passwords
* API keys
* And more!

Project Info:

* [Overview](user-overview.md)
* [Why is this important?](why-is-this-important.md)
* [Support/Community](support.md)
* [How BB encrypts](encryption.md)
* [OS Compatibility](compatibility.md)
* [Installation Instructions](installation.md)
* [Alternatives](alternatives.md)

User Info:

* [Enabling Blackbox on a Repo](enable-repo.md)
* [Enroll a file](enable-repo.md)
* [Full Command List](full-command-list.md)
* [Add/Remove users](admin-ops.md)
* [Add/Remove files](file-ops.md)
* [Advanced techiques](advanced.md)
* [Use with Role Accounts](role-accounts.md)
* [Backwards Compatibility](backwards-compatibility.md)
* [Replacing expired keys](expired-keys.md)
* [Git Tips](git-tips.md)
* [SubVersion Tips](subversion-tips.md)
* [GnuPG tips](gnupg-tips.md)
* [Use with Ansible](with-ansible.md)
* [Use with Puppet](with-puppet.md)

For contributors:

* [Developer Info](dev.md)
* [Code overview](dev-code-overview.md)
* [HOWTO: Add new OS support](dev-add-os-support.md)
* [HOWTO: Add new VCS support](dev-add-vcs-support.md)


A slide presentation about an older release [is on SlideShare](http://www.slideshare.net/TomLimoncelli/the-blackbox-project-sfae).

Join our mailing list: [https://groups.google.com/d/forum/blackbox-project](https://groups.google.com/d/forum/blackbox-project)


License
=======

This content is released under the MIT License.
See the [LICENSE.txt](LICENSE.txt) file.
