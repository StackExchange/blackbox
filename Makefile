SHELL=/bin/sh

PKGNAME=stack_blackbox

all:
	@echo 'Menu:'
	@echo '  make packages          Make RPM packages'
	@echo '  make install           (incomplete)'

install:
	@echo 'To install, copy the files from bin to somewhere in your PATH.'
	@echo 'Or, if you use RPMs, "make packages" and install the result.'

# The default package type is RPM.
packages: packages-rpm

#
# RPM builds
#

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

# Add other package types here.

#
# System Test:
#

confidence:
	@if [[ -e ~/.gnupg ]]; then echo ERROR: '~/.gnupg should not exist. If it does, bugs may polute your .gnupg configuration. If the code has no bugs everything will be fine. Do you feel lucky?'; false ; fi
	@if which >/dev/null gpg-agent ; then pkill gpg-agent ; rm -rf /tmp/tmp.* ; fi
	@export PATH=~/gitwork/blackbox/bin:/usr/local/bin:/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/sbin ; 
		cd ~/gitwork/blackbox && tools/confidence_test.sh
	@if which >/dev/null gpg-agent ; then pkill gpg-agent ; fi
	@if [[ -e ~/.gnupg ]]; then echo ERROR: '~/.gnupg was created which means the scripts might be poluting GnuPG configuration.  Fix this bug.'; false ; fi
