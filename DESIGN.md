BlackBox Internals
==================

The goal of the Go rewrite is to improve the usability and
maintainability of Blackbox, meanwhile make it easier to implement new 

The system is built in distinct layers: view, controller, model.

Suppose there is a subcommand "`foo`".  `blackbox.go` parses the
user's command line args and calls `cmdFoo()`, which is given
everything it needs to do the operation.  For example, it is given the
filenames the user specified exactly; even if an empty list means "all
files", at this layer the empty list is passed to the function.

`cmdFoo()` contains the business logic of how the operation should be
done: usually iterating over filenames and calling verb(s) for each
one.  For example if an empty file list means "all files", this is the
layer that enumerates the files.

`cmdFoo()` is implemented in the file `cmd_foo.go`.  The caller of
`cmdFoo()` should provide all data it needs to get the job done.
`cmdFoo()` doesn't refer to global flags, they are passed to the
function as parameters.  Therefore the function has zero side-effects
(except possibly logging) and can be called as library functions by
other systems.  This is the external (binary) API which should be
relatively stable.

`cmdFoo()` calls verbs that are in `bbutil/`.  Some of those verbs are
actually interfaces. For example, any VCS-related verbs are actually a
Go interface which might be implemented one of many ways (Git,
Subversion, Mercurial), GPG-functions may be implemented by shelling
out to `gpg.exe` or by using Go's gpg library.

They layers look like this:

| View | `blackbox.go` | Parses User Commands, calls controller |
| Controller | `cmd_*.go` | The business logic. Iterates and calls verbs |
| Model | `pkg/bbutil` | Verbs |
| Interfaces | `pkg/*` | Interfaces and their implementations |

At least that's the goal.  We'll see how well we can achieve this.
