
import { useState 
       }
  from 'react';

import { useSelectedRouteServer }
  from 'app/components/routeservers/Provider';

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
    { /* <QuickLinks /> */ }
    QUICKLINKS

    { /* <Protocols protocol="bgp" routeserverId={this.props.params.routeserverId} /> */ }
          PROTOCOLS
        </div>
        <div className="col-lg-3 col-md-12 col-aside-details">
          <div className="card">
            STATUS
            {/*
            <Status routeserverId={this.props.params.routeserverId}
                    cacheStatus={this.props.cacheStatus} />
              */}
          </div>
        </div>
      </div>
    </div>
  );
}

export default NeighborsPage;
