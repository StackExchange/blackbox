How to use the secrets with Puppet?
===================================

# Entire files:

Entire files, such as SSL certs and private keys, are treated just
like regular files. You decrypt them any time you push a new release
to the puppet master.

Example of an encrypted file named `secret_file.key.gpg`

* Plaintext file is: `modules/${module_name}/files/secret_file.key`
* Encrypted file is: `modules/${module_name}/files/secret_file.key.gpg`
* Puppet sees it as: `puppet:///modules/${module_name}/secret_file.key`

Puppet code that stores `secret_file.key` in `/etc/my_little_secret.key`:

```
file { '/etc/my_little_secret.key':
    ensure  => 'file',
    owner   => 'root',
    group   => 'puppet',
    mode    => '0760',
    source  => "puppet:///modules/${module_name}/secret_file.key",  # No ".gpg"
}
```

# Small strings:

For small strings such as passwords and API keys, it makes sense
to store them in an (encrypted) YAML file which is then made
available via hiera.

For example, we use a file called `blackbox.yaml`. You can access the
data in it using the hiera() function.

*Setup:*

Edit `hiera.yaml` to include "blackbox" to the search hierarchy:

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

The variable `$the_password` will contain "my secret password" and can be used anywhere strings are used.
