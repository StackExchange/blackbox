Release v1.20220610

NOTE: I don't have a lot of time to commit to this project.  I'd gladly accept help, especially
with improving the testing on various operating systems.

Major feature: macOS users rejoice!  Incompatibility with macOS Monterey 12.3 is fixed! (#347)

* Add .gitattributes during repo initialization (#352)
* Update zgen reference to zgenom (#350)
* Improve test data generation (#348)
* Fix 'chmod' for macOS Monterey 12.3 (#347)


Release v1.20200429

NOTE: While there is now support for NetBSD and SunOS/SmartOS, the
release process only tests on macOS and CentOS7 because that's all I
have access to.

* Fix tools that break when ".." or "." are used in a path (#304)
* Respect PREFIX variable for copy-install (#294)
* Documentation: Add pkgsrc install instructions (#292)
* Improve support for Windows (#291)
* Clarify gpg version usage (#290)
* Many documentation fixes
* DOCUMENTATION: Promote 'getting started' to a section, enumerate steps (#283)
* Commit changes to gitignore when deregistering (#282)
* Add support for NetBSD and SunOS (SmartOS)
* Defend against ShellShock


Release v1.20181219

* New OS support: Add support for NetBSD and SunOS (SmartOS)
* Testing: Improve confidence test.
* .blackbox is now the default config directory for new repos. (#272)
* Add blackbox_decrypt_file (#270)
* Improved compatibility: change"/bin/[x]" to "/usr/bin/env [x]" (#265)
* Add blackbox_less. (#263)
* add nix method of install (#261)
* Linked setting up of GPG key (#260)


Release v1.20180618

* Restore `make manual-install` with warning. (#258)

Release v1.20180615

* Standardize on .blackbox for config. Use keyrings/live for backwards compatibility.
* Store keys in .blackbox directory (#218)
* Suggest committing changes to pubring.gpg when running blackbox_removeadmin (#248)
* Fix typo (#246)
* Improve installation instructions (#244)
* Fix replacing-expired-keys link in README (#241)
* Fix problems when gpg2 is installed next to gpg (#237)
* Many documentation corrections, updates, etc.
* Exclude default keyring from import (#223)
* .gitattributes not always updated (PR#146)
* Fix bugs related to updating .gitattributes (PR#146)
* Update readme with CircleCI link (#216)
* Run the tests on a CI (#215)
* Fixed Alpine compatibility (chmod) (#212)
* direct repobase message to stderr (#204)
* Improve Windows compatibility
* NEW: .gitattributes Set Unix-only files to eol=lf
* Silence 'not changed' output during keychain import (#200)
* Improve FreeBSD compatibility
* shred_file() outputs warning message to stderr. (#192)
* Don't complain about GPG_AGENT_INFO if using newer gpg-agent (#189)
* [FreeBSD] Fix use of chmod (#180)
* Requiring a file to be entered to finish editing (#175)
* Remove the key from the keyring when removing an admin (#173)
* Add FreeBSD support (#172)
* Add list admins commandline tool. (#170)
ignore backup files and secring.gpg in $BLACKBOXDATA (#169)
Allow parallel shredding of files (#167)
* Add/improve Mingw support
* Make "make confidence" less fragile
* And a lot, lot more.

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

Release v1.20170127

* Starting CHANGELOG.
