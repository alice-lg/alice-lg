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
FILES=$(shell find . -depth 1 ! -name "*_test.go" -name "*.go")

all: $(TARGET)
	@echo "Built $(VERSION) @ $(TARGET)"

deps:
	GO111MODULE=on go get -v .

osx-dev: deps
	GO111MODULE=on go run $(FILES)

osx: deps
	GO111MODULE=on GOARCH=$(ARCH) GOOS=darwin go build $(LDFLAGS) -o $(PROG)-osx-$(ARCH)

linux: deps
	GO111MODULE=on GOARCH=$(ARCH) GOOS=linux go build $(LDFLAGS) -o $(PROG)-linux-$(ARCH)

bundle:
	rice embed-go

test:
	GO111MODULE=on go test -v
	cd api/ && GO111MODULE=on go test -v
	cd caches/ && GO111MODULE=on go test -v
	cd sources/birdwatcher && GO111MODULE=on go test -v


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
