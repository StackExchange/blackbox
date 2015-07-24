#!/usr/bin/env bash

#
# _blackbox_common_test.sh -- Unit tests of functions from _blackbox_common.sh
#

set -e
. "${0%/*}/_blackbox_common.sh"
. tools/test_functions.sh

PHASE 'Test cp-permissions: TestA'
touch TestA TestB TestC TestD
chmod 0347 TestA
chmod 0700 TestB
chmod 0070 TestC
chmod 0070 TestD
cp_permissions TestA TestB TestC
# NOTE: cp_permissions is not touching TestD.
assert_file_perm '--wxr--rwx' TestA
assert_file_perm '--wxr--rwx' TestB
assert_file_perm '--wxr--rwx' TestC
assert_file_perm '----rwx---' TestD  # TestD doesn't change.
rm -f TestA TestB TestC TestD

echo '========== DONE.'
