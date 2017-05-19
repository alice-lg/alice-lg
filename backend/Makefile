#
# Ecix Alice Makefile
#

PROG=alice-lg
ARCH=amd64

APP_VERSION=$(shell cat ../VERSION)
VERSION=$(APP_VERSION)_$(shell git rev-parse --short HEAD)

BUILD_SERVER=''

SYSTEM_INIT=systemd

DIST=DIST/
REMOTE_DIST=$(PROG)-$(DIST)

RPM=$(PROG)-$(VERSION)-1.x86_64.rpm

LOCAL_RPMS=RPMS

# OS Detection
UNAME=$(shell uname)
ifeq ($(UNAME), Darwin)
  TARGET=osx
else
  TARGET=linux
endif


LDFLAGS=-ldflags="-X main.version=$(APP_VERSION)"

all: $(TARGET)
	@echo "Built $(VERSION) @ $(TARGET)"

osx:
	GOARCH=$(ARCH) GOOS=darwin go build $(LDFLAGS) -o $(PROG)-osx-$(ARCH)

linux:
	GOARCH=$(ARCH) GOOS=linux go build $(LDFLAGS) -o $(PROG)-linux-$(ARCH)


build_server:
ifeq ($(BUILD_SERVER), '')
	$(error BUILD_SERVER not configured)
endif

dist: clean linux

	mkdir -p $(DIST)opt/ecix/birdwatcher/bin
	mkdir -p $(DIST)etc/ecix

ifeq ($(SYSTEM_INIT), systemd)
	# Installing systemd services
	mkdir -p $(DIST)usr/lib/systemd/system/
	cp install/systemd/* $(DIST)usr/lib/systemd/system/.
else
	# Installing upstart configuration
	mkdir -p $(DIST)/etc/init/
	cp install/upstart/init/* $(DIST)etc/init/.
endif


	# Copy config and startup script
	cp etc/ecix/* DIST/etc/ecix/.
	rm -f DIST/etc/ecix/*.local.*

	# Copy bin
	cp $(PROG)-linux-$(ARCH) DIST/opt/ecix/birdwatcher/bin/.


release: linux

	mkdir -p ../birdseye-static/birdwatcher-builds/$(APP_VERSION)/
	cp birdwatcher-linux-amd64 ../birdseye-static/birdwatcher-builds/$(APP_VERSION)/
	rm -f ../birdseye-static/birdwatcher-builds/latest
	cd ../birdseye-static/birdwatcher-builds && ln -s $(APP_VERSION) latest 


rpm: dist

	# Clear tmp failed build (if any)
	mkdir -p $(LOCAL_RPMS)

	# Create RPM from dist
	fpm -s dir -t rpm -n $(PROG) -v $(VERSION) -C $(DIST) \
		--config-files /etc/ecix/birdwatcher.conf \
		opt/ etc/

	mv $(RPM) $(LOCAL_RPMS)


remote_rpm: build_server dist

	mkdir -p $(LOCAL_RPMS)

	# Copy distribution to build server
	ssh $(BUILD_SERVER) -- rm -rf $(REMOTE_DIST)
	scp -r $(DIST) $(BUILD_SERVER):$(REMOTE_DIST)
	ssh $(BUILD_SERVER) -- fpm -s dir -t rpm -n $(PROG) -v $(VERSION) -C $(REMOTE_DIST) \
		--config-files /etc/ecix/birdwatcher.conf \
		opt/ etc/

	# Get rpm from server
	scp $(BUILD_SERVER):$(RPM) $(LOCAL_RPMS)/.


clean:
	rm -f $(PROG)-osx-$(ARCH)
	rm -f $(PROG)-linux-$(ARCH)
	rm -rf $(DIST)
