# Ideas for BlackBox Version 2

I'm writing this to solicit feedback and encourage discussion.

Here are my thoughts on a "version 2" of BlackBox.  This is where
I list ideas that would require major changes to the system. They
might break backwards compatibility, though usually not.

BlackBox grew from a few simple shell scripts used at StackOverflow.com
to a larger system used by dozens (hundreds?) of organizations. Not
all the design decisions were "forward looking".

These are the things I'd like to change someday.

[TOC]

## Change the commmand names

There should be one program, with subcommands that have names that make more sense:

* `blackbox init`
* `blackbox register <filename> <...>`
* `blackbox deregister <filename> <...>`
* `blackbox edit <filename> <...>`
* `blackbox decrypt <filename> <...>`
* `blackbox encrypt <filename> <...>`
* `blackbox decrypt_all`
* `blackbox addadmin <key>`
* `blackbox removeadmin <key>`
* `blackbox cat <filename> <...>`
* `blackbox diff <filename> <...>`
* `blackbox list_files`
* `blackbox list_admins`
* `blackbox shred_all`
* `blackbox update_all`
* `blackbox whatsnew`

Backwards compatibility: The old commands would simply call the new commands.

## Change the "keyrings" directory

The name "keyrings" was unfortunate.  First, it should probably begin with a ".".  Second, it stores more than just keyrings.  Lastly, I'm finding that in most cases we want many repos to refer to the same keyring, which is not supported very well.

A better system would be:

1. If `$BLACKBOX_CONFIG` is set, use that directory.
2. If the repo base directory has a file called ".blackbox_external", read that file as if you are reading `$BLACKBOX_CONFIG`
3. If the repo base directory has a "keyrings" directory, use that.
4. If the repo base directory has a ".blackboxconfig" directory, use that.

Some thoughts on .blackbox_external:
I'm not sure what the format should be, but I want it to be simple and expandable.  It should support support "../../dir/name" and "/long/path".  However some day we may want to include a Git URL and have the system automatically get the keychain from it. That means the format has to be something like directory:../dir/name so that later we can add git:the_url.


Backwards compatibility: "keyrings" would be checked before .blackbox

## Repo-less mode

I can't imagine storing files that aren't in a repo. I just put everything in repos lately. I use it more than I use NFS.  That said, I have received feedback that people would like the ability to disable automatic committing of files.

I prefer the file commits to be automatic because when they were manual, people often accidentally committed the plaintext file instead of the GPG file.  Fixing such mistakes is a PITA and, of yeah, a big security nightmare.

That said, I'm willing to have a "repo-less" mode.

When this mode is triggered, no add/commit/ignore tasks are done.  The search for the keyrings directory still uses `$BLACKBOX_CONFIG` but if that is unset it looks for .blackbox_config in the current directory, then recursively ".." until we hit "/".

I think (but I'm not sure) this would benefit the entire system because it would force us to re-think what VCS actions are done when.

I think (but I'm not sure) that a simple way to implement this would be to add an environment variable that overrides the automatic VCS detection. When set to "none", all VCS operations would basically become no-ops.  (This could be done by writing a plug-in that does nothing for all the vcs_* calls)

Backwards compatibility: This would add a "none" VCS, not remove any existing functionality.


## Is "bash" the right language?

`bash` is fairly universal. It even exists on Windows.  However it is not the right language for large systems. Writing the acceptance tests is quite a bear.  Managing ".gitignore" files in bash is impossible and the current implementation fails in many cases.

`python` is my second favorite language. It would make the code cleaner and more testable. However it is not installed everywhere.  I would also want to write it in Python3 (why start a new project in Python2?) but sadly Python3 is less common.  It is a chicken vs. egg situation.

`go` is my favorite language. I could probably rewrite this in go in a weekend. However, now the code is compiled, not interpreted. Therefore we lose the ability to just "git clone" and have the tools you want.  Not everyone has a Go compiler installed on every machine.

The system is basically unusable on Windows without Cygwin or MINGW.  A rewrite in python or go would make it work better on Windows, which currently requires Cygwin or MinGW (which is a bigger investment than installing Python). On the other hand, maybe Ubuntu-on-Windows makes that a non-issue.

As long as the code is in `bash` the configuration files like `blackbox-files.txt` and `blackbox-admins.txt` have problems.  Filenames with carriage returns aren't supported.  If this was in Python/Go/etc. those files could be json or some format with decent quoting and we could handle funny file names better. On the other hand, maybe it is best that we don't support funny filenames... we shouldn't enable bad behavior.

How important is itto blackbox users that the system is written in "bash"?


## ditch the project and use git-crypt

People tell me that git-crypt is better because, as a plug-in, automagically supports "git diff", "git log" and "git blame".

However, I've never used it so I don't have any idea whether git-crypt is any better than blackbox.

Of course, git-crypt doesn't work with SVN, HG, or any other VCS.  Is blackbox's strong point the fact that it support so many VCS systems?  To be honest, it originally only supported HG and GIT because I was at a company that used HG but then changed to GIT.  Supporting anything else was thanks to contributors. Heck, HG support hasn't even been tested recently (by me) since we've gone all git where I work.

How important is this to BlackBox users?
