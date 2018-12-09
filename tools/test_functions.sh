#!/usr/bin/env bash

# NB: This is copied from _blackbox_common.sh
function get_pubring_path() {
  : "${KEYRINGDIR:=keyrings/live}" ;
  if [[ -f "${KEYRINGDIR}/pubring.gpg" ]]; then
    echo "${KEYRINGDIR}/pubring.gpg"
  else
    echo "${KEYRINGDIR}/pubring.kbx"
  fi
}

function PHASE() {
  echo '********************'
  echo '********************'
  echo '*********' """$@"""
  echo '********************'
  echo '********************'
}

function md5sum_file() {
  # Portably generate the MD5 hash of file $1.
  case $(uname -s) in
    Darwin | FreeBSD )
      md5 -r "$1" | awk '{ print $1 }'
      ;;
    NetBSD )
      md5 -q "$1"
      ;;
    SunOS )
      digest -a md5 "$1"
      ;;
    Linux )
      md5sum "$1" | awk '{ print $1 }'
      ;;
    CYGWIN* )
      md5sum "$1" | awk '{ print $1 }'
      ;;
    * )
      echo 'ERROR: Unknown OS. Exiting.'
      exit 1
      ;;
  esac
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
    echo "PWD=$(/usr/bin/env pwd -P)"
    #echo "LS START"
    #ls -la
    #echo "LS END"
    exit 1
  fi
}
function assert_file_md5hash() {
  local file="$1"
  local wanted="$2"
  assert_file_exists "$file"
  local found
  found=$(md5sum_file "$file")
  if [[ "$wanted" != "$found" ]]; then
    echo "ASSERT FAILED: $file hash wanted=$wanted found=$found"
    exit 1
  fi
}
function assert_file_group() {
  local file="$1"
  local wanted="$2"
  local found
  assert_file_exists "$file"

  case $(uname -s) in
    Darwin | FreeBSD | NetBSD )
      found=$(stat -f '%Dg' "$file")
      ;;
    Linux | SunOS )
      found=$(stat -c '%g' "$file")
      ;;
    CYGWIN* )
      echo "ASSERT_FILE_GROUP: Running on Cygwin. Not being tested."
      return 0
      ;;
    * )
      echo 'ERROR: Unknown OS. Exiting.'
      exit 1
      ;;
  esac

  echo "DEBUG: assert_file_group X${wanted}X vs. X${found}X"
  echo "DEBUG:" $(which stat)
  if [[ "$wanted" != "$found" ]]; then
    echo "ASSERT FAILED: $file chgrp group wanted=$wanted found=$found"
    exit 1
  fi
}
function assert_file_perm() {
  local wanted="$1"
  local file="$2"
  local found
  assert_file_exists "$file"

  case $(uname -s) in
    Darwin | FreeBSD | NetBSD )
      found=$(stat -f '%Sp' "$file")
      ;;
    # NB(tlim): CYGWIN hasn't been tested. It might be more like Darwin.
    Linux | CYGWIN* | SunOS )
      found=$(stat -c '%A' "$file")
      ;;
    * )
      echo 'ERROR: Unknown OS. Exiting.'
      exit 1
      ;;
  esac

  echo "DEBUG: assert_file_perm X${wanted}X vs. X${found}X"
  echo "DEBUG:" $(which stat)
  if [[ "$wanted" != "$found" ]]; then
    echo "ASSERT FAILED: $file chgrp perm wanted=$wanted found=$found"
    exit 1
  fi
}
function assert_line_not_exists() {
  local target="$1"
  local file="$2"
  assert_file_exists "$file"
  if grep -F -x -s -q >/dev/null "$target" "$file" ; then
    echo "ASSERT FAILED: line '$target' should not exist in file $file"
    echo "==== file contents: START $file"
    cat "$file"
    echo "==== file contents: END $file"
    exit 1
  fi
}
function assert_line_exists() {
  local target="$1"
  local file="$2"
  assert_file_exists "$file"
  if ! grep -F -x -s -q >/dev/null "$target" "$file" ; then
    echo "ASSERT FAILED: line '$target' should exist in file $file"
    echo "==== file contents: START $file"
    cat "$file"
    echo "==== file contents: END $file"
    exit 1
  fi
}
