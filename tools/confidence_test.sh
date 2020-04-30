#!/usr/bin/env bash

blackbox_home=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )/../bin
export PATH="${blackbox_home}:/usr/local/bin:/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/sbin:/opt/local/bin:/usr/pkg/bin:/usr/pkg/gnu/bin:${blackbox_home}"

export LANG=C.UTF-8  # Required ro "gpg --export" to work properly.

# TODO(tlim): The tests are hardcoded to this directory. This should be fixed.
export BLACKBOXDATA=keyrings/live

# This script requires many utilities, some are not
# required by the usual blackbox scripts.  Test to make
# sure we have them all.
e=false
for i in blackbox_addadmin blackbox_list_admins blackbox_register_new_file \
  cat git gpg gpg-agent mkdir pinentry pinentry-tty rm tar which ; do
  if ! which >/dev/null 2>&1  $i ; then
    echo ERROR: Command not in PATH: $i
    e=true
  fi
done
if $e ; then
  echo 'Exiting. Please install the above commands.'
  echo 'This script requires many utilities not required by blackbox itself.'
  echo PATH="$PATH"
  exit 1
fi

set -e
. _stack_lib.sh
. tools/test_functions.sh

PHASE 'UNIT TESTS'

_blackbox_common_test.sh

PHASE 'SYSTEM TESTS'

make_tempdir test_repository
cd "$test_repository"

make_self_deleting_tempdir fake_alice_home
make_self_deleting_tempdir fake_bob_home
export GNUPGHOME="$fake_alice_home"
#echo 'pinentry-program' "$(which pinentry-insecure-fake)" >> "$GNUPGHOME/gpg-agent.conf"
echo 'pinentry-program' "$(which pinentry-tty)" >> "$GNUPGHOME/gpg-agent.conf"
#pinentry-program "${blackbox_home}/pinentry-fake-insecure"
eval "$(gpg-agent --homedir "$fake_alice_home" --daemon)"
GPG_AGENT_INFO_ALICE="$GPG_AGENT_INFO"

export GNUPGHOME="$fake_bob_home"
#echo 'pinentry-program' "$(which pinentry-insecure-fake)" >> "$GNUPGHOME/gpg-agent.conf"
echo 'pinentry-program' "$(which pinentry-tty)" >> "$GNUPGHOME/gpg-agent.conf"
eval "$(gpg-agent --homedir "$fake_bob_home" --daemon)"
GPG_AGENT_INFO_BOB="$GPG_AGENT_INFO"

function become_alice() {
  export GNUPGHOME="$fake_alice_home"
  export GPG_AGENT_INFO="$GPG_AGENT_INFO_ALICE"
  echo BECOMING ALICE: GNUPGHOME="$GNUPGHOME AGENT=$GPG_AGENT_INFO"
  mkdir -p .git ; touch .git/config
  git init
  git config user.name "Alice Example"
  git config user.email alice@example.com
}

function become_bob() {
  export GNUPGHOME="$fake_bob_home"
  export GPG_AGENT_INFO="$GPG_AGENT_INFO_BOB"
  mkdir -p .git ; touch .git/config
  git config user.name "Bob Example"
  git config user.email bob@example.com
}

PHASE 'Alice creates a repo.  She creates secret.txt.'

become_alice
echo 'this is my secret' >secret.txt


PHASE 'Alice wants to be part of the secret system.'
PHASE 'She creates a GPG key...'

make_self_deleting_tempfile gpgconfig
cat >"$gpgconfig" <<EOF
%echo Generating a basic OpenPGP key
Key-Type: RSA
Subkey-Type: RSA
Name-Real: Alice Example
Name-Comment: my password is the lowercase letter a
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
git commit -m'NEW ADMIN: alice@example.com' keyrings/live/pubring.??? keyrings/live/trustdb.gpg keyrings/live/blackbox-admins.txt

make_self_deleting_tempfile adminlist1
blackbox_list_admins >"$adminlist1"
assert_file_md5hash "$adminlist1" "aa1db827772e1d51d453b844394b7617"


PHASE 'Bob arrives.'

become_bob


PHASE 'Bob creates a gpg key.'

