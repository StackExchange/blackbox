#!/usr/bin/env bash

#
# Common constants and functions used by the blackbox_* utilities.
#

# Usage:
#
#   set -e
#   source "${0%/*}/_blackbox_common.sh"

# Load additional useful functions
source "${0%/*}"/_stack_lib.sh

# Where are we?
: "${BLACKBOX_HOME:="$(cd "${0%/*}" ; pwd)"}" ;

# Where in the VCS repo should the blackbox data be found?
: "${BLACKBOXDATA:=keyrings/live}" ;   # If BLACKBOXDATA not set, set it.


# If $EDITOR is not set, set it to "vi":
: "${EDITOR:=vi}" ;

# Allow overriding gpg command
: "${GPG:=gpg}" ;

function physical_directory_of() {
  local d=$(dirname "$1")
  local f=$(basename "$1")
  (cd "$d" && echo "$(pwd -P)/$f" )
}

# Set REPOBASE to the top of the repository
# Set VCS_TYPE to 'git', 'hg', 'svn' or 'unknown'
if which >/dev/null 2>/dev/null git && git rev-parse --show-toplevel >/dev/null 2>&1 ; then
  VCS_TYPE=git
  REPOBASE=$(git rev-parse --show-toplevel)
elif [ -d ".svn" ] ; then
  # Find topmost dir with .svn sub-dir
  parent=""
  grandparent="."
  while [ -d "$grandparent/.svn" ]; do
    parent=$grandparent
    grandparent="$parent/.."
  done

  REPOBASE=$(cd "$parent" ; pwd)
  VCS_TYPE=svn
elif which >/dev/null 2>/dev/null hg && hg root >/dev/null 2>&1 ; then
  # NOTE: hg has to be tested last because it always "succeeds".
  VCS_TYPE=hg
  REPOBASE=$(hg root 2>/dev/null)
else
  # We aren't in a repo at all.  Assume the cwd is the root
  # of the tree.
  VCS_TYPE=unknown
  REPOBASE="$(pwd)"
fi
export VCS_TYPE
export REPOBASE=$(physical_directory_of "$REPOBASE")
# FIXME: Verify this function by checking for .hg or .git
# after determining what we believe to be the answer.

KEYRINGDIR="$REPOBASE/$BLACKBOXDATA"
BB_ADMINS_FILE="blackbox-admins.txt"
BB_ADMINS="${KEYRINGDIR}/${BB_ADMINS_FILE}"
BB_FILES_FILE="blackbox-files.txt"
BB_FILES="${KEYRINGDIR}/${BB_FILES_FILE}"
SECRING="${KEYRINGDIR}/secring.gpg"
: "${DECRYPT_UMASK:=0022}" ;
# : ${DECRYPT_UMASK:=o=} ;

# Checks if $1 is 0 bytes, and if $1/keyrings
# is a directory
function is_blackbox_repo() {
  if [[ -n "$1" ]] && [[ -d "$1/keyrings" ]]; then
    return 0 # Yep, its a repo
  else
    return 1
  fi
}

# Return error if not on cryptlist.
function is_on_cryptlist() {
  # Assumes $1 does NOT have the .gpg extension
  file_contains_line "$BB_FILES" "$(vcs_relative_path "$1")"
}

# Exit with error if a file exists.
function fail_if_exists() {
  if [[ -f "$1" ]]; then
    echo ERROR: "$1" exists.  "$2" >&2
    echo Exiting... >&2
    exit 1
  fi
}

# Exit with error if a file is missing.
function fail_if_not_exists() {
  if [[ ! -f "$1" ]]; then
    echo ERROR: "$1" not found.  "$2" >&2
    echo Exiting... >&2
    exit 1
  fi
}

# Exit we we aren't in a VCS repo.
function fail_if_not_in_repo() {
  if [[ $VCS_TYPE = "unknown" ]]; then
    echo "ERROR: This must be run in a VCS repo: git, hg, or svn." >&2
    echo Exiting... >&2
    exit 1
  fi
}

# Exit with error if filename is not registered on blackbox list.
function fail_if_not_on_cryptlist() {
  # Assumes $1 does NOT have the .gpg extension

  local name="$1"

  if ! is_on_cryptlist "$name" ; then
    echo "ERROR: $name not found in $BB_FILES" >&2
    echo "PWD=$(/bin/pwd)" >&2
    echo 'Exiting...' >&2
    exit 1
  fi
}

