#!/usr/bin/env bash

# Turn the Portfile.template into a Portfile.
# Usage:
#   mk_portfile.sh TEMPLATE OUTPUTFILE VERSION

set -e
blackbox_home=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
source ${blackbox_home}/../bin/_stack_lib.sh

TEMPLATEFILE=tools/Portfile.template
OUTPUTFILE=Portfile
PORTVERSION=${1?"Arg 1 must be a version number like 1.20150222 (with no v)"} ; shift

# Add the version number to the template.
sed <"$TEMPLATEFILE" >"$OUTPUTFILE" -e 's/@@VERSION@@/'"$PORTVERSION"'/g'

# Test it. Record the failure in $checksumout
fgrep >/dev/null -x 'file:///var/tmp/ports' /opt/local/etc/macports/sources.conf || sudo sed -i -e '1s@^@file:///var/tmp/ports\'$'\n@' /opt/local/etc/macports/sources.conf
rm -rf /var/tmp/ports
mkdir -p /var/tmp/ports/security/vcs_blackbox
cp Portfile /var/tmp/ports/security/vcs_blackbox
( cd /var/tmp/ports && sudo portindex )
make_self_deleting_tempfile checksumout
set +e
sudo port -v checksum vcs_blackbox > "$checksumout" 2>/dev/null
ret=$?

# If it failed, grab the checksums. Then re-process the template with them.
if [[ $ret != 0 ]]; then
  RMD160=$(awk <"$checksumout" '/^Distfile checksum: .*rmd160/ { print $NF }')
  SHA256=$(awk <"$checksumout" '/^Distfile checksum: .*sha256/ { print $NF }')
  echo RMD160=$RMD160
  echo SHA256=$SHA256
  echo
  if [[ $RMD160 != '' && $SHA256 != '' ]]; then
    sed <"$TEMPLATEFILE" >"$OUTPUTFILE" -e 's/@@VERSION@@/'"$PORTVERSION"'/g' -e 's/@@RMD160@@/'"$RMD160"'/g' -e 's/@@SHA256@@/'"$SHA256"'/g'
    cp Portfile /var/tmp/ports/security/vcs_blackbox
    ( cd /var/tmp/ports && sudo portindex )
    sudo port -v checksum vcs_blackbox
  fi
fi

# Generate the diff
cp /opt/local/var/macports/sources/rsync.macports.org/release/tarballs/ports/security/vcs_blackbox/Portfile /var/tmp/ports/security/vcs_blackbox/Portfile.orig
( cd /var/tmp/ports/security/vcs_blackbox && diff --ignore-matching-lines='Id:' -u Portfile.orig Portfile ) > Portfile-vcs_blackbox.diff
open -R Portfile-vcs_blackbox.diff

echo
echo 'portfile is in:'
echo '                /var/tmp/ports/security/vcs_blackbox/Portfile'
echo 'cleanup:'
echo '                sudo vi /opt/local/etc/macports/sources.conf'

echo "
PLEASE OPEN A TICKET WITH THIS INFORMATION:
    https://trac.macports.org/newticket
    Summary: vcs_blackbox @$PORTVERSION Update to latest upstream
    Description: 
New upstream of vcs_blackbox.
github.setup and checksums updated.
    Type: update
    Component: ports
    Port: vcs_blackbox
    Keywords: maintainer haspatch
"
echo 'Attach: Portfile-vcs_blackbox.diff'
