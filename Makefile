SHELL=/bin/sh
BIN=tools

all:
	@echo 'Menu:'
	@echo '  make packages          Make all RPM packages'
	@echo '  make install           (incomplete)

packages:
	PKGRELEASE="$${PKGRELEASE}" $(BIN)/build_rpm.sh stack_blackbox tools/rpm_filelist.txt

install:
	@echo 'To install, copy the files from bin to somewhere in your PATH.'
	@echo 'Or, if you use RPMs, "make packages" and install the result.'
