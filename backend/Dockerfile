
FROM golang:1.10

# Add project (for prefetching dependencies)
ADD . /go/src/github.com/alice-lg/alice-lg/backend

RUN cd /go/src/github.com/alice-lg/alice-lg/backend && go get -v .

RUN go get github.com/GeertJohan/go.rice/rice
RUN go install github.com/GeertJohan/go.rice/rice

WORKDIR /go/src/github.com/alice-lg/alice-lg
VOLUME ["/go/src/github.com/alice-lg/alice-lg"]