# Exit with error if keychain contains secret keys.
function fail_if_keychain_has_secrets() {
  if [[ -s ${SECRING} ]]; then
    echo 'ERROR: The file' "$SECRING" 'should be empty.' >&2
    echo 'Did someone accidentally add this private key to the ring?' >&2
    echo 'Exiting...' >&2
    exit 1
  fi
}

function get_pubring_path() {
  if [[ -f "${KEYRINGDIR}/pubring.gpg" ]]; then
    echo "${KEYRINGDIR}/pubring.gpg"
  else
    echo "${KEYRINGDIR}/pubring.kbx"
  fi
}

# Output the unencrypted filename.
function get_unencrypted_filename() {
  echo "$(dirname "$1")/$(basename "$1" .gpg)" | sed -e 's#^\./##'
}

# Output the encrypted filename.
function get_encrypted_filename() {
  echo "$(dirname "$1")/$(basename "$1" .gpg).gpg" | sed -e 's#^\./##'
}

# Prepare keychain for use.
function prepare_keychain() {
  echo '========== Importing keychain: START' >&2
  $GPG --import "$(get_pubring_path)" 2>&1 | egrep -v 'not changed$' >&2
  echo '========== Importing keychain: DONE' >&2
}

# Add file to list of encrypted files.
function add_filename_to_cryptlist() {
  # If the name is already on the list, this is a no-op.
  # However no matter what the datestamp is updated.
  
  # https://github.com/koalaman/shellcheck/wiki/SC2155
  local name
  name=$(vcs_relative_path "$1")

  if file_contains_line "$BB_FILES" "$name" ; then
    echo "========== File is registered. No need to add to list."
  else
    echo "========== Adding file to list."
    touch "$BB_FILES"
    sort -u -o "$BB_FILES" <(echo "$name") "$BB_FILES"
  fi
}

# Removes a file from the list of encrypted files
function remove_filename_from_cryptlist() {
  # If the name is not already on the list, this is a no-op.

  # https://github.com/koalaman/shellcheck/wiki/SC2155
  local name
  name=$(vcs_relative_path "$1")

  if ! file_contains_line "$BB_FILES" "$name" ; then
    echo "========== File is not registered. No need to remove from list."
  else
    echo "========== Removing file from list."
    remove_line "$BB_FILES" "$name"
  fi
}

# Print out who the current BB ADMINS are:
function disclose_admins() {
  echo "========== blackbox administrators are:"
  cat "$BB_ADMINS"
}

# Encrypt file, overwriting .gpg if it exists.
function encrypt_file() {
  local unencrypted
  local encrypted
  unencrypted="$1"
  encrypted="$2"

  echo "========== Encrypting: $unencrypted" >&2
  $GPG --use-agent --yes --trust-model=always --encrypt -o "$encrypted"  $(awk '{ print "-r" $1 }' < "$BB_ADMINS") "$unencrypted" >&2
  echo '========== Encrypting: DONE' >&2
}

# Decrypt .gpg file, asking "yes/no" before overwriting unencrypted file.
function decrypt_file() {
  local encrypted
  local unencrypted
  local old_umask
  encrypted="$1"
  unencrypted="$2"

  echo "========== EXTRACTING $unencrypted" >&2

  old_umask=$(umask)
  umask "$DECRYPT_UMASK"
  $GPG --use-agent -q --decrypt -o "$unencrypted" "$encrypted" >&2
  umask "$old_umask"
}

# Decrypt .gpg file, overwriting unencrypted file if it exists.
function decrypt_file_overwrite() {
  local encrypted
  local unencrypted
  local old_hash
  local new_hash
  local old_umask
  encrypted="$1"
  unencrypted="$2"

  if [[ -f "$unencrypted" ]]; then
    old_hash=$(md5sum_file "$unencrypted")
  else
    old_hash=unmatchable
  fi

  old_umask=$(umask)
  umask "$DECRYPT_UMASK"
  $GPG --use-agent --yes -q --decrypt -o "$unencrypted" "$encrypted" >&2
  umask "$old_umask"

  new_hash=$(md5sum_file "$unencrypted")
  if [[ "$old_hash" != "$new_hash" ]]; then
    echo "========== EXTRACTED $unencrypted" >&2
  fi
}

