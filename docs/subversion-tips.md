Subversion Tips
===============

NOTE: This is from v1.  Can someone that uses Subversion check
this and update it?


The current implementation will store the blackbox in `/keyrings` at
the root of the entire repo.  This will create an issue between
environments that have different roots (i.e. checking out `/` on
development vs `/releases/foo` in production). To get around this, you
can `export BLACKBOX_REPOBASE=/path/to/repo` and set a specific base
for your repo.

This was originally written for git and supports a two-phase commit,
in which `commit` is a local commit and "push" sends the change
upstream to the version control server when something is registered or
deregistered with the system.  The current implementation will
immediately `commit` a file (to the upstream subversion server) when
you execute a `blackbox_*` command.

