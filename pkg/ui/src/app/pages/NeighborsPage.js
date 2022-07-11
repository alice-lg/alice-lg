
import { useState 
       }
  from 'react';

import { useSelectedRouteServer }
  from 'app/components/routeservers/Provider';
import RouteServerStatusProvider
  from 'app/components/routeservers/StatusProvider';

import Status
  from 'app/components/routeservers/Status';

import NeighborsProvider
  from 'app/components/neighbors/Provider';
import Neighbors
  from 'app/components/neighbors/Neighbors';
import QuickLinks
  from 'app/components/neighbors/QuickLinks';

import PageHeader
  from 'app/components/page/Header';
import SearchInput
  from 'app/components/search/Input';

/** 
 * The NeighborsPage renders a list of all peers on
 * the route server.
 *
 * A search field for quick filtering is provided
 */
const NeighborsPage = () => {
  const [filterValue, setFilterValue] = useState("");
  const routeServer = useSelectedRouteServer();

  if (!routeServer) { return null; } // nothing to do here

  return (
    <NeighborsProvider routeServerId={routeServer.id}>
    <RouteServerStatusProvider routeServerId={routeServer.id}>
    <div className="routeservers-page">
      <PageHeader>
       <span className="status-name">{routeServer.name}</span>
      </PageHeader>

      <div className="row details-main">
        <div className="col-main col-lg-9 col-md-12">
          <div className="card">
            <SearchInput
              value={filterValue}
              placeholder="Filter by Neighbor, ASN or Description"
              onChange={(e) => setFilterValue(e.target.value)}
            />
          </div>
          <QuickLinks />
          <Neighbors filter={filterValue} />
        </div>
        <div className="col-lg-3 col-md-12 col-aside-details">
          <div className="card">
            <Status />
          </div>
        </div>
      </div>
    </div>
    </RouteServerStatusProvider>
    </NeighborsProvider>
  );
}

export default NeighborsPage;
