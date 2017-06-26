
# Alice LG - Your friendly looking glass

Alice is a frontend to the API exposed by 
services implementing Barry O'Donovan's
[birds-eye API design](https://github.com/inex/birds-eye-design/) to
[the BIRD routing daemon](http://bird.network.cz/):

 * INEX Birdseye API (https://github.com/inex/birdseye)
 * Birdwatcher (https://github.com/ecix/birdwatcher)


The project was started at the
[RIPE IXP Tools Hackathon](https://atlas.ripe.net/hackathon/ixp-tools/) 
just prior to [RIPE73](https://ripe73.ripe.net/) in Madrid, Spain.


## Building Alice from scratch

Alice requires a working (and configured) `golang` installation
for the backend.
The frontend requires `yarn` and `gulp` for building.


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