# Shred a file.  If shred binary does not exist, delete it.
function shred_file() {
  local name
  local CMD
  local OPT
  name="$1"

  if which shred >/dev/null 2>/dev/null ; then
    CMD=shred
    OPT=-u
  elif which srm >/dev/null 2>/dev/null ; then
    #NOTE: srm by default uses 35-pass Gutmann algorithm
    CMD=srm
    OPT=-f
  else
    CMD=rm
    OPT=-f
  fi

  $CMD $OPT -- "$name"
}

# $1 is the name of a file that contains a list of files.
# For each filename, output the individual subdirectories
# leading up to that file. i.e. one one/two one/two/three
function enumerate_subdirs() {
  local listfile
  local dir
  local filename
  listfile="$1"

  while read filename; do
    dir=$(dirname "$filename")
    while [[ $dir != '.' && $dir != '/' ]]; do
      echo "$dir"
      dir=$(dirname "$dir")
    done
  done <"$listfile" | sort -u
}
 

# chdir to the base of the repo.
function change_to_vcs_root() {
  # if vcs_root not explicitly defined, use $REPOBASE

  local rbase=${1:-$REPOBASE} # use $1 but if unset use $REPOBASE

  cd "$rbase"

}

# $1 is a string pointing to a directory.  Outputs a
# list of  valid blackbox repos,relative to $1
function enumerate_blackbox_repos() {
  if [[ -z "$1" ]]; then
    echo "enumerate_blackbox_repos: ERROR: No Repo provided to Enumerate"
    exit 1
  fi

  # https://github.com/koalaman/shellcheck/wiki/Sc2045
  for dir in $1*/; do
    if is_blackbox_repo "$dir"; then
      echo "$dir"
    fi
  done
}

# Output the path of a file relative to the repo base
function vcs_relative_path() {
  # Usage: vcs_relative_path file
  local name="$1"
  #python -c 'import os ; print(os.path.relpath("'"$(pwd -P)"'/'"$name"'", "'"$REPOBASE"'"))'
  local p=$( printf "%s" "$( pwd -P )/${1}" | sed 's#//*#/#g' )
  local name="${p#$REPOBASE}"
  name=$( printf "%s" "$name" | sed 's#^/##g' | sed 's#/$##g' )
  printf "%s" "$name"
}

# Removes a line from a text file
function remove_line() {
  local tempfile

  make_self_deleting_tempfile tempfile

  # Ensure source file exists
  touch "$1"
  grep -Fsxv "$2" "$1" > "$tempfile" || true

  # Using cat+rm instead of cp will preserve permissions/ownership
  cat "$tempfile" > "$1"
}

# Determine if a file contains a given line
function file_contains_line() {
  # $1: the file
  # $2: the line
  grep -xsqF "$2" "$1"
}

#
# Portability Section:
#

#
# Abstract the difference between Linux and Mac OS X:
#

function md5sum_file() {
  # Portably generate the MD5 hash of file $1.
  case $(uname -s) in
    Darwin )
      md5 -r "$1" | awk '{ print $1 }'
      ;;
    Linux | CYGWIN* | MINGW* )
      md5sum "$1" | awk '{ print $1 }'
      ;;
    * )
      echo 'ERROR: Unknown OS. Exiting. (md5sum_file)'
      exit 1
      ;;
  esac
}

function cp_permissions() {
  # Copy the perms of $1 onto $2 .. end.
  case $(uname -s) in
    Darwin )
      chmod $( stat -f '%p' "$1" ) "${@:2}"
      ;;
    Linux | CYGWIN* | MINGW* )
      chmod --reference "$1" "${@:2}"
      ;;
    * )
      echo 'ERROR: Unknown OS. Exiting. (cp_permissions)'
      exit 1
      ;;
  esac
}


#
# Abstract the difference between git and hg:
#

# Is this file in the current repo?
function is_in_vcs() {
  is_in_$VCS_TYPE "$@"
}
# Mercurial
function is_in_hg() {
  local filename
  filename="$1"

  if hg locate "$filename" ; then
    echo true
  else
    echo false
  fi
}
# Git:
function is_in_git() {
  local filename
  filename="$1"

  if git ls-files --error-unmatch >/dev/null 2>&1 -- "$filename" ; then
    echo true
  else
    echo false
  fi
}
# Subversion
function is_in_svn() {
  local filename
  filename="$1"

  if svn list "$filename" ; then
    echo true
  else
    echo false
  fi
}
# Perforce
function is_in_p4() {
  local filename
  filename="$1"

  if p4 list "$filename" ; then
    echo true
  else
    echo false
  fi
}
# No repo
function is_in_unknown() {
  echo true
}


