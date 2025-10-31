
#
# Alice - The friendly BGP looking glass
#

# Build frontend first
FROM node:latest AS ui

# Install dependencies 
WORKDIR /src/alice-lg/ui
ADD ui/package.json .
ADD ui/yarn.lock .

RUN yarn install

# Add frontend
ADD ui/ .

# Build frontend
RUN yarn build 

# Build the backend
FROM golang:1.24 AS backend

# Install dependencies
WORKDIR /src/alice-lg
ADD go.mod .
ADD go.sum .
RUN go mod download

ADD . .

# Add client
COPY --from=ui /src/alice-lg/ui/build ui/build

WORKDIR /src/alice-lg/cmd/alice-lg
RUN make alpine

FROM alpine:latest

RUN apk add -U tzdata

COPY --from=backend /src/alice-lg/cmd/alice-lg/alice-lg-linux-amd64 /usr/bin/alice-lg
RUN ls -lsha /usr/bin/alice-lg

EXPOSE 7340:7340
CMD ["/usr/bin/alice-lg"]
