How to use the secrets with Ansible?
===================================

Ansible Vault provides functionality for encrypting both entire files
and strings stored within files; however, keeping track of the
password(s) required for decryption is not handled by this module.

Instead one must specify a password file when running the playbook.

Ansible example for password file: `my_secret_password.txt.gpg`

```
ansible-playbook --vault-password-file my_secret_password.txt site.yml
```

Alternatively, one can specify this in the
`ANSIBLE_VAULT_PASSWORD_FILE` environment variable.

