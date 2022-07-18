
import { useQuery }
  from 'app/components/query';

import { useSelectedRouteServer }
  from 'app/components/routeservers/Provider';

import Status
  from 'app/components/routeservers/Status';

import Neighbors
  from 'app/components/neighbors/Neighbors';
import QuickLinks
  from 'app/components/neighbors/QuickLinks';

import PageHeader
  from 'app/components/page/Header';
import SearchQueryInput
  from 'app/components/query/SearchQueryInput';



/** 
 * The NeighborsPage renders a list of all peers on
 * the route server.
 *
 * A search field for quick filtering is provided
 */
const NeighborsPage = () => {
  const routeServer = useSelectedRouteServer();
  if (!routeServer) { return null; } // nothing to do here

  return (
    <div className="routeservers-page">
      <PageHeader>
       <span className="status-name">{routeServer.name}</span>
      </PageHeader>

      <div className="row details-main">
        <div className="col-main col-lg-9 col-md-12">
          <div className="card">
        
            <SearchQueryInput
              placeholder="Filter by Neighbor, ASN or Description" />
          </div>
          <QuickLinks />
          <Neighbors />
        </div>
        <div className="col-lg-3 col-md-12 col-aside-details">
          <div className="card">
            <Status />
          </div>
        </div>
      </div>
    </div>
  );
}

export default NeighborsPage;
