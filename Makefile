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
