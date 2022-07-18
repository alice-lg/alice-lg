
import { Link
       , useParams
       }
  from 'react-router-dom';

import { humanizedJoin }
  from 'app/utils/text'

import { intersect
       , resolve
       }
  from 'app/utils/lists'

import { useQuery }
  from 'app/components/query';
import { useConfig }
  from 'app/components/config/Provider';
import { useSelectedRouteServer }
  from 'app/components/routeservers/Provider';
import { NeighborProvider
       , useNeighbor
       }
  from 'app/components/neighbors/NeighborProvider';

import PageHeader
  from 'app/components/page/Header';
import SearchQueryInput
  from 'app/components/query/SearchQueryInput';


const FILTERABLE_COLUMNS = [
  "gateway", "network"
];


const filterableColumnsText = (columns, order) => {
  const filterable = resolve(columns, intersect(order, FILTERABLE_COLUMNS));
  return humanizedJoin(filterable, "or");
}


const RoutesHeader = () => {
  const routeServer = useSelectedRouteServer();
  const neighbor = useNeighbor();

  if (!routeServer || !neighbor) {
    return null;
  }

  return (
    <>
      <Link to={`/routeservers/${routeServer.id}`}>
        {routeServer.name}
      </Link>
      <span className="spacer">&raquo;</span>
        {neighbor.description}
    </>
  );
}


/**
 * RoutesPage renders the page with all routes for a neighbor
 * on a route server
 */
const RoutesPage = () => {
  const { neighborId } = useParams();
  const [query, setQuery] = useQuery({q: ""});
  const config = useConfig();

  const routesColumns = config.routes_columns;
  const routesColumnsOrder = config.routes_columns_order;

  let pageClass = "routeservers-page";
  /*
   * TODO: find better solution.
  if (this.props.localRelatedPeers.length > 1) {
    pageClass += " has-related-peers";
  }
  */
  
  const filterPlaceholder = "Filter by " +
    filterableColumnsText(routesColumns, routesColumnsOrder);


  return (
    <NeighborProvider neighborId={neighborId}>
    <div className={pageClass}>
      <PageHeader>
        <RoutesHeader />
      </PageHeader>
      <div className="row details-main">
        <div className="col-main col-lg-9 col-md-12">

          <div className="card">
            [Related Peers Tab]
            <SearchQueryInput
              placeholder={filterPlaceholder} />
            { /*
            <RelatedPeersTabs
              peers={this.props.localRelatedPeers}
              protocolId={this.props.params.protocolId}
              routeserverId={this.props.params.routeserverId} />
             */}
          </div>

        </div>
      </div>
    </div>
    </NeighborProvider>
  );

}

/*
 
      <BgpAttributesModal />
      <div className="row details-main">
        <div className="col-main col-lg-9 col-md-12">

          <div className="card">
            <RelatedPeersTabs
              peers={this.props.localRelatedPeers}
              protocolId={this.props.params.protocolId}
              routeserverId={this.props.params.routeserverId} />
            <SearchInput
              value={this.props.filterValue}
              placeholder={filterPlaceholder}
              onChange={(e) => this.setFilter(e.target.value)}  />
          </div>

          <QuickLinks routes={this.props.routes} />

          <RoutesViewEmpty routes={this.props.routes}
                           hasQuery={!!this.props.filterValue}
                           loadNotExported={this.props.loadNotExported} />
          <RoutesView
              type={ROUTES_FILTERED}
              routeserverId={this.props.params.routeserverId}
              protocolId={this.props.params.protocolId} />

          {this.props.receivedLoading && <RoutesLoadingIndicator />}

          <RoutesView
              type={ROUTES_RECEIVED}
              routeserverId={this.props.params.routeserverId}
              protocolId={this.props.params.protocolId} />

          {this.props.notExportedLoading && <RoutesLoadingIndicator />}

          <RoutesView
              type={ROUTES_NOT_EXPORTED}
              routeserverId={this.props.params.routeserverId}
              protocolId={this.props.params.protocolId} />


        </div>
        <div className="col-lg-3 col-md-12 col-aside-details">
          <div className="card">
            <Status routeserverId={this.props.params.routeserverId}
                    cacheStatus={cacheStatus} />
          </div>
          <FiltersEditor makeLinkProps={makeLinkProps}
                         linkProps={this.props.linkProps}
                         filtersApplied={this.props.filtersApplied}
                         filtersAvailable={this.props.filtersAvailable} />
          <RelatedPeersCard
            neighbors={this.props.allRelatedPeers}
            rsId={this.props.params.routeserverId} 
            protocolId={this.props.params.protocolId} />
        </div>
      </div>
    </div>
*/

export default RoutesPage;
