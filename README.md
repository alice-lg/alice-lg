# Alice-LG - Your friendly looking glass
__"No, no! The adventures first, explanations take such a dreadful time."__  
_Lewis Carroll, Alice's Adventures in Wonderland & Through the Looking-Glass_

Take a look at Alice-LG production examples at:
- https://lg.de-cix.net/
- https://lg.ams-ix.net
- https://lg.bcix.de/
- https://lg.megaport.com/
- https://lg.netnod.se/
- https://alice-rs.linx.net/
- https://lg.ix.br/

And checkout the API at:
- https://lg.de-cix.net/api/v1/config
- https://lg.de-cix.net/api/v1/routeservers
- https://lg.de-cix.net/api/v1/routeservers/rs1_fra_ipv4/status
- https://lg.de-cix.net/api/v1/routeservers/rs1_fra_ipv4/neighbors
- https://lg.de-cix.net/api/v1/routeservers/rs1_fra_ipv4/neighbors/R194_106/routes
- https://lg.de-cix.net/api/v1/lookup/prefix?q=217.115.0.0


## Breaking Changes

### 6.0.0

With the new functional react UI, the `DOMContentLoaded` event can no
longer be used for injecting additional content. Please use
`Alice.onLayoutReady(function(main) { ... });` instead.

### 5.1.0

The spelling of "neighbors" is now harmonized. Please update your config and
replace e.g. neighbour.asn with neighbor.asn (in case of java script errors).

## Explanations
Alice-LG is a BGP looking glass which gets its data from external APIs.

