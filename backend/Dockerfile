
FROM golang:1.16

# Add project (for prefetching dependencies)
ADD . /go/src/github.com/alice-lg/alice-lg/backend

RUN cd /go/src/github.com/alice-lg/alice-lg/backend && go get -v .

WORKDIR /go/src/github.com/alice-lg/alice-lg
VOLUME ["/go/src/github.com/alice-lg/alice-lg"]

