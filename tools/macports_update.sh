#! /usr/bin/env bash

# Generate the Macports Portfile and test it.

PORTVERSION=${1?"Arg 1 must be a version number like 1.20150222 (with no v)"} ; shift

tools/mk_portfile.sh tools/Portfile.template Portfile "$PORTVERSION"

echo
echo 'Adding file:///var/tmp/ports to the start of /opt/local/etc/macports/sources.conf'
echo
fgrep >/dev/null -x 'file:///var/tmp/ports' /opt/local/etc/macports/sources.conf || sudo sed -i -e '1s@^@file:///var/tmp/ports\'$'\n@' /opt/local/etc/macports/sources.conf

echo
echo 'Installing port using local repo.'
echo
sudo port uninstall vcs_blackbox
rm -rf /var/tmp/ports
mkdir -p /var/tmp/ports/security/vcs_blackbox
cp Portfile /var/tmp/ports/security/vcs_blackbox
cd /var/tmp/ports && portindex
sudo port clean --all vcs_blackbox
sudo port install vcs_blackbox

#echo
#echo 'Removing file:///var/tmp/ports from /opt/local/etc/macports/sources.conf'
#echo
#sudo sed -i -e '\@^file:///var/tmp/ports@d' /opt/local/etc/macports/sources.conf
echo 'You may wish to: sudo vi /opt/local/etc/macports/sources.conf'
