Code Overview
=============

Here is how the code is laid out.

TODO(tlim): Add a diagram of the layers

```
cmd/blackbox/   The command line tool.
             blackbox.go   main()
             cli.go        Definition of all subcommands and flags
             drive.go      Processes flags and calls functions in verbs.go
                   NOTE: These are the only files that are aware of the
                         flags.  Everything else gets the flag data passed to it
                         as a parameter. This way the remaining system can be
                         used as a module.

pkg/box/        High-level functions related to "the black box".
        verbs.go       One function per subcommand.
        box.go         Functions for manipulating the files in .blackbox
        boxutils.go    Helper functions for the above.

pkg/bblog/      Module that provides logging facilities.
pkg/bbutil/     Functions that are useful to box, plug-ins, etc.
pkg/tainedname/ Module for printing filenames escaped for Bash.

models/vcs.go        The interface that defines a VCS plug-in.
models/crypters.go   The interface that defines a GPG plug-in.

pkg/crypters/   Plug-ins for GPG functionality.
pkg/crypters/gnupg   Plug-in that runs an external gpg binary (found via $PATH)

pkg/vcs/        Plug-ins for VCS functionality.
pkg/vcs/none        Repo-less mode.
pkg/vcs/git         Git mode.
```
