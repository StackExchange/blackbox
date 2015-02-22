#!/usr/bin/env bash

# Turn the Portfile.template into a Portfile.
# Usage:
#   mk_portfile.sh TEMPLATE OUTPUTFILE VERSION

# FIXME(tal): This code may be broken.  URL may be incorrect.

set -e
blackbox_home=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
source ${blackbox_home}/../bin/_stack_lib.sh

TEMPLATEFILE=${1?"Arg 1 must be the template."} ; shift
OUTPUTFILE=${1?"Arg 2 must be the outputfile."} ; shift
PORTVERSION=${1?"Arg 3 must be a version number like 1.20150222 (with no v)"} ; shift

make_self_deleting_tempfile bbtar
echo URL="https://github.com/StackExchange/blackbox/archive/v${PORTVERSION}.tar.gz"
curl -L "https://github.com/StackExchange/blackbox/archive/v${PORTVERSION}.tar.gz" -o ${bbtar}
RMD160=$(openssl dgst -rmd160 ${bbtar}  | awk '{ print $NF }')
SHA256=$(openssl dgst -sha256 ${bbtar}  | awk '{ print $NF }')
echo PORTVERSION=$PORTVERSION
echo RMD160=$RMD160
echo SHA256=$SHA256

sed <"$TEMPLATEFILE" >"$OUTPUTFILE" -e 's/@@VERSION@@/'"$PORTVERSION"'/g' -e 's/@@RMD160@@/'"$RMD160"'/g' -e 's/@@SHA256@@/'"$SHA256"'/g'