# Add a file to the repo (but don't commit it).
function vcs_add() {
  vcs_add_$VCS_TYPE "$@"
}
# Mercurial
function vcs_add_hg() {
  hg add "$@"
}
# Git
function vcs_add_git() {
  git add "$@"
}
# Subversion
function vcs_add_svn() {
  svn add --parents "$@"
}
# Perfoce
function vcs_add_p4() {
  p4 add "$@"
}
# No repo
function vcs_add_unknown() {
  :
}


# Commit a file to the repo
function vcs_commit() {
  vcs_commit_$VCS_TYPE "$@"
}
# Mercurial
function vcs_commit_hg() {
  hg commit -m "$@"
}
# Git
function vcs_commit_git() {
  git commit -m "$@"
}
# Subversion
function vcs_commit_svn() {
  svn commit -m "$@"
}
# Perforce
function vcs_commit_p4() {
  p4 submit -d "$@"
}
# No repo
function vcs_commit_unknown() {
  :
}


# Remove file from repo, even if it was deleted locally already.
# If it doesn't exist yet in the repo, it should be a no-op.
function vcs_remove() {
  vcs_remove_$VCS_TYPE "$@"
}
# Mercurial
function vcs_remove_hg() {
  hg rm -A -- "$@"
}
# Git
function vcs_remove_git() {
  git rm --ignore-unmatch -f -- "$@"
}
# Subversion
function vcs_remove_svn() {
  svn delete "$@"
}
# Perforce
function vcs_remove_p4() {
  p4 delete "$@"
}
# No repo
function vcs_remove_unknown() {
  :
}

# Get a path for the ignore file if possible in current vcs
function vcs_ignore_file_path() {
  vcs_ignore_file_path_$VCS_TYPE
}
# Mercurial
function vcs_ignore_file_path_hg() {
  echo "$REPOBASE/.hgignore"
}
# Git
function vcs_ignore_file_path_git() {
  echo "$REPOBASE/.gitignore"
}


# Ignore a file in a repo.  If it was already ignored, this is a no-op.
function vcs_ignore() {
  local file
  for file in "$@"; do
    vcs_ignore_$VCS_TYPE "$file"
  done
}
# Mercurial
function vcs_ignore_hg() {
  vcs_ignore_generic_file "$(vcs_ignore_file_path)" "$file"
}
# Git
function vcs_ignore_git() {
  vcs_ignore_generic_file "$(vcs_ignore_file_path)" "$file"
}
# Subversion
function vcs_ignore_svn() {
  svn propset svn:ignore "$(vcs_relative_path "$file")"
}
# Perforce
function vcs_ignore_p4() {
  :
}
# No repo
function vcs_ignore_unknown() {
  :
}
# Generic - add line to file
function vcs_ignore_generic_file() {
  local file
  file="$(vcs_relative_path "$2")"
  file="${file/\$\//}"
  file="$(echo "/$file" | sed 's/\([\*\?]\)/\\\1/g')"
  if ! file_contains_line "$1" "$file" ; then
    echo "$file" >> "$1"
    vcs_add "$1"
  fi
}


# Notice (un-ignore) a file in a repo.  If it was not ignored, this is
# a no-op
function vcs_notice() {
  local file
  for file in "$@"; do
    vcs_notice_$VCS_TYPE "$file"
  done
}
# Mercurial
function vcs_notice_hg() {
  vcs_notice_generic_file "$REPOBASE/.hgignore" "$file"
}
# Git
function vcs_notice_git() {
  vcs_notice_generic_file "$REPOBASE/.gitignore" "$file"
}
# Subversion
function vcs_notice_svn() {
  svn propdel svn:ignore "$(vcs_relative_path "$file")"
}
# Perforce
function vcs_notice_p4() {
  :
}
# No repo
function vcs_notice_unknown() {
  :
}
# Generic - remove line to file
function vcs_notice_generic_file() {
  local file
  file="$(vcs_relative_path "$2")"
  file="${file/\$\//}"
  file="$(echo "/$file" | sed 's/\([\*\?]\)/\\\1/g')"
  if file_contains_line "$1" "$file" ; then
    remove_line "$1" "$file"
    vcs_add "$1"
  fi
  if file_contains_line "$1" "${file:1}" ; then
    echo "WARNING:  Found a non-absolute ignore match in $1"
    echo "WARNING:  Confirm the pattern is intended to only exclude $file"
    echo "WARNING:  If so, manually update the ignore file"
  fi
}
