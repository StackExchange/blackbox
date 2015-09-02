SHELL=/bin/sh

PKGNAME=stack_blackbox

all:
	@echo 'Menu:'
	@echo '  make update            Update any generated files'
	@echo '  make packages          Make RPM packages'
	@echo '  make packages-deb      Make DEB packages'
	@echo '  make install           (incomplete)'

install:
	@echo 'To install, copy the files from bin to somewhere in your PATH.'
	@echo 'Or, if you use RPMs, "make packages" and install the result.'

# The default package type is RPM.
packages: packages-rpm

#
# RPM builds
#

# NOTE: mk_rpm_fpmdir.stack_blackbox.txt is the master list of files.  All
# other packages should generate their list from it.

packages-rpm:
	cd tools && PKGRELEASE="$${PKGRELEASE}" PKGDESCRIPTION="Safely store secrets in git/hg/svn repos using GPG encryption" ./mk_rpm_fpmdir stack_blackbox mk_rpm_fpmdir.stack_blackbox.txt

packages-rpm-debug:
	@echo BUILD:
	@PKGRELEASE=99 make packages
	@echo ITEMS TO BE PACKAGED:
	find ~/rpmbuild-$(PKGNAME)/installroot -type f
	@echo ITEMS ACTUALLY IN PACKAGE:
	@rpm -qpl $$(cat ~/rpmbuild-$(PKGNAME)/bin-packages.txt)

local-rpm:
	@PKGRELEASE=1 make packages
	-@sudo rpm -e $(PKGNAME)
	sudo rpm -i $$(cat ~/rpmbuild-$(PKGNAME)/bin-packages.txt)

lock-rpm:
	sudo yum versionlock add $(PKGNAME)

unlock-rpm:
	sudo yum versionlock clear

#
# Manual install
#
manual-install:
	@echo 'Symlinking files from ./bin to /usr/local/bin'
	@cd bin && for f in `find . -type f -iname "*" ! -iname "Makefile"`; do ln -fs `pwd`/$$f /usr/local/bin/$$f; done
	@echo 'Done.'

manual-uninstall:
	@echo 'Removing blackbox files from /usr/local/bin'
	@cd bin && for f in `find . -type f -iname "*" ! -iname "Makefile"`; do rm /usr/local/bin/$$f; done
	@echo 'Done.'

#
# DEB builds
#

packages-deb:	tools/mk_deb_fpmdir.stack_blackbox.txt
	cd tools && PKGRELEASE="$${PKGRELEASE}" PKGDESCRIPTION="Safely store secrets in git/hg/svn repos using GPG encryption" ./mk_deb_fpmdir stack_blackbox mk_deb_fpmdir.stack_blackbox.txt

# Make mk_deb_fpmdir.vcs_blackbox.txt from mk_rpm_fpmdir.stack_blackbox.txt:
tools/mk_deb_fpmdir.stack_blackbox.txt: tools/mk_rpm_fpmdir.stack_blackbox.txt
	sed -e '/^#/d' -e 's@/usr/blackbox/bin/@/usr/bin/@g' -e '/profile.d-usrblackbox.sh/d' <tools/mk_rpm_fpmdir.stack_blackbox.txt >$@

packages-deb-debug:	tools/mk_deb_fpmdir.stack_blackbox.txt
	@echo BUILD:
	@PKGRELEASE=99 make packages-deb
	@echo ITEMS TO BE PACKAGED:
	find ~/debbuild-$(PKGNAME)/installroot -type f
	@echo ITEMS ACTUALLY IN PACKAGE:
	@dpkg --contents $$(cat ~/debbuild-$(PKGNAME)/bin-packages.txt)

local-deb:
	@PKGRELEASE=1 make packages
	-@sudo dpkg -e $(PKGNAME)
	sudo dpkg -i $$(cat ~/rpmbuild-$(PKGNAME)/bin-packages.txt)

#
# MacPorts builds
#
# To test:
# rm -rf /tmp/foo ; mkdir -p /tmp/foo;make packages-macports DESTDIR=/tmp/foo;find /tmp/foo -ls

# Make mk_macports.vcs_blackbox.txt from mk_rpm_fpmdir.stack_blackbox.txt:
tools/mk_macports.vcs_blackbox.txt: tools/mk_rpm_fpmdir.stack_blackbox.txt
	sed  -e '/^#/d' -e 's@/usr/blackbox/bin/@bin/@g' -e '/profile.d-usrblackbox.sh/d' <tools/mk_rpm_fpmdir.stack_blackbox.txt >$@

# MacPorts expects to run: make packages-macports DESTDIR=${destroot}
packages-macports: tools/mk_macports.vcs_blackbox.txt
	mkdir -p $(DESTDIR)/bin
	cd tools && ./mk_macports mk_macports.vcs_blackbox.txt

# stow is a pretty easy way to manage simple local installs on GNU systems
install-stow:
	mkdir -p /usr/local/stow/blackbox/bin
	cp bin/* /usr/local/stow/blackbox/bin
	rm /usr/local/stow/blackbox/bin/Makefile
	cd /usr/local/stow; stow -R blackbox
uninstall-stow:
	cd /usr/local/stow; stow -D blackbox
	rm -rf /usr/local/stow/blackbox

# Add other package types here.

#
# Updates
#
update: tools/mk_deb_fpmdir.stack_blackbox.txt tools/mk_macports.vcs_blackbox.txt

clean:
	rm tools/mk_deb_fpmdir.stack_blackbox.txt tools/mk_macports.vcs_blackbox.txt

#
# System Test:
#
confidence:
	@if [ -e ~/.gnupg ]; then echo ERROR: '~/.gnupg should not exist. If it does, bugs may polute your .gnupg configuration. If the code has no bugs everything will be fine. Do you feel lucky?'; false ; fi
	@if which >/dev/null gpg-agent ; then pkill gpg-agent ; rm -rf /tmp/tmp.* ; fi
	@export PATH="$(PWD)/bin:/usr/local/bin:/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/sbin:/opt/local/bin:$(PATH)" ; tools/confidence_test.sh
		tools/confidence_test.sh
	@if which >/dev/null gpg-agent ; then pkill gpg-agent ; fi
	@if [ -e ~/.gnupg ]; then echo ERROR: '~/.gnupg was created which means the scripts might be poluting GnuPG configuration.  Fix this bug.'; false ; fi
