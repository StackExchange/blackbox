#!/bin/bash

#
# blackbox_edit_end.sh -- Re-encrypt file after edits.
#

source blackbox_common.sh
set -e

fail_if_bad_environment
unencrypted_file=$(get_unencrypted_filename "$1")
encrypted_file=$(get_encrypted_filename "$1")
echo ========== PLAINFILE "$unencrypted_file"
echo ========== ENCRYPTED "$encrypted_file"

fail_if_not_on_cryptlist "$unencrypted_file"
fail_if_not_exists "$unencrypted_file" "No unencrypted version to encrypt!"
fail_if_keychain_has_secrets

encrypt_file "$unencrypted_file" "$encrypted_file"
shred_file "$unencrypted_file"

echo "========== UPDATED ${encrypted_file}"
echo "Likely next step:"
echo "    git commit -m\"${encrypted_file} updated\" $encrypted_file"
