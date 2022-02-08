
#
# Alice - The friendly BGP looking glass
#

# Build frontend first
FROM node:11 AS frontend

# Install dependencies 
WORKDIR /src/alice-lg/client
ADD client/package.json .
ADD client/yarn.lock .

RUN npm install -g gulp@4.0.0
RUN npm install -g gulp-cli
RUN yarn install

# Add frontend
WORKDIR /src/alice-lg/client
ADD client .

# Build frontend
RUN DISABLE_LOGGING=1 NODE_ENV=production /usr/local/bin/gulp

# Build the backend
FROM golang:1.16 AS backend

# Install dependencies
WORKDIR /src/alice-lg
ADD go.mod .
ADD go.sum .
RUN go mod download

ADD . .

# Add client
COPY --from=frontend /src/alice-lg/client/build client/build

WORKDIR /src/alice-lg/cmd/alice-lg
RUN make alpine

FROM alpine:latest

RUN apk add -U tzdata

COPY --from=backend /src/alice-lg/cmd/alice-lg/alice-lg-linux-amd64 /usr/bin/alice-lg
RUN ls -lsha /usr/bin/alice-lg

EXPOSE 7340:7340
CMD ["/usr/bin/alice-lg"]
