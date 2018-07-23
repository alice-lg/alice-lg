#
# Alice LG Backend 
# ----------------
#

PROG=alice-lg
ARCH=amd64

APP_VERSION=$(shell cat ../VERSION)
VERSION=$(APP_VERSION)_$(shell git rev-parse --short HEAD)

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

deps:
	go get -v .

osx: deps
	GOARCH=$(ARCH) GOOS=darwin go build $(LDFLAGS) -o $(PROG)-osx-$(ARCH)

linux: deps
	GOARCH=$(ARCH) GOOS=linux go build $(LDFLAGS) -o $(PROG)-linux-$(ARCH)

bundle:
	rice embed-go

test:
	go test -v
	cd caches/ && go test -v
	cd sources/birdwatcher && go test -v


dev: clean all

prod: clean bundle $(TARGET)
	@echo "Build $(VERSION) [production] @ $(TARGET)"


clean:
	rm -f rice-box.go
	rm -f $(PROG)-osx-$(ARCH)
	rm -f $(PROG)-linux-$(ARCH)

coverage:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
