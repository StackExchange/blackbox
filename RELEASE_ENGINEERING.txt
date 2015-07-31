# Branches and Tags:

There are 3 branches/tags:

* **HEAD:** The cutting edge of development.
* **tag stable:** Stable enough for use by most people.
* **tag production:** Burned in long enough that we are confident it can be widely adopted.

If you are packaging Blackbox for distribution, you should track the *tag production*.  You might also want to provide a separate package that tracks *tag stable:* for early adopters.

# Build Tasks

# Stable Releases

Marking the software to be "stable":

Step 1. Tag it.

```
git pull
git tag -d stable
git push origin :stable
git tag stable
git push origin tag stable
```

Step 2. Mark your calendar 1 week from today to check
to see if this should be promoted to production.


# Production Releases

If no bugs have been reported a full week after a stable tag has been pushed, mark the release to be "production".

```
git fetch
git checkout stable
git tag -d production
git push origin :production
git tag production
git push origin tag production
R="v1.$(date +%Y%m%d)"
git tag "$R"
git push origin tag "$R"
```

# Updating MacPorts (automatic)

Step 1: Generate the Portfile

```
tools/macports_report_upgrade.sh  1.20150222
```

This script will generate a file called `Portfile-vcs_blackbox.diff` and instructions on how to submit it as a update request.

Step 2: Submit the update request.

Submit the diff file as a bug as instructed. The instructions should look like this:

* PLEASE OPEN A TICKET WITH THIS INFORMATION:
    https://trac.macports.org/newticket
* Summary: `vcs_blackbox @1.20150222 Update to latest upstream`
* Description: ```New upstream of vcs_blackbox.
github.setup and checksums updated.```
* Type: `update`
* Component: `ports`
* Port: `vcs_blackbox`
* Keywords: `maintainer haspatch`
* Attach this file: `Portfile-vcs_blackbox.diff`

Step 3: Watch for the update to happen.

# Updating MacPorts (manual)

This is the old, manual, procedure.  If the automated procedure fails to work, these notes may or may not be helpful.

The ultimate result of the script should be the output of `diff -u Portfile.orig Portfile` which is sent as an attachment to MacPorts.  The new `Portfile` should have these changes:

1. The `github.setup` line should have a new version number.
2. The `checksums` line(s) should have updated checksums.

How to generate the checksums?

The easiest way is to to make a Portfile with incorrect checksums, then run `sudo port -v checksum vcs_blackbox` to see what they should have been.  Fix the file, and try again until the checksum command works.

Next run `port lint vcs_blackbox` and make sure it has no errors.

Some useful commands:

Change repos in sources.conf:
```
sudo vi /opt/local/etc/macports/sources.conf
  Add this line early in the file:
  file:///var/tmp/ports
```

Add a local repo:
```
fgrep >/dev/null -x 'file:///var/tmp/ports' /opt/local/etc/macports/sources.conf || sudo sed -i -e '1s@^@file:///var/tmp/ports\'$'\n@' /opt/local/etc/macports/sources.conf
```

Remove the local repo:
```
sudo sed -i -e '\@^file:///var/tmp/ports@d' /opt/local/etc/macports/sources.conf
```

Test a Portfile:
``` 
sudo port uninstall vcs_blackbox
sudo port clean --all vcs_blackbox
rm -rf ~/.macports/opt/local/var/macports/sources/rsync.macports.org/release/tarballs/ports/security/vcs_blackbox/
rm -rf /var/tmp/ports
mkdir -p /var/tmp/ports/security/vcs_blackbox
cp Portfile /var/tmp/ports/security/vcs_blackbox
cd /var/tmp/ports && portindex
sudo port -v checksum vcs_blackbox
sudo port install vcs_blackbox
```
