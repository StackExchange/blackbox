Release v1.20170127

* Starting CHANGELOG.


Release v1.20170309

* "make test" is an alias for "make confidence"
* macOS: make_tempdir must create shorter paths
* Fix "make confidence" for newer version of Git
* README.md: Add info about our new mailing list

Release v1.20170611

* confidence_test.sh verifies external tools exist
* confidence_test.sh more reliable for non-UTF8 users
* "make test" no longer prompts for passwords
* blackbox works better when target directory lives in root (#194)
* Add confidence_test.sh tests for admin operations
* blackbox_list_admins fails (#193)
* confidence_test.sh works better on FreeBSD
* tools/confidence_test.sh: now works with gnupg-2.0 and gnupg-2.1
* Blackbox now officially supports both gnupg-2.0 and gnupg-2.1
* blackbox_shred_all_files: BUGFIX: Does not shred files with spaces
* blackbox_removeadmin: disable gpg's confirmation
* Sync mk_rpm_fpmdir from master
