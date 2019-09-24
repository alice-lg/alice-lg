
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
FROM golang:1.12 AS backend

# Install dependencies
WORKDIR /src/alice-lg
ADD go.mod .
ADD go.sum .
RUN go mod download
RUN go install github.com/GeertJohan/go.rice/rice

# Add client
COPY --from=frontend /src/alice-lg/client/build client/build

# Build backend
WORKDIR /src/alice-lg/backend
ADD backend .
RUN rice embed-go

RUN go build -o alice-lg-linux-amd64 -ldflags="-X main.version=4.0.3"

EXPOSE 7340:7340
CMD ["/src/alice-lg/backend/alice-lg-linux-amd64"]
