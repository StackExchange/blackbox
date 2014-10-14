#!/bin/bash

# build_rpm.sh - Build an RPM of these files.  (uses FPM)

# Usage:
#   make_rpm.sh PACKAGENAME MANIFEST1 MANIFEST2 ...

# Example:
#   Make a package foopkg manifest.txt
# Where "manifest.txt" contains:
#   exec /usr/bin/foo        foo/foo
#   exec /usr/bin/bar        bar/bar.sh
#   read /usr/man/man1/bar.1 bar/bar.1.man
#   0444 /etc/foo.conf       bar/foo.conf
#
# Col1  chmod-style permissions or "exec" for 0755, "read" for 0744.
# Col2  Installation location.
# Col3  Source of the file.

set -e

# Parameters for this RPM:
PACKAGENAME=${1?"First arg must be the package name."}
shift

# Defaults that can be overridden via env variables:
# All packages are 1.0 unless otherwise specifed:
: ${PKGVERSION:=1.0} ;
# If there is no iteration setting, assume "1":
: ${PKGRELEASE:=1}

# The RPM is output here: (should be a place that can be wiped)
OUTPUTDIR="${HOME}/rpmbuild-$PACKAGENAME"
# Our build system expects to find the list of artifacts here:
RPM_BIN_LIST="${OUTPUTDIR}/bin-packages.txt"

# -- Now the real work can be done.

# Clean the output dir.
rm -rf "$OUTPUTDIR"
mkdir -p "$OUTPUTDIR/installroot"

# Copy the files into place:
cat """$@""" | grep -v '^$' | while read -a arr ; do
  PERM="${arr[0]}"
  DEST="${arr[1]}"
  SRC="${arr[2]}"
  echo ========== "$PERM $DEST"
  case $PERM in
    \#*)  continue ;;   # Skip comments.
    exec) PERM=0755 ;;
    read) PERM=0744 ;;
    *) ;;
  esac
  FULLDEST="$OUTPUTDIR/installroot/${arr[1]}"
  install -D -T -b -m "$PERM" -T "$SRC" "$FULLDEST"
done

# Build the RPM:
cd "$OUTPUTDIR" && fpm -s dir -t rpm \
  -a x86_64 \
  --epoch '0' \
  -n "${PACKAGENAME}" \
  --version "${PKGVERSION}" \
  --iteration "${PKGRELEASE}" \
  --description 'Safely store secrets in Git/Hg repos using GPG encryption' \
  -C "$OUTPUTDIR/installroot" \
  .

# Our build system expects to find the list of all packages created
# in bin-packages.txt.  Generate that list:
find "$OUTPUTDIR" -maxdepth 1 -name '*.rpm' >"$RPM_BIN_LIST"
# Output the list for debugging purposes:
echo ========== "$RPM_BIN_LIST"
cat "$RPM_BIN_LIST"
