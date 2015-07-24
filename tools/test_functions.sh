#!/usr/bin/env bash

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
    Darwin )
      md5 -r "$1" | awk '{ print $1 }'
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
    echo "PWD=$(/bin/pwd -P)"
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
    Darwin|FreeBSD )
      found=$(stat -f '%Sg' "$file")
      ;;
    Linux )
      found=$(stat -c '%G' "$file")
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

  if [[ "$wanted" != "$found" ]]; then
    echo "ASSERT FAILED: $file chgrp wanted=$wanted found=$found"
    exit 1
  fi
}
function assert_file_perm() {
  local wanted="$1"
  local file="$2"
  local found
  assert_file_exists "$file"

  case $(uname -s) in
    Darwin|FreeBSD )
      found=$(stat -f '%Sp' "$file")
      ;;
    # NB(tlim): CYGWIN hasn't been tested. It might be more like Darwin.
    Linux | CYGWIN* )
      found=$(stat -c '%A' "$file")
      ;;
    * )
      echo 'ERROR: Unknown OS. Exiting.'
      exit 1
      ;;
  esac

  if [[ "$wanted" != "$found" ]]; then
    echo "ASSERT FAILED: $file chgrp wanted=$wanted found=$found"
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
