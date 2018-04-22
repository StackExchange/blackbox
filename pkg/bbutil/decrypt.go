package bbutil

import (
	"fmt"
	"os"
)

// DecryptFile decrypts a single file.
func (bbu *RepoInfo) DecryptFile(filename, group string, overwrite bool) error {

	// change_to_vcs_root

	fmt.Fprintf(os.Stderr, "WOULD DECRYPT: %v %q %q\n", overwrite, group, filename)

	// export PATH=/usr/bin:/bin:"$PATH"

	// # Decrypt:
	// echo '========== Decrypting new/changed files: START'
	// while IFS= read <&99 -r unencrypted_file; do
	//   encrypted_file=$(get_encrypted_filename "$unencrypted_file")
	//   decrypt_file_overwrite "$encrypted_file" "$unencrypted_file"
	//   cp_permissions "$encrypted_file" "$unencrypted_file"
	//   if [[ ! -z "$FILE_GROUP" ]]; then
	//     chmod g+r "$unencrypted_file"
	//     chgrp "$FILE_GROUP" "$unencrypted_file"
	//   fi
	// done 99<"$BB_FILES"

	// echo '========== Decrypting new/changed files: DONE'

	return nil
}
