#
# Alice LG Backend 
# ----------------
#

PROG=alice-lg
ARCH=amd64

APP_VERSION=$(shell cat ../../VERSION)
VERSION=$(APP_VERSION)_$(shell git rev-parse --short HEAD)


# OS Detection
UNAME=$(shell uname)
ifeq ($(UNAME), Darwin)
  TARGET=osx
else
  TARGET=linux
endif


LDFLAGS=-ldflags="-X github.com/alice-lg/alice-lg/pkg/config.Version=$(APP_VERSION)"
LDFLAGS_STATIC=-ldflags="-X github.com/alice-lg/alice-lg/pkg/config.Version=$(APP_VERSION) -extldflags '-static'"

FILES=$(shell find . -depth 1 ! -name "*_test.go" -name "*.go")

all: deps $(TARGET)
	@echo "Built $(VERSION) @ $(TARGET)"

deps:
	GO111MODULE=on go get -v .

osx-dev:
	GO111MODULE=on go run $(FILES)

osx:
	GO111MODULE=on GOARCH=$(ARCH) GOOS=darwin go build $(LDFLAGS) -o $(PROG)-osx-$(ARCH)

linux:
	GO111MODULE=on GOARCH=$(ARCH) GOOS=linux go build $(LDFLAGS) -o $(PROG)-linux-$(ARCH)

alpine:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux go build \
		-a -installsuffix cgo \
		$(LDFLAGS) -o $(PROG)-linux-$(ARCH)

static:
	GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=$(ARCH) \
		go build $(CFLAGS) \
		-a $(LDFLAGS_STATIC) \
		-o $(PROG)-linux-$(ARCH)


test:
	GO111MODULE=on go test -v
	cd ../../pkg/api/ && GO111MODULE=on go test -v
	cd ../../pkg/caches/ && GO111MODULE=on go test -v
	cd ../../pkg/sources/birdwatcher && GO111MODULE=on go test -v


dev: clean all

prod: clean $(TARGET)
	@echo "Build $(VERSION) [production] @ $(TARGET)"


clean:
	rm -f $(PROG)-osx-$(ARCH)
	rm -f $(PROG)-linux-$(ARCH)

coverage:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
