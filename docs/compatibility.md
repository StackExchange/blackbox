Compatibility
=============

# Compatibility with Blackbox v1

The command names all changed from v1 to v2.  The `binv2` directory
includes shell scripts that provide full backwards compatibility.

# Supported Architectures

Blackbox supports a plug-in archtecture to easily support multiple VCS
system.  Current support is for:

## Supported VCS/DVCS systems

* git
* "none" (repo-less use is supported)
* WOULD LOVE VOLUNTEERS TO HELP ADD SUPPORT FOR: hg, svn, p4

## Supported GPG versions

* Git 1.x and 2.0
* Git 2.2 and higher
* WOULD LOVE VOLUNTEERS TO HELP ADD SUPPORT FOR:
  golang.org/x/crypto/openpgp (this would make the code have no
  external dependencies)

## Supported Operating systems

Blackbox should work on any Linux system with GnuPG installed.
Blackbox simply looks for `gpg` in `$PATH`.

Windows: It should work (but has not been extensively tested) on
Windows WSL2.

# Automated testing

While many combinations work, we do automated tests
on these combinations.  If any of these fail it blocks the release:

* macOS: GnuPG 2.2 executables from https://gpgtools.org/
* CentOS: GnuPG 2.0.x executables from the "base" or "updates" repo.

Windows native: VOLUNTEER NEEDED to make a native Windows version
(should be rather simple as Go does most of the work)

NOTE: Version 1 worked on CentOS/RedHat, macOS, Gygwin, WinGW, NetBSD,
and SmartOS.  Hopefully we can achieve that broad level of support in
the future.  Any system that is supported by the Go language and
has GuPG 2.0.x or higher binaries available should be easy to achieve.
We'd also like to have automated testing for the same.

# Windows Support

BlackBox assumes that `blackbox-admins.txt` and `blackbox-files.txt` will have
LF line endings. Windows users should be careful to configure Git or other systems
to not convert or "fix" those files.

If you use Git, add the following lines to your `.gitattributes` file:

    **/blackbox-admins.txt text eol=lf
    **/blackbox-files.txt text eol=lf

The `blackbox init` (and newer versions of `blackbox_initialize`)
will create an appropriate `.gitattributes` file for you.

# Cygwin

TODO: List what packages are required for building the software.

TODO: List what packages are required for running the software.


# MinGW

MinGW (comes with Git for Windows) support requires the following:

TODO: FILL IN any requirements
