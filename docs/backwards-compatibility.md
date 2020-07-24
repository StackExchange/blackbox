Backwards Compatibility
=======================

# Where is the configuration stored? .blackbox vs. keyrings/live

Blackbox stores its configuration data in the `.blackbox` subdirectory.  Older
repos use `keyrings/live`.  For backwards compatibility either will work.

All documentation refers to `.blackbox`.

You can convert an old repo by simply renaming the directory:

```
mv keyrings/live .blackbox
rmdir keyrings
```

There is no technical reason to convert old repos except that it is less
confusing to users.

This change was made in commit 60e782a0, release v1.20180615.


# How blackbox fines the config directory:

## Creating the repo:

`blackbox init` creates the config directory in the root
of the repo.  Here's how it picks the name:

- If `$BLACKBOX_TEAM` is set, `.blackbox-$BLACKBOX_TEAM` is used.
- If the flag `--team <teamname>` is set, it uses `.blackbox-<teamname>`
- Otherwise, it uses `.blackbox`

When searching for the configuration directory, the following
locations are checked. First match wins.

- `.blackbox-$BLACKBOX_TEAM` (only if `$BLACKBOX_TEAM` is set)
- The value of `--config value` (if the flag is set)
- `$BLACKBOX_CONFIGDIR` (the preferred env. variable to use)
- `$BLACKBOXDATA` (for backwards compatibility with v1)
- `.blackbox`
- `keyrings/live` (for backwards compatibility)

NOTE: The env variables and `--config` should be set to the full path
to the config directory (i.e.: `/Users/tom/gitstuff/myrepo/.blackbox`).
If it is set to a relative directory (i.e. `.blackbox` or
`../myrepo/.blackbox`) most commands will break.

NOTE: Why the change from `$BLACKBOXDATA` to `$BLACKBOX_CONFIGDIR`?  We want
all the env. variables to begin with the prefix `BLACKBOX_`.  If v1
supported another name, that is still supported. If you are starting
with v2 and have no other users using v1, please use the `BLACKBOX_`
prefix.