cat >"$gpgconfig" <<EOF
%echo Generating a basic OpenPGP key
Key-Type: RSA
Subkey-Type: RSA
Name-Real: Bob Example
Name-Comment: my password is the lowercase letter b
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
git commit -m'NEW ADMIN: alice@example.com' keyrings/live/pubring.??? keyrings/live/trustdb.gpg keyrings/live/blackbox-admins.txt

make_self_deleting_tempfile adminlist2
blackbox_list_admins >"$adminlist2"
assert_file_md5hash "$adminlist2" "0b82b0b3c96e6e5dd5faf60493fe5cf7"

PHASE 'Alice does the second part to enroll bob.'
become_alice

PHASE 'She enrolls bob.'
#gpg --import $(get_pubring_path)
## TODO(tlim) That --import can be eliminated... maybe?

PHASE 'She enrolls secrets.txt.'
blackbox_register_new_file secret.txt
assert_file_missing secret.txt
assert_file_exists secret.txt.gpg
assert_line_exists '/secret.txt' .gitignore

PHASE 'She cats secrets.txt.gpg'
make_self_deleting_tempfile catsecret
blackbox_cat secret.txt.gpg > $catsecret
assert_line_exists 'this is my secret' $catsecret

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


PHASE 'Alice copies files to a non-repo directory. (NO REPO)'

# Copy the repo entirely:
make_self_deleting_tempdir fake_alice_filedir
tar cf - . | ( cd "$fake_alice_filedir" && tar xpvf - )
# Remove the .git directory
rm -rf "$fake_alice_filedir/.git"
(
cd "$fake_alice_filedir"
assert_file_missing '.git'
assert_file_exists 'secret.txt.gpg'
assert_file_missing 'secret.txt'
blackbox_postdeploy
assert_file_missing '.git'
assert_file_exists 'secret.txt.gpg'
assert_file_exists 'secret.txt'
assert_file_md5hash secret.txt "08a3fa763a05c018a38e9924363b97e7"

PHASE 'Alice shreds these non-repo files. (NO REPO)'
blackbox_shred_all_files
assert_file_missing '.git'
assert_file_exists 'secret.txt.gpg'
assert_file_missing 'secret.txt'

PHASE 'Alice decrypts secrets.txt (NO REPO).'
blackbox_edit_start secret.txt
assert_file_exists secret.txt
assert_file_exists secret.txt.gpg
assert_file_md5hash secret.txt "08a3fa763a05c018a38e9924363b97e7"

PHASE 'Alice edits secrets.txt. (NO REPO EDIT)'
echo 'NOREPO EDIT' >secret.txt
assert_file_md5hash secret.txt "d3e6bbdfc76fae7fd0a921f3408db1d1"
blackbox_edit_end secret.txt
assert_file_missing secret.txt
assert_file_exists secret.txt.gpg

PHASE 'Alice decrypts secrets.txt (NO REPO EDIT).'
blackbox_edit_start secret.txt
assert_file_exists secret.txt
assert_file_exists secret.txt.gpg
assert_file_md5hash secret.txt "d3e6bbdfc76fae7fd0a921f3408db1d1"
)

PHASE 'appears.'
become_bob

#PHASE 'Bob makes sure he has all new keys.'
#gpg --import keyrings/live/pubring.???

# Pick a GID to use:
# This users's default group:
DEFAULT_GID_NUM=$(id -g)
# Pick a group that is not the default group:
TEST_GID_NUM=$(grep -v "$DEFAULT_GID_NUM" /etc/group | cut -d: -f3 | sort -rn | head -1)
echo "DEFAULT_GID_NUM=$DEFAULT_GID_NUM"
echo "TEST_GID_NUM=$TEST_GID_NUM"

PHASE 'Bob postdeploys... default.'
blackbox_postdeploy
assert_file_exists secret.txt
assert_file_exists secret.txt.gpg
assert_file_md5hash secret.txt "08a3fa763a05c018a38e9924363b97e7"
assert_file_group secret.txt "$DEFAULT_GID_NUM"

