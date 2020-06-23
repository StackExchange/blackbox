BlackBox
========

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

* [Overview](user-overview)
* [Why is this important?](why-is-this-important)
* [Support/Community](support)
* [How BB encrypts](encryption)
* [OS Compatibility](compatibility)
* [Installation Instructions](installation)
* [Alternatives](alternatives)

User Info:

* [Enabling Blackbox on a Repo](enable-repo)
* [Enroll a file](enable-repo)
* [Full Command List](full-command-list)
* [Add/Remove users](admin-ops)
* [Add/Remove files](file-ops)
* [Advanced techiques](advanced)
* [Use with Role Accounts](role-accounts)
* [Backwards Compatibility](backwards-compatibility)
* [Replacing expired keys](expired-keys)
* [Git Tips](git-tips)
* [SubVersion Tips](subversion-tips)
* [GnuPG tips](gnupg-tips)
* [Use with Ansible](with-ansible)
* [Use with Puppet](with-puppet)

For contributors:

* [Developer Info](dev)
* [Code overview](dev-code-overview)
* [Add new OS support]()
* [Add new VCS support]()


A slide presentation about an older release [is on SlideShare](http://www.slideshare.net/TomLimoncelli/the-blackbox-project-sfae).

Join our mailing list: [https://groups.google.com/d/forum/blackbox-project](https://groups.google.com/d/forum/blackbox-project)


License
=======

This content is released under the MIT License.
See the [LICENSE.txt](LICENSE.txt) file.
