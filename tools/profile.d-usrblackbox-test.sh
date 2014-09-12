#!/bin/sh

# Test profile.d-usrblackbox.sh

# Make sure profile.d-usrblackbox.sh works.

# Test variations including /usr/blackbox/bin is not in the path, is
# already there in the front, middle, or end, and tests if the path has :
# in weird places (front, middle, or both).

# To run the test:
# bash tools/profile.d-usrblackbox-test.sh | fgrep --color /usr/blackbox/bin
# sh tools/profile.d-usrblackbox-test.sh | fgrep --color /usr/blackbox/bin

# NOTE: profile.d-usrblackbox.sh is written to be so small that it fits as an "inline" file.
# https://lwn.net/Articles/468678/
# To remove the last newline in the file:
#  perl -i -pe 'chomp if eof' profile.d-usrblackbox.sh

for p in \
  '/usr/local/bin:/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/sbin'                     \
  '/usr/blackbox/bin:/usr/local/bin:/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/sbin'   \
  '/usr/local/bin:/bin:/usr/blackbox/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/sbin'   \
  '/usr/local/bin:/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/sbin:/usr/blackbox/bin'   \
  '/Apple spaces/local/bin:/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/sbin'            \
  ; do

  export PATH="$p"
  . tools/profile.d-usrblackbox.sh
  echo NEW: "$PATH"

  export PATH=":$p"
  . tools/profile.d-usrblackbox.sh
  echo NEW: "$PATH"

  export PATH="$p:"
  . tools/profile.d-usrblackbox.sh
  echo NEW: "$PATH"

  export PATH=":$p:"
  . tools/profile.d-usrblackbox.sh
  echo NEW: "$PATH"

done
