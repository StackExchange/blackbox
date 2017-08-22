# Library functions for bash scripts at Stack Exchange.

# NOTE: This file is open sourced. Do not put Stack-proprietary code here.

# Usage:
#
#   set -e
#   . _stack_lib.sh

# ----- Utility Functions -----

function debugmsg() {
  # Log to stderr.
  echo 1>&2 LOG: "$@"
  :
}

function logit() {
  # Log to stderr.
  echo 1>&2 LOG: "$@"
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

function create_self_deleting_tempfile() {
  local filename

  case $(uname -s) in
    Darwin | FreeBSD )
      : "${TMPDIR:=/tmp}" ;
      filename=$(mktemp -t _stacklib_.XXXXXXXX )
      ;;
    Linux | CYGWIN* | MINGW* )
      filename=$(mktemp)
      ;;
    * )
      echo 'ERROR: Unknown OS. Exiting. (create_self_deleting_tempfile)'
      exit 1
      ;;
  esac

  add_on_exit rm -f "$filename"
  echo "$filename"
}

function create_self_deleting_tempdir() {
  local filename

  case $(uname -s) in
    Darwin | FreeBSD )
      : "${TMPDIR:=/tmp}" ;
      filename=$(mktemp -d -t _stacklib_.XXXXXXXX )
      ;;
    Linux | CYGWIN* | MINGW* )
      filename=$(mktemp -d)
      ;;
    * )
      echo 'ERROR: Unknown OS. Exiting. (create_self_deleting_tempdir)'
      exit 1
      ;;
  esac

  add_on_exit rm -rf "$filename"
  echo "$filename"
}

# Securely and portably create a temporary file that will be deleted
# on EXIT.  $1 is the variable name to store the result.
function make_self_deleting_tempfile() {
  local __resultvar="$1"
  local name

  case $(uname -s) in
    Darwin | FreeBSD )
      : "${TMPDIR:=/tmp}" ;
      name=$(mktemp -t _stacklib_.XXXXXXXX )
      ;;
    Linux | CYGWIN* | MINGW* )
      name=$(mktemp)
      ;;
    * )
      echo 'ERROR: Unknown OS. Exiting. (make_self_deleting_tempfile)'
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
    Darwin | FreeBSD )
      : "${TMPDIR:=/tmp}" ;
      # The full path to the temp directory must be short.
      # This is used by blackbox's testing suite to make a fake GNUPGHOME,
      # which needs to fit within sockaddr_un.sun_path (see unix(7)).
      name=$(mktemp -d -t SO )
      ;;
    Linux | CYGWIN* | MINGW* )
      name=$(mktemp -d)
      ;;
    * )
      echo 'ERROR: Unknown OS. Exiting. (make_tempdir)'
      exit 1
      ;;
  esac

  eval $__resultvar="$name"
}

function make_self_deleting_tempdir() {
  local __resultvar="$1"
  local dname

  make_tempdir dname

  add_on_exit rm -rf "$dname"
  eval $__resultvar="$dname"
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
    Darwin | FreeBSD )
      if [[ $(stat -f'%i' / ) == $(stat -f'%i' . ) ]] ; then
        echo 'SECURITY ALERT: The current directory is the root directory.'
        echo 'Exiting...'
        exit 1
      fi
      ;;
    Linux | CYGWIN* | MINGW* )
      if [[ $(stat -c'%i' / ) == $(stat -c'%i' . ) ]] ; then
        echo 'SECURITY ALERT: The current directory is the root directory.'
        echo 'Exiting...'
        exit 1
      fi
      ;;
    * )
      echo 'ERROR: Unknown OS. Exiting. (fail_if_in_root_directory)'
      exit 1
      ;;
  esac
}

function semverParseInto() {
    local RE='[^0-9]*\([0-9]*\)[.]\([0-9]*\)[.]\([0-9]*\)\([0-9A-Za-z-]*\)'
    #MAJOR
    eval $2=`echo $1 | sed -e "s#$RE#\1#"`
    #MINOR
    eval $3=`echo $1 | sed -e "s#$RE#\2#"`
    #MINOR
    eval $4=`echo $1 | sed -e "s#$RE#\3#"`
    #SPECIAL
    eval $5=`echo $1 | sed -e "s#$RE#\4#"`
}

function semverEQ() {
    local MAJOR_A=0
    local MINOR_A=0
    local PATCH_A=0
    local SPECIAL_A=0

    local MAJOR_B=0
    local MINOR_B=0
    local PATCH_B=0
    local SPECIAL_B=0

    semverParseInto $1 MAJOR_A MINOR_A PATCH_A SPECIAL_A
    semverParseInto $2 MAJOR_B MINOR_B PATCH_B SPECIAL_B

    if [ $MAJOR_A -ne $MAJOR_B ]; then
        return 1
    fi

    if [ $MINOR_A -ne $MINOR_B ]; then
        return 1
    fi

    if [ $PATCH_A -ne $PATCH_B ]; then
        return 1
    fi

    if [[ "_$SPECIAL_A" != "_$SPECIAL_B" ]]; then
        return 1
    fi


    return 0

}

function semverLT() {
    local MAJOR_A=0
    local MINOR_A=0
    local PATCH_A=0
    local SPECIAL_A=0

    local MAJOR_B=0
    local MINOR_B=0
    local PATCH_B=0
    local SPECIAL_B=0

    semverParseInto $1 MAJOR_A MINOR_A PATCH_A SPECIAL_A
    semverParseInto $2 MAJOR_B MINOR_B PATCH_B SPECIAL_B

    if [ $MAJOR_A -lt $MAJOR_B ]; then
        return 0
    fi

    if [[ $MAJOR_A -le $MAJOR_B  && $MINOR_A -lt $MINOR_B ]]; then
        return 0
    fi

    if [[ $MAJOR_A -le $MAJOR_B  && $MINOR_A -le $MINOR_B && $PATCH_A -lt $PATCH_B ]]; then
        return 0
    fi

    if [[ "_$SPECIAL_A"  == "_" ]] && [[ "_$SPECIAL_B"  == "_" ]] ; then
        return 1
    fi
    if [[ "_$SPECIAL_A"  == "_" ]] && [[ "_$SPECIAL_B"  != "_" ]] ; then
        return 1
    fi
    if [[ "_$SPECIAL_A"  != "_" ]] && [[ "_$SPECIAL_B"  == "_" ]] ; then
        return 0
    fi

    if [[ "_$SPECIAL_A" < "_$SPECIAL_B" ]]; then
        return 0
    fi

    return 1

}

function semverGT() {
    semverEQ $1 $2
    local EQ=$?

    semverLT $1 $2
    local LT=$?

    if [ $EQ -ne 0 ] && [ $LT -ne 0 ]; then
        return 0
    else
        return 1
    fi
}

