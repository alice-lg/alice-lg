
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
	mkdir -p ./ui/build
	touch ./ui/build/UI_BUILD_STUB
	go test ./pkg/...
	rm ./ui/build/UI_BUILD_STUB


clean:
	rm -f bin/alice-lg-linux-amd64
	rm -f bin/alice-lg-osx-amd64
	rm -rf $(DIST)
	rm ./ui/build/UI_BUILD_STUB


.PHONY: backend ui clean

