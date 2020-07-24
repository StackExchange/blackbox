How is the encryption done?
===========================

GPG has many different ways to encrypt a file. BlackBox uses the mode
that lets you specify a list of keys that can decrypt the message.

If you have 5 people ("admins") that should be able to access the
secrets, each creates a GPG key and adds their public key to the
keychain. The GPG command used to encrypt the file lists all 5 key
names, and therefore any 1 key can decrypt the file.

Blackbox stores a copy of the public keys of all admins. It never
stores the private keys.

To remove someone's access, remove that admin's key name (i.e. email
address) from the list of admins and re-encrypt all the files. They
can still read the .gpg file (assuming they have access to the
repository) but they can't decrypt it any more.

*What if they kept a copy of the old repo before you removed access?*
Yes, they can decrypt old versions of the file. This is why when an
admin leaves the team, you should change all your passwords, SSL
certs, and so on. You should have been doing that before BlackBox,
right?

*Why don't you use symmetric keys?* In other words, why mess with all
this GPG key stuff and instead why don't we just encrypt all the files
with a single passphrase. Yes, GPG supports that, but then we are
managing a shared password, which is fraught with problems. If someone
"leaves the team" we would have to communicate to everyone a new
password. Now we just have to remove their key. This scales better.

*How do automated processes decrypt without asking for a password?*
GPG requires a passphrase on a private key. However, it permits the
creation of subkeys that have no passphrase. For automated processes,
create a subkey that is only stored on the machine that needs to
decrypt the files. For example, at Stack Exchange, when our Continuous
Integration (CI) system pushes a code change to our Puppet masters,
they run `blackbox decrypt --all --overwrite` to decrypt all the files.
The user that
runs this code has a subkey that doesn't require a passphrase. Since
we have many masters, each has its own key. And, yes, this means our
Puppet Masters have to be very secure. However, they were already
secure because, like, dude... if you can break into someone's puppet
master you own their network.

*If you use Puppet, why didn't you just use hiera-eyaml?* There are 4
reasons:

1. This works with any Git or Mercurial repo, even if you aren't using Puppet.
2. hiera-eyaml decrypts "on demand" which means your Puppet Master now uses a lot of CPU to decrypt keys every time it is contacted. It slows down your master, which, in my case, is already slow enough.
3. This works with binary files, without having to ASCIIify them and paste them into a YAML file. Have you tried to do this with a cert that is 10K long and changes every few weeks? Ick.
4. hiera-eyaml didn't exist when I wrote this. (That's the real reason.)

