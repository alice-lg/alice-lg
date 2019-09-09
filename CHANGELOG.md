
# Changelog

## 4.0.2 (2019-09-09)

* Fixed issue with multitable bird: `getMasterPipeName` returned incorrect
  pipe.

## 4.0.1 (2019-03-07)

* Enhance the neighbors store to perform uncached requests for peer status
  on every request. A timeout with fallback to cached data is applied in order
  too keep the response times low.
* Add caching to Neighbors()

## 4.0.0 (2019-02-22)

Breaking Changes: Birdwatcher 2.0

Support for birdwatcher route server API implementation version 2.0.0 and above.  
This new implementation of birdwatcher only provides the direct output of the
birdc comands and eliminates complex endpoints that fetch data from multiple
birdc responses. The aggregation of data, based on the particular route server
setup in use is now implemented in Alice-LG.
Therefore the birdwatcher source can be configured with a new config parameter
'type', which specifies a processing strategy for the ingested data which
corresponds to a particular layout of the routing daemon (BIRD) configuration
(e.g. single-table, multi-table or something even more custom). For developers
it is made easy to add new configuration types.

The neighbor summary has been removed, since much of it's data can be requested
from the new birdwatcher endpoints in alternative ways.

The config option from birdwatcher "PeerTablePrefix" and "PipeProtocolPrefix"
have been carried over to Alice-LG. These constants may be defined on a
per route server basis and are used to generate the request URLs for the
route server (birdwatcher) API in case of multi-table setup.

In addition this version contains the following bug-fixes and features:
* Fix a bug in Neighbors(), a peer that is down would cause a runtime error
* Fix the cache, it would still store entries even if disabled
* Fix a bug affecting the cache (subsequent modification of entries)
* Remove additional caches to avoid duplicate caching and save memory
* Save memory by periodically expiring entries with a housekeeping routine
* Change extended communities format to (string, string, string)

## 3.4.4 (2019-01-29)

* Loading indicators in frontend for received routes and filtered routes
* Consequent use of 'alice-lg' instead of 'alicelg'
* Use yarn to manage node dependencies
* Improved error handling on misconfigured or unavailable sources
* Fix issues related to pagination of results
* Make search for IPv6 prefixes with netmask work
* Add example for routeserver ids as strings to config and overall improvement

## 3.4.0 (2018-12-09)

* Removed baseUrl from frontend

* Introduced static routeserver ids


## 3.0.0 (2018-10-03)

### Breaking changes:

* The API endpoints is now include the API version,
  e.g. /api/v1/status, /api/v1/routeservers, ...

* The API is now consistently using 'neighbors' instead of 'neighbours'

* Reject reasons are now configured in BGP community
  notataion: 1234:65666:1 = My filter reason


## 2.3.0 (2018-09-10)

### New Features:

* Sortable columns in neighbors table
* Show prefix 'flags': Best Route and Blackhole
* BGP-Communities are now human readable
* Added links to related peers in routes view
* Added quick links to routes received, filtered and not-exported
* Added quick links to bgp sessions established and down
* Routes not exported can now be configured to be loaded on demand
* Routes can now be configured to be paginated
* Added support for birdwatcher neighbor summary capabilities
* Information about the cache-state was added
* Skin / Theme support

### Fixes:

* Performance improvements by eliminating copy operations
* Time information in the API is now normalized to UTC
* We improved the error handling a bit


## 2.2.6 (2018-01-31)

* Improved logging of missing birdwatcher modules
* Fixed bugs and improved documentation
