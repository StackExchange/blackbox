#!/bin/bash

#
# blackbox_postdeploy.sh -- Decrypt all blackbox files.
#

: ${BASEDIR:=/etc/puppet} ;
: ${CHGRP:=chgrp} ;

cd "$BASEDIR"
export PATH=/usr/bin:/bin:"$BASEDIR"/bin:"$PATH"

source blackbox_common.sh
set -e

prepare_keychain

# Decrypt:
echo '========== Decrypting new/changed files: START'
while read unencrypted_file; do
  encrypted_file=$(get_encrypted_filename "$unencrypted_file")
  decrypt_file_overwrite "$encrypted_file" "$unencrypted_file"
  chmod g+r,o-rwx "$unencrypted_file"
  $CHGRP puppet "$unencrypted_file"
done <"$BB_FILES"
echo '========== Decrypting new/changed files: DONE'