PHASE 'Bob postdeploys... with a GID.'
blackbox_postdeploy "$TEST_GID_NUM"
assert_file_exists secret.txt
assert_file_exists secret.txt.gpg
assert_file_md5hash secret.txt "08a3fa763a05c018a38e9924363b97e7"
assert_file_group secret.txt "$TEST_GID_NUM"

PHASE 'Bob cleans up the secret.'
rm secret.txt

PHASE 'Bob removes Alice.'
blackbox_removeadmin alice@example.com
assert_line_not_exists 'alice@example.com' keyrings/live/blackbox-admins.txt

make_self_deleting_tempfile adminlist3
blackbox_list_admins >"$adminlist3"
assert_file_md5hash "$adminlist3" "aadbfafd76ea66ff40dbfd239a69067f"

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

PHASE 'Bob exposes a secret in the repo.'
echo 'this is my exposed secret' >mistake.txt
git add mistake.txt
git commit -m'Oops I am committing a secret to the repo.' mistake.txt

PHASE 'Bob corrects it by registering it.'
blackbox_register_new_file mistake.txt
assert_file_missing mistake.txt
assert_file_exists mistake.txt.gpg
# NOTE: It is still in the history. That should be corrected someday.
assert_line_exists '/mistake.txt' .gitignore

PHASE 'Bob enrolls my/path/to/relsecrets.txt.'
mkdir my my/path my/path/to
echo 'New secret' > my/path/to/relsecrets.txt
cd my/path/to
blackbox_register_new_file relsecrets.txt
assert_file_missing relsecrets.txt
assert_file_exists relsecrets.txt.gpg
assert_file_missing .gitignore
assert_file_exists ../../../.gitignore
assert_line_exists '/my/path/to/relsecrets.txt' ../../../.gitignore

PHASE 'Bob decrypts relsecrets.txt.'
cd ..
blackbox_edit_start to/relsecrets.txt
assert_file_exists to/relsecrets.txt
assert_file_exists to/relsecrets.txt.gpg
assert_file_md5hash to/relsecrets.txt "c47f9c3c8ce03d895b883ac22384cb67"
cd ../..

PHASE 'Bob enrolls !important!.txt'
echo A very important file >'!important!.txt'
blackbox_register_new_file '!important!.txt'
assert_file_missing '!important!.txt'
assert_file_exists '!important!.txt'.gpg
assert_line_exists '/!important!.txt' .gitignore

PHASE 'Bob enrolls #andpounds.txt'
echo A very commented file >'#andpounds.txt'
blackbox_register_new_file '#andpounds.txt'
assert_file_missing '#andpounds.txt'
assert_file_exists '#andpounds.txt'.gpg
assert_line_exists '/#andpounds.txt' .gitignore

PHASE 'Bob enrolls stars*bars?.txt'
echo A very wild and questioned file >'stars*bars?.txt'
blackbox_register_new_file 'stars*bars?.txt'
assert_file_missing 'stars*bars?.txt'
assert_file_exists 'stars*bars?.txt'.gpg
assert_line_exists '/stars\*bars\?.txt' .gitignore

PHASE 'Bob enrolls space space.txt'
echo A very spacey file >'space space.txt'
blackbox_register_new_file 'space space.txt'
assert_file_missing 'space space.txt'
assert_file_exists 'space space.txt'.gpg
assert_line_exists '/space space.txt' .gitignore

PHASE 'Bob checks out stars*bars?.txt.'
blackbox_edit_start 'stars*bars?.txt'
assert_file_exists 'stars*bars?.txt'
assert_file_exists 'stars*bars?.txt'
assert_file_md5hash 'stars*bars?.txt' "448e018faade28cede2bf6f33c3c2dfb"

PHASE 'Bob checks out space space.txt.'
blackbox_edit_start 'space space.txt'
assert_file_exists 'space space.txt'
assert_file_exists 'space space.txt'
assert_file_md5hash 'space space.txt' "de1d4e4a07046f81af5d3c0194b78742"

