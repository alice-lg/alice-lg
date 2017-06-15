# Alice-LG - Your friendly looking glass

Alice-LG is a BGP looking glass which gets its data from external APIs.

Take a look at an Alice-LG production example:
- https://lg.ecix.net/
- https://lg.ecix.net/api/config
- https://lg.ecix.net/api/routeservers
- https://lg.ecix.net/api/routeservers/0
- https://lg.ecix.net/api/routeservers/0/neighbours
- https://lg.ecix.net/api/routeservers/0/neighbours/ID109_AS31078/routes
- https://lg.ecix.net/api/routeservers/0/lookup/prefix?q=217.115.15.0/20

Currently Alice-LG supports the following APIs:
- [birdwatcher API](https://github.com/ecix/birdwatcher) for [BIRD](http://bird.network.cz/)

This project was a direct result of the [RIPE IXP Tools Hackathon](https://atlas.ripe.net/hackathon/ixp-tools/) 
just prior to [RIPE73](https://ripe73.ripe.net/) in Madrid, Spain.

Major thanks to Barry O'Donovan who built the original [INEX Bird's Eye](https://github.com/inex/birdseye) BIRD API of which Alice_LG is a spinnoff

## Alice RPMs

## Building Alice from scratch
### TLDR CentOS 7:

mkdir -p ~/go/bin ~/go/pkg ~/go/src

### TLDR Ubuntu:
apt-get install golang
mkdir -p ~/go/bin ~/go/pkg ~/go/src


### Installing and configuring golang
Alice requires a working (and configured) `golang` installation
for the backend. If you are already set up for go then just skip ahead!

A full guide on setting up golang can be found at: https://golang.org/doc/install




The frontend requires `npm` for building.


Clone this repository in your go workspace and type
`make`

This will download all required *go* and *js* dependencies
and will start building alice.


## Installation

For systemwide deployment it is advised to add the contents
of the local `etc/` to your system's `/etc`
directory.



## Configuration

An example configuration can be found under
[etc/alicelg/alice.example.conf](https://github.com/ecix/alice/blob/master/etc/alicelg/alice.example.conf).

You can copy it to any of the following locations:

    etc/alicelg/alice.conf # local
    etc/alicelg/alice.local.conf # local as well
    /etc/alicelg/alice.conf # global


You will have to at least edit it to add bird API servers:

    [source.0]
    name = rs1.example.com (IPv4)
    [source.0.birdwatcher]
    api = http://rs1.example.com:29184/
    # Optional:
    show_last_reboot = true
    timezone = UTC

    [source.1]
    name = rs1.example.com (IPv6)
    [source.1.birdwatcher]
    api = http://rs1.example.com:29186/


## Running

Launch the server by running

    ./bin/alice-lg-linux-amd64


## Deployment

We added a `Makefile` for packaging Alice as an RPM using [fpm](https://github.com/jordansissel/fpm).

If you have all tools available locally, you can just type:

    make rpm

If you want to build the package on a remote machine, just use

    make remote_rpm BUILD_SERVER=my-rpm-building-server.example.com

which will copy the dist to the remote server and executes fpm via ssh.

You can specify which system integration to use:
Set the `SYSTEM_INIT` variable to `upstart` or `systemd` (default)
prior to building the RPM.

    make remote_rpm BUILD_SERVER=rpmbuild.example.com SYSTEM_INIT=upstart



## Hacking

The client is a Single Page React Application.
All sources are available in `client/`. 

Install build tools as needed:

    npm install -g gulp-cli


Create a fresh UI build with

    cd client/
    make client

This will install all dependencies and run `gulp`.

While working on the UI you might want to use `make watch`,
which will keep the `gulp watch` task up and running.



