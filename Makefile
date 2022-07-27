
#
# Build the Alice Looking Glass locally
#

PROG=alice-lg
ARCH=amd64

VERSION=$(shell cat ./VERSION)


all: alice

test: ui_test backend_test

alice: ui backend
	cp cmd/alice-lg/alice-lg-* bin/

ui:
	$(MAKE) -C ui/

ui_test:
	$(MAKE) -C ui/ test

backend:
	$(MAKE) -C cmd/alice-lg/ static

backend_test:
	go test ./pkg/...


clean:
	rm -f bin/alice-lg-linux-amd64
	rm -f bin/alice-lg-osx-amd64
	rm -rf $(DIST)


.PHONY: backend ui clean

