SHELL=/bin/sh
BIN=tools

PKGNAME=stack_blackbox

all:
	@echo 'Menu:'
	@echo '  make packages          Make all RPM packages'
	@echo '  make install           (incomplete)

packages:
	PKGRELEASE="$${PKGRELEASE}" $(BIN)/build_rpm.sh stack_blackbox tools/rpm_filelist.txt

install:
	@echo 'To install, copy the files from bin to somewhere in your PATH.'
	@echo 'Or, if you use RPMs, "make packages" and install the result.'

packages-debug:
	@echo BUILD:
	@PKGRELEASE=99 make packages
	@echo ITEMS TO BE PACKAGED:
	find ~/rpmbuild-$(PKGNAME)/installroot -type f
	@echo ITEMS ACTUALLY IN PACKAGE:
	@rpm -qpl $$(cat ~/rpmbuild-$(PKGNAME)/bin-packages.txt)

local:
	@PKGRELEASE=1 make packages
	-@sudo rpm -e $(PKGNAME)
	sudo rpm -i $$(cat ~/rpmbuild-$(PKGNAME)/bin-packages.txt)

lock:
	sudo yum versionlock add $(PKGNAME)

unlock:
	sudo yum versionlock clear

test:
	echo "You don't want to run this."
	exit 1
	pkill gpg-agent ; rm -rf /tmp/tmp.*
	export PATH=~/gitwork/blackbox/bin:/usr/local/bin:/bin:/usr/bin:/usr/local/sbin:/usr/sbin:/sbin ; \
		cd ~/gitwork/blackbox && tools/confidence_test.sh
	@if [[ -e ~/.gnupg ]]; then echo ERROR: '~/.gnupg' should not exist. If it does, this means test test suite may be poluting your actual .gnupg configuration. ; false ; fi
