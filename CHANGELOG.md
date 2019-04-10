
# Changelog

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
