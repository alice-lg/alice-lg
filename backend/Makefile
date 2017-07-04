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
	go get .

osx: deps
	GOARCH=$(ARCH) GOOS=darwin go build $(LDFLAGS) -o $(PROG)-osx-$(ARCH)

linux: deps
	GOARCH=$(ARCH) GOOS=linux go build $(LDFLAGS) -o $(PROG)-linux-$(ARCH)


bundle:
	rice embed-go


dev: clean all

clean:
	rm -f rice-box.go
	rm -f $(PROG)-osx-$(ARCH)
	rm -f $(PROG)-linux-$(ARCH)

coverage:
	go test -coverprofile=coverage.out
	go tool cover -func=coverage.out
