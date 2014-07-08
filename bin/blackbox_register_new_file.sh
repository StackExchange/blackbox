#!/bin/bash

#
# blackbox_register_new_file.sh -- Enroll a new file in the blackbox system.
#
# Takes a previously unencrypted file and enters it into the blackbox
# system.  It will be kept in HG as an encrypted file.  On deployment
# to the puppet masters, it will be decrypted.  The puppet masters
# refer to the unencrypted filename.

source bin/blackbox_common.sh
set -e

fail_if_bad_environment
unencrypted_file=$(get_unencrypted_filename "$1")
encrypted_file=$(get_encrypted_filename "$1")

if [[ $1 == $encrypted_file ]]; then
  echo ERROR: Please only register unencrypted files.
  exit 1
fi

echo ========== PLAINFILE "$unencrypted_file"
echo ========== ENCRYPTED "$encrypted_file"

fail_if_not_exists "$unencrypted_file" "Please specify an existing file."
fail_if_exists "$encrypted_file" "Will not overwrite."

prepare_keychain
encrypt_file "$unencrypted_file" "$encrypted_file"
add_filename_to_cryptlist "$unencrypted_file"

# TODO(tlim): The code below should be rewritten to check
# for HG vs. GIT use and DTRT depending.

# Is the unencrypted file already in HG? (ie. are we correcting a bad situation)
SECRETSEXPOSED=$(is_in_hg ${unencrypted_file})
echo "========== CREATED: ${encrypted_file}"
echo "========== UPDATING HG:"
shred_file "$unencrypted_file"
if $SECRETSEXPOSED ; then
  hg rm -A "$unencrypted_file"
  hg add "$encrypted_file"
  COMMIT_FILES="$BB_FILES $encrypted_file $unencrypted_file"
else
  COMMIT_FILES="$BB_FILES $encrypted_file"
fi
echo 'NOTE: "already tracked!" messages are safe to ignore.'
hg add $BB_FILES $encrypted_file
hg commit -m"registered in blackbox: ${unencrypted_file}" $COMMIT_FILES
echo "========== UPDATING HG: DONE"
echo "Local repo updated.  Please push when ready."
echo "    hg push"
