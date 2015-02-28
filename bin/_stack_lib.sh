# Library functions for bash scripts at Stack Exchange.

# Usage:
#   set -e
#   . _stack_lib.sh

# ----- Utility Functions -----

function debugmsg() {
  # Log to stderr.
  echo 1>&2 LOG: """$@"""
  :
}

function logit() {
  # Log to stderr.
  echo 1>&2 LOG: """$@"""
}

function fail_out() {
    echo "FAILED:" "$*"
    echo 'Exiting...'
    exit 1
}

# on_exit and add_on_exit from http://www.linuxjournal.com/content/use-bash-trap-statement-cleanup-temporary-files
# Usage:
#   add_on_exit rm -f /tmp/foo
#   add_on_exit echo "I am exiting"
#   tempfile=$(mktemp)
#   add_on_exit rm -f "$tempfile"
function on_exit()
{
    for i in "${on_exit_items[@]}"
    do
        eval $i
    done
}

function add_on_exit()
{
    local n=${#on_exit_items[*]}
    on_exit_items[$n]="$*"
    if [[ $n -eq 0 ]]; then
        trap on_exit EXIT
    fi
}

# Securely and portably create a temporary file that will be deleted
# on EXIT.  $1 is the variable name to store the result.
function make_self_deleting_tempfile() {
  local __resultvar="$1"
  local name

  case $(uname -s) in
    Darwin )
      : ${TMPDIR:=/tmp} ;
      name=$(mktemp -t _stacklib_ )
      ;;
    Linux )
      name=$(mktemp)
      ;;
    CYGWIN* )
      name=$(mktemp)
      ;;
    * )
      echo 'ERROR: Unknown OS. Exiting.'
      exit 1
      ;;
  esac

  add_on_exit rm -f "$name"
  eval $__resultvar="$name"
}

function make_tempdir() {
  local __resultvar="$1"
  local name

  case $(uname -s) in
    Darwin )
      : "${TMPDIR:=/tmp}" ;
      name=$(mktemp -d -t _stacklib_ )
      ;;
    Linux )
      name=$(mktemp -d)
      ;;
    CYGWIN* )
      name=$(mktemp -d)
      ;;
    * )
      echo 'ERROR: Unknown OS. Exiting.'
      exit 1
      ;;
  esac

  eval $__resultvar="$name"
}

function make_self_deleting_tempdir() {
  local __resultvar="$1"
  local dirname

  make_tempdir dirname

  add_on_exit rm -rf "$dirname"
  eval $__resultvar="$dirname"
}

function fail_if_not_running_as_root() {
  if [[ $EUID -ne 0 ]]; then
    echo 'ERROR: This command should only be run as root.'
    echo 'Exiting...'
    exit 1
  fi
}

function fail_if_in_root_directory() {
  # Verify nobody has tricked us into being in "/".
  case $(uname -s) in
    Darwin )
      if [[ $(stat -f'%i' / ) == $(stat -f'%i' . ) ]] ; then
        echo 'SECURITY ALERT: The current directory is the root directory.'
        echo 'Exiting...'
        exit 1
      fi
      ;;
    Linux )
      if [[ $(stat -c'%i' / ) == $(stat -c'%i' . ) ]] ; then
        echo 'SECURITY ALERT: The current directory is the root directory.'
        echo 'Exiting...'
        exit 1
      fi
      ;;
    CYGWIN* )
      if [[ $(stat -c'%i' / ) == $(stat -c'%i' . ) ]] ; then
        echo 'SECURITY ALERT: The current directory is the root directory.'
        echo 'Exiting...'
        exit 1
      fi
      ;;
    * )
      echo 'ERROR: Unknown OS. Exiting.'
      exit 1
      ;;
  esac
}
