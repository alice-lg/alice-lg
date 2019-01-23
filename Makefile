
#
# Build the Alice Looking Glass
# -----------------------------
#
#

PROG=alice-lg
ARCH=amd64

SYSTEM_INIT=systemd

# == END BUILD CONFIGURATION ==

VERSION=$(shell cat ./VERSION)

# Specify build server for remotely building the RPM
# you can do this when you invoke the makefile
# using:
#    make remote_rpm BUILD_SERVER=build-rpm.example.com
BUILD_SERVER=''

DIST=DIST/
REMOTE_DIST=$(PROG)-$(DIST)

RPM=$(PROG)-$(VERSION)-1.x86_64.rpm

LOCAL_RPMS=RPMS

all: alice

client_dev:
	$(MAKE) -C client/

client_prod:
	$(MAKE) -C client/ client_prod

backend_dev: client_dev
	$(MAKE) -C backend/


backend_prod: client_prod
	$(MAKE) -C backend/ bundle
	$(MAKE) -C backend/ linux


alice: client_prod backend_prod
	mv backend/alice-lg-* bin/


dist: clean alice

	mkdir -p $(DIST)opt/alice-lg/alice-lg/bin
	mkdir -p $(DIST)etc/alice-lg

	# Adding post install script
	cp install/scripts/after_install $(DIST)/.

ifeq ($(SYSTEM_INIT), systemd)
	# Installing systemd services
	mkdir -p $(DIST)usr/lib/systemd/system/
	cp install/systemd/* $(DIST)usr/lib/systemd/system/.
else
	# Installing upstart configuration
	mkdir -p $(DIST)/etc/init/
	cp install/upstart/* $(DIST)etc/init/.
endif

	# Copy example configuration
	cp etc/alice-lg/alice.conf.example $(DIST)/etc/alice-lg/alice.conf.example

	# Copy application
	cp bin/$(PROG)-linux-$(ARCH) DIST/opt/alice-lg/alice-lg/bin/.


rpm: dist

	# Clear tmp failed build (if any)
	mkdir -p $(LOCAL_RPMS)

	# Create RPM from dist
	fpm -s dir -t rpm -n $(PROG) -v $(VERSION) -C $(DIST) \
		--architecture $(ARCH) \
		--config-files /etc/alice-lg/alice.conf.example \
		--after-install $(DIST)/after_install \
		opt/ etc/

	mv $(RPM) $(LOCAL_RPMS)


build_server:
ifeq ($(BUILD_SERVER), '')
	$(error BUILD_SERVER not configured)
endif

remote_rpm: build_server dist

	mkdir -p $(LOCAL_RPMS)

	# Copy distribution to build server
	ssh $(BUILD_SERVER) -- rm -rf $(REMOTE_DIST)
	scp -r $(DIST) $(BUILD_SERVER):$(REMOTE_DIST)
	ssh $(BUILD_SERVER) -- fpm -s dir -t rpm -n $(PROG) -v $(VERSION) -C $(REMOTE_DIST) \
		--architecture $(ARCH) \
		--config-files /etc/alice-lg/alice.conf.example \
		--after-install $(REMOTE_DIST)/after_install \
		opt/ etc/

	# Get rpm from server
	scp $(BUILD_SERVER):$(RPM) $(LOCAL_RPMS)/.


clean:
	rm -f bin/alice-lg-linux-amd64
	rm -f bin/alice-lg-osx-amd64
	rm -rf $(DIST)
