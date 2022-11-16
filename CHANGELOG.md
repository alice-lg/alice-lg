
# Changelog

## 6.1.0

 * Added memory pools for deduplicating route information.
   This drastically reduces the memory consumption.

 * The BGP info modal can now be dismissed by pressing `esc`.

 * Bugfixes: 
   - Fixed parsing and handling of ext community filters.


## 6.0.0 (2022-11-10)
 
 * Pure functional react UI!

   Frontend is now using `create-react-app` for scripts and
   contexts instead of redux.

   **Theme compatibility**

   - Stylesheets are compatible
   - Content API is compatible
   - API now provides `Alice.onLayoutReady((page) => ... )`
     callback. This should be used to install additional
     DOM event listeners for extensions.

     So, if you want to inject additional dom nodes into
     the UI and used something like:
     
     `document.addEventListener("DOMContentLoaded", function() { ... }`

     you now need to use the `Alice.onLayoutReady(function(main) { ... })`
     callback.


## 5.1.1 (2022-06-21)
  
 * Improved search query validation.

 * Fixed http status response when validation fails.
   Was Internal Server Error (500), now: Bad Request (400).

 * Memory-Store is now using sync.Map to avoid timeouts
   due to aggressive locking.

## 5.1.0 (2022-06-02)

 * **BREAKING CHANGE** The spelling of "neighbors" is now harmonized.
   Please update your config and replace e.g. `neighbour.asn` 
   with `neighbor.asn` (in case of java script errors).
   This also applies to the API.

   In the config `neighbors_store_refresh_interval` needs to be updated.

 * Parallel route / neighbor store refreshs: Route servers are not
   longer queried sequentially. A jitter is applied to not hit all
   servers exactly at once.

 * Parallelism can be tuned through the config parameters:
    [server]

    routes_store_refresh_parallelism = 5
    neighbors_store_refresh_parallelism = 10000

   A value of 1 is a sequential refresh.

 * Postgres store backend: Not keeping routes and neighbors in
   memory might reduce the memory footprint.

 * Support for alternative pipe in `multi_table` birdwatcher
   configurations.

 * Reduced memory footprint by precomputing route details


## 5.0.1 (2021-11-01)

* Fixed parsing extended communities in openbgpd source causing a crash.

## 5.0.0 (2021-10-09)

* OpenBGPD support! Thanks to the Route Server Support Foundation
  for sponsoring this feature!

* Backend cleanup and restructured go codebase.
  This should improve a bit working with containers.

* Fixed links to the IRR Explorer.

## 4.3.0 (2021-04-15)

* Added configurable main table

## 4.2.0 (2020-07-29)

* Added GoBGP processing_timeout source config option

## 4.1.0 (2019-12-23)

* Added related neighbors feature

## 4.0.2, 4.0.3 (2019-09-09)

* Fixed issue with multitable bird: `getMasterPipeName` returned incorrect
  pipe.

* Fixed state check in multitable bird source with bird2.

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
