#!/usr/bin/env bash

export PATH=/home/tlimoncelli/gitwork/blackbox/bin:/usr/lib64/qt-3.3/bin:/usr/local/bin:/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/sbin

. _stack_lib.sh

set -e

function PHASE() {
  echo '===================='
  echo '===================='
  echo '=========' """$@"""
  echo '===================='
  echo '===================='
}

function assert_file_missing() {
  if [[ -e "$1" ]]; then
    echo "ASSERT FAILED: ${1} should not exist."
    exit 1
  fi
}

function assert_file_exists() {
  if [[ ! -e "$1" ]]; then
    echo "ASSERT FAILED: ${1} should exist."
    exit 1
  fi
}
function assert_file_md5hash() {
  local file="$1"
  local wanted="$2"
  local found=$(md5sum <"$file" | cut -d' ' -f1 )
  assert_file_exists "$file"
  if [[ "$wanted" != "$found" ]]; then
    echo "ASSERT FAILED: $file hash wanted=$desired found=$found"
    exit 1
  fi
}

make_tempdir test_repository
cd "$test_repository"

make_self_deleting_tempdir fake_alice_home
make_self_deleting_tempdir fake_bob_home
export GNUPGHOME="$fake_alice_home"
eval $(gpg-agent --homedir "$fake_alice_home" --daemon)
GPG_AGENT_INFO_ALICE="$GPG_AGENT_INFO"

export GNUPGHOME="$fake_bob_home"
eval $(gpg-agent --homedir "$fake_alice_home" --daemon)
GPG_AGENT_INFO_BOB="$GPG_AGENT_INFO"

function become_alice() {
  export GNUPGHOME="$fake_alice_home"
  export GPG_AGENT_INFO="$GPG_AGENT_INFO_ALICE"
  echo BECOMING ALICE: GNUPGHOME=$GNUPGHOME AGENT=$GPG_AGENT_INFO
  git config --global user.name "Alice Example"
  git config --global user.email alice@example.com
}

function become_bob() {
  export GNUPGHOME="$fake_alice_home"
  export GPG_AGENT_INFO="$GPG_AGENT_INFO_ALICE"
  git config --global user.name "Bob Example"
  git config --global user.email bob@example.com
}


PHASE 'Alice creates a repo.  She creates secret.txt.'

become_alice
git init
echo 'this is my secret' >secret.txt


PHASE 'Alice wants to be part of the secret system.'
PHASE 'She creates a GPG key...'

make_self_deleting_tempfile gpgconfig
cat >"$gpgconfig" <<EOF
%echo Generating a basic OpenPGP key
Key-Type: default
Subkey-Type: default
Name-Real: Alice Example
Name-Comment: with weak passphrase
Name-Email: alice@example.com
Expire-Date: 0
Passphrase: a
# Do a commit here, so that we can later print "done" :-)
%commit
%echo done
EOF
gpg --no-permission-warning --batch --gen-key "$gpgconfig"

#gpg --delete-key bob@example.com || true
#gpg --delete-key alice@example.com || true


PHASE 'Initializes BB...'

blackbox_initialize yes
git commit -m'INITIALIZE BLACKBOX' keyrings .gitignore


PHASE 'and adds herself as an admin.'

blackbox_addadmin alice@example.com
git commit -m'NEW ADMIN: alice@example.com' keyrings/live/pubring.gpg keyrings/live/trustdb.gpg keyrings/live/blackbox-admins.txt


PHASE 'Bob arrives.'

become_bob


PHASE 'Bob creates a gpg key.'

cat >"$gpgconfig" <<EOF
%echo Generating a basic OpenPGP key
Key-Type: default
Subkey-Type: default
Name-Real: Bob Example
Name-Comment: with weak passphrase
Name-Email: bob@example.com
Expire-Date: 0
Passphrase: b
# Do a commit here, so that we can later print "done" :-)
%commit
%echo done
EOF
gpg --no-permission-warning --batch --gen-key "$gpgconfig"

echo '========== Bob enrolls himself too.'

blackbox_addadmin bob@example.com
git commit -m'NEW ADMIN: alice@example.com' keyrings/live/pubring.gpg keyrings/live/trustdb.gpg keyrings/live/blackbox-admins.txt

PHASE 'Alice does the second part to enroll bob.'
become_alice

PHASE 'She enrolls bob.'
gpg --import keyrings/live/pubring.gpg
# TODO(tlim) That --import can be eliminated... maybe?

PHASE 'She enrolls secrets.txt.'
blackbox_register_new_file secret.txt
assert_file_missing secret.txt
assert_file_exists secret.txt.gpg

PHASE 'She decrypts secrets.txt.'
blackbox_edit_start secret.txt
assert_file_exists secret.txt
assert_file_exists secret.txt.gpg
assert_file_md5hash secret.txt "69923af35054e09cff786424e7b287aa"

PHASE 'She edits secrets.txt.'
echo 'this is MY NEW SECRET' >secret.txt
blackbox_edit_end secret.txt
assert_file_missing secret.txt
assert_file_exists secret.txt.gpg


PHASE 'Bob appears.'
become_bob

PHASE 'Bob makes sure he has all new keys.'

gpg --import keyrings/live/pubring.gpg

PHASE 'Bob postdeploys.'
blackbox_postdeploy
assert_file_exists secret.txt
assert_file_exists secret.txt.gpg
assert_file_md5hash secret.txt "08a3fa763a05c018a38e9924363b97e7"

PHASE 'Bob cleans up the secret.'
rm secret.txt

PHASE 'Bob removes alice.'
blackbox_removeadmin alice@example.com
if grep -xs >dev/null 'alice@example.com' keyrings/live/blackbox-admins.txt ; then
  echo "ASSERT FAILED: alice@example.com should be removed from keyrings/live/blackbox-admins.txt"
  echo ==== file start
  cat keyrings/live/blackbox-admins.txt
  echo ==== file end
  exit 1
fi

PHASE 'Bob reencrypts files so alice can not access them.'
blackbox_update_all_files

PHASE 'Bob decrypts secrets.txt.'
blackbox_edit_start secret.txt
assert_file_exists secret.txt
assert_file_exists secret.txt.gpg
assert_file_md5hash secret.txt "08a3fa763a05c018a38e9924363b97e7"

PHASE 'Bob edits secrets.txt.'
echo 'BOB BOB BOB BOB' >secret.txt
blackbox_edit_end secret.txt
assert_file_missing secret.txt
assert_file_exists secret.txt.gpg

PHASE 'Bob decrypts secrets.txt VERSION 3.'
blackbox_edit_start secret.txt
assert_file_exists secret.txt
assert_file_exists secret.txt.gpg
assert_file_md5hash secret.txt "beb0b0fd5701afb6f891de372abd35ed"

# TODO(tlim): Add test to make sure that now alice can NOT decrypt.

#
# ASSERTIONS
#

if [[ -e $HOME/.gnupg ]]; then
  echo "ASSERT FAILED: $HOME/.gnupg should not exist."
  exit 1
fi

find * -ls
echo cd "$test_repository"
echo rm "$test_repository"
echo DONE.