PHASE 'Bob shreds all exposed files.'
assert_file_exists 'my/path/to/relsecrets.txt'
assert_file_exists 'secret.txt'
blackbox_shred_all_files
which blackbox_shred_all_files
assert_file_missing '!important!.txt'
assert_file_missing '#andpounds.txt'
assert_file_missing 'mistake.txt'
assert_file_missing 'my/path/to/relsecrets.txt'
assert_file_missing 'secret.txt'
assert_file_missing 'space space.txt'
assert_file_missing 'stars*bars?.txt'
assert_file_exists '!important!.txt.gpg'
assert_file_exists '#andpounds.txt.gpg'
assert_file_exists 'mistake.txt.gpg'
assert_file_exists 'my/path/to/relsecrets.txt.gpg'
assert_file_exists 'secret.txt.gpg'
assert_file_exists 'space space.txt.gpg'
assert_file_exists 'stars*bars?.txt.gpg'

PHASE 'Bob updates all files.'
blackbox_update_all_files
assert_file_missing '!important!.txt'
assert_file_missing '#andpounds.txt'
assert_file_missing 'mistake.txt'
assert_file_missing 'my/path/to/relsecrets.txt'
assert_file_missing 'secret.txt'
assert_file_missing 'space space.txt'
assert_file_missing 'stars*bars?.txt'
assert_file_exists '!important!.txt.gpg'
assert_file_exists '#andpounds.txt.gpg'
assert_file_exists 'mistake.txt.gpg'
assert_file_exists 'my/path/to/relsecrets.txt.gpg'
assert_file_exists 'secret.txt.gpg'
assert_file_exists 'space space.txt.gpg'
assert_file_exists 'stars*bars?.txt.gpg'

PHASE 'Bob DEregisters mistake.txt'
touch 'mistake.txt'
blackbox_deregister_file 'mistake.txt.gpg'
assert_file_exists 'keyrings/live/blackbox-admins.txt'
assert_file_exists 'keyrings/live/blackbox-files.txt'
assert_line_not_exists 'mistake.txt' 'keyrings/live/blackbox-files.txt'
assert_file_missing 'mistake.txt.gpg'
assert_file_exists 'mistake.txt'
# Now remove 'mistake.txt' to leave the area clean.
rm 'mistake.txt'

PHASE 'Bob enrolls multiple files: multi1.txt and multi2.txt'
echo 'One singular sensation.' >'multi1.txt'
echo 'Another singular sensation.' >'multi2.txt'
blackbox_register_new_file 'multi1.txt' 'multi2.txt'
assert_file_missing 'multi1.txt'
assert_file_exists 'multi1.txt'.gpg
assert_line_exists '/multi1.txt' .gitignore
assert_file_missing 'multi2.txt'
assert_file_exists 'multi2.txt'.gpg
assert_line_exists '/multi2.txt' .gitignore

PHASE 'Alice returns. She should be locked out'
assert_file_missing 'secret.txt'
become_alice
PHASE 'Alice tries to decrypt secret.txt. Is blocked.'
if blackbox_edit_start secret.txt ; then
  echo 'ERROR: Alice was able to decrypt secret.txt!  She should have been blocked.'
  exit 1
else
  echo 'NOTE: Alice was not able to decrypt secret.txt as expected.'
fi

PHASE 'Bob returns.  Tries to update all files with a corrupt blackbox-admins.txt'
become_bob
# Corrupt the blackbox-admins.txt list:
echo 'abba@notarealuser.com' >> keyrings/live/blackbox-admins.txt
# Make sure it fails.
if blackbox_update_all_files; then
  echo '!!!!! blackbox_update_all_files should have failed and it did NOT.'
  exit 1
fi
# Cleanup:
blackbox_removeadmin abba@notarealuser.com

make_self_deleting_tempfile adminlist4
blackbox_list_admins >"$adminlist4"
assert_file_md5hash "$adminlist4" "aadbfafd76ea66ff40dbfd239a69067f"


# TODO: Create a new directory. "git clone" the repo into it.

#
# ASSERTIONS
#

echo '========== Verifying .gnupg was not accidentally created.'

if [[ -e $HOME/.gnupg ]]; then
  echo "ASSERT FAILED: $HOME/.gnupg should not exist."
  exit 1
fi

echo '========== DONE with tests.  Outputing some diagnostics:'

find .git?* * -type f -ls
echo cd "$test_repository"
echo rm -rf "$test_repository"
echo 'SUCCESS! Doing final clean-up then exiting.'
