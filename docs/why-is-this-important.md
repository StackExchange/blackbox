Why encrypt your secrets?
=========================

OBVIOUSLY we don't want secret things like SSL private keys and
passwords to be leaked.

NOT SO OBVIOUSLY when we store "secrets" in a VCS repo like Git or
Mercurial, suddenly we are less able to share our code with other
people. Communication between subteams of an organization is hurt. You
can't collaborate as well. Either you find yourself emailing
individual files around (yuck!), making a special repo with just the
files needed by your collaborators (yuck!!), or just deciding that
collaboration isn't worth all that effort (yuck!!!).

The ability to be open and transparent about our code, with the
exception of a few specific files, is key to the kind of collaboration
that DevOps and modern IT practitioners need to do.