Currently Alice-LG supports the following APIs:
- [birdwatcher API](https://github.com/alice-lg/birdwatcher) for [BIRD](http://bird.network.cz/)
- [GoBGP](https://osrg.github.io/gobgp/)

### Birdwatcher
Normally you would first install the [birdwatcher API](https://github.com/alice-lg/birdwatcher) directly on the machine(s) where you run [BIRD](http://bird.network.cz/) on
and then install Alice-LG on a seperate public facing server and point her to the afore mentioned [birdwatcher API](https://github.com/alice-lg/birdwatcher).

This project was a direct result of the [RIPE IXP Tools Hackathon](https://atlas.ripe.net/hackathon/ixp-tools/)
just prior to [RIPE73](https://ripe73.ripe.net/) in Madrid, Spain.

Major thanks to Barry O'Donovan who built the original [INEX Bird's Eye](https://github.com/inex/birdseye) BIRD API of which Alice-LG is a spinnoff

### GoBGP
Alice-LG supports direct integration with GoBGP instances using gRPC.  See the configuration section for more detail.

### OpenBGPD

Alice-LG supports OpenBGP via [`bgplgd`](https://github.com/cjeker/bgplgd) 
and [`openbgpd-state-server`](https://github.com/alice-lg/openbgpd-state-server).

## Building Alice-LG from scratch
__These examples include setting up your Go environment, if you already have set that up then you can obviously skip that__

### CentOS:
First add the following lines at the end of your `~/.bash_profile`:
```bash
GOPATH=$HOME/go
export GOPATH
PATH=$PATH:$GOPATH/bin
export PATH
```
Now run:
```bash
source ~/.bash_profile

# Install frontend build dependencies
sudo yum install golang npm
sudo npm install --global yarn

mkdir -p ~/go/bin ~/go/pkg ~/go/src/github.com/alice-lg/

cd ~/go/src/github.com/alice-lg
git clone https://github.com/alice-lg/alice-lg.git

cd alice-lg
make
```
Your Alice-LG source will now be located at `~/go/src/github.com/alice-lg/alice-lg` and your alice-LG executable should be at `~/go/src/github.com/alice-lg/alice-lg/bin/alice-lg-linux-amd64`

## Configuration

An example configuration can be found at
[etc/alice-lg/alice.example.conf](https://github.com/alice-lg/alice-lg/blob/master/etc/alice-lg/alice.example.conf).

You can copy it to any of the following locations:

    etc/alice-lg/alice.conf        # local
    etc/alice-lg/alice.local.conf  # local
    /etc/alice-lg/alice.conf       # global


You will have to edit the configuration file as you need to point Alice-LG to the correct backend source.  Multiple sources can be configured.

[Birdwatcher](https://github.com/alice-lg/birdwatcher):
```ini
[source.rs1-example-v4]
name = rs1.example.com (IPv4)
[source.rs1-example-v4.birdwatcher]
api = http://rs1.example.com:29184/
neighbors_refresh_timeout = 2
# show_last_reboot = true
# timezone = UTC
# type = single_table / multi_table
type = multi_table
# not needed for single_table
peer_table_prefix = T
pipe_protocol_prefix = M

[source.rs1-example-v6]
name = rs1.example.com (IPv6)
[source.rs1-example-v6.birdwatcher]
api = http://rs1.example.com:29186/
```

[GoBGP](https://osrg.github.io/gobgp/):
```ini
[source.rs2-example]
name = rs2.example.com
group = AMS

[source.rs2-example.gobgp]
# Host is the IP (or DNS name) and port for the remote GoBGP daemon
host = rs2.example.com:50051
# ProcessingTimeout is a timeout in seconds configured per gRPC call to a given GoBGP daemon
processing_timeout = 300
```

[OpenBGPD](https://www.openbgpd.org/) via `openbgpd-state-server`:
```ini
[source.rs-example]
name = rs-example.openbgpd-state-server

[source.rs-example.openbgpd-state-server]
api = http://rs23.example.net:29111/api

# Optional response cache time in seconds
# Default: disabled (0)
cache_ttl = 100
```

[OpenBGPD](https://www.openbgpd.org/) via `bgplgd`:
```ini
[source.rs-example]
name = rs-example.openbgpd-bgplgd

[source.rs-example.openbgpd-bgplgd]
api = http://rs23.example.net/bgplgd

# Optional response cache time in seconds
# Default: disabled (0)
cache_ttl = 100
```

## Running

Launch the server by running

    ./bin/alice-lg-linux-amd64


## Customization

Alice now supports custom themes!
In your alice.conf, you now can specify a theme by setting:

    [theme]
    path = /path/to/my/alice-theme

with the optional parameter (the "mountpoint" of the theme)
    url_base = /theme


You can put assets (images, fonts, javscript, css) in
this folder.

Stylesheets and Javascripts are automatically included in
the client's html and are served from the backend.

Alice provides early stages of an extension API, which is for now
only used to modify the content of the welcome screen,
by providing a javascript in your theme containing:

```javascript
Alice.updateContent({
    welcome: {
        title: "My Awesome Looking Glass",
        tagline: "powered by Alice"
    }
});

```

A callback for running custom javascript after the base application
was initialized can be installed using:

```javascript
Alice.onLayoutReady(function(page) {
    // page is the layout HTML root element
});
```

For an example check out: https://github.com/alice-lg/alice-theme-example

## Hacking

The client is a Single Page React Application.
All sources are available in `ui/`.

`Yarn` is required for building the UI.

Create a fresh UI build with
```bash
cd ui/
make
```

This will install all dependencies with `yarn install` and run `yarn build`.

As this is a `create-react-app` application, react-scripts are present
and you can just run a development server using `yarn start`.

All this available as a containerized environment:

Running `docker-compose up` in the `./dev` will build and start the
backend and run a webpack dev server for the UI.

The UI is then available on http://localhost:3000/ and on http://localhost:7340/
the backend will serve the API.


## Sponsors

The development of Alice is now sponsored by
<p align="center">
    <a href="https://www.de-cix.net" target="_blank">
        <img src="doc/images/DE-CIX_Logo_2016_small.png?raw=true" alt="DE-CIX Logo" title="DE-CIX" />
    </a>
</p>

Many thanks go out to [ECIX](https://www.ecix.net), where this project originated and was backed over the last two years.

Support for **OpenBGPD** was sponsored by the [Route Server Support Foundation](https://www.rssf.nl/).

