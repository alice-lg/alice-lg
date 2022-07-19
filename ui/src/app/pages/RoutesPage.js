
import { useMemo }
  from 'react';

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

import { useQueryParams }
  from 'app/context/query';
import { useConfig }
  from 'app/context/config';
import { useRouteServer }
  from 'app/context/route-servers';
import { NeighborProvider
       , RelatedNeighborsProvider
       , useNeighbor
       , useLocalRelatedPeers
       }
  from 'app/context/neighbors';
import { RoutesReceivedProvider
       , RoutesFilteredProvider
       , RoutesNotExportedProvider
       }
  from 'app/context/routes';

import PageHeader
  from 'app/components/page/Header';
import Status
  from 'app/components/status/Status';
import LocalRelatedPeersTabs
  from 'app/components/neighbors/LocalRelatedPeersTabs';
import RelatedPeersCard
  from 'app/components/neighbors/RelatedPeersCard';
import Routes 
  from 'app/components/routes/Routes';
import SearchQueryInput
  from 'app/components/search/SearchQueryInput';


const FILTERABLE_COLUMNS = [
  "gateway", "network"
];


const filterableColumnsText = (columns, order) => {
  const filterable = resolve(columns, intersect(order, FILTERABLE_COLUMNS));
  return humanizedJoin(filterable, "or");
}


const RoutesPageHeader = () => {
  const routeServer = useRouteServer();
  const neighbor = useNeighbor();
  if (!routeServer || !neighbor) {
    return null;
  }
  return (
    <PageHeader>
      <Link to={`/routeservers/${routeServer.id}`}>
        {routeServer.name}
      </Link>
      <span className="spacer">&raquo;</span>
        {neighbor.description}
    </PageHeader>
  );
}

const RoutesPageSearch = () => {
  const config = useConfig();
  const filterPlaceholder = useMemo(() => (
    "Filter by " + filterableColumnsText(
      config.routes_columns, config.routes_columns_order)
  ), [config]);

  return (
    <SearchQueryInput placeholder={filterPlaceholder} />
  );
}

const RoutesPageContent = () => {
  const localRelatedPeers = useLocalRelatedPeers();

  let pageClass = "routeservers-page";
  if (localRelatedPeers.length > 1) {
    pageClass += " has-related-peers";
  }

  return (
    <div className={pageClass}>
      <RoutesPageHeader />
      <div className="row details-main">
        <div className="col-main col-lg-9 col-md-12">

          <div className="card">
            <LocalRelatedPeersTabs />
            <RoutesPageSearch />
          </div>

          <Routes />

        </div>
        <div className="col-lg-3 col-md-12 col-aside-details">
          <div className="card">
            <Status />
          </div>
          <RelatedPeersCard />
          { /* 
          <FiltersEditor makeLinkProps={makeLinkProps}
                         linkProps={this.props.linkProps}
                         filtersApplied={this.props.filtersApplied}
                         filtersAvailable={this.props.filtersAvailable} />
          <RelatedPeersCard
            neighbors={this.props.allRelatedPeers}
            rsId={this.props.params.routeserverId} 
            protocolId={this.props.params.protocolId} />
            */ }
        </div>
      </div>
    </div>
  );
}

/**
 * RoutesPage renders the page with all routes for a neighbor
 * on a route server
 */
const RoutesPage = () => {
  const { neighborId, routeServerId } = useParams();
  const query = useQueryParams({
    pr: 0,
    pf: 0,
    ne: 0,
    q: "",
  });
  const notExportedEnabled = query.ne === 1;

  return (
    <NeighborProvider neighborId={neighborId}>
    <RelatedNeighborsProvider>

    <RoutesNotExportedProvider
      routeServerId={routeServerId}
      neighborId={neighborId}
      page={query.pn}
      query={query.q}
      enabled={notExportedEnabled}>
    <RoutesFilteredProvider
      routeServerId={routeServerId}
      neighborId={neighborId}
      query={query.q}
      page={query.pf}>
    <RoutesReceivedProvider 
      routeServerId={routeServerId}
      neighborId={neighborId}
      query={query.q}
      page={query.pr}>{/* innermost used for api status */}

      <RoutesPageContent />

    </RoutesReceivedProvider>
    </RoutesFilteredProvider>
    </RoutesNotExportedProvider>

    </RelatedNeighborsProvider>
    </NeighborProvider>
  );
}

/*
 
      <BgpAttributesModal />
      <div className="row details-main">
        <div className="col-main col-lg-9 col-md-12">

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
      </div>
    </div>
*/

export default RoutesPage;
