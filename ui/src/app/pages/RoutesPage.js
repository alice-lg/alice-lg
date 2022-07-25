
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

import { useConfig 
       , RoutesTableConfigProvider
       }
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
       , RouteDetailsProvider
       , useRoutesLoading
       , useNotExportedEnabled
       }
  from 'app/context/routes';
import { useSearchQuery }
  from 'app/context/search';
import { usePageQuery }
  from 'app/context/pagination';
import { RoutesFiltersProvider
       , useFiltersQuery
       }
  from 'app/context/filters';

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
import WaitingCard
  from 'app/components/spinners/WaitingCard';
import FiltersEditor
  from 'app/components/filters/FiltersEditor';


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
  const isLoading = useRoutesLoading();

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
          <WaitingCard isLoading={isLoading} />
          <RelatedPeersCard />
          <FiltersEditor />
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
  const config = useConfig();

  const { neighborId, routeServerId } = useParams();

  const page = usePageQuery();
  const search = useSearchQuery();
  const [filters] = useFiltersQuery();

  const notExportedEnabled = useNotExportedEnabled();

  // Setup context and render content
  return (
    <NeighborProvider neighborId={neighborId}>
    <RelatedNeighborsProvider>
    <RouteDetailsProvider>
    <RoutesTableConfigProvider
      columns={config.routes_columns}
      columnsOrder={config.routes_columns_order}>

    <RoutesNotExportedProvider
      routeServerId={routeServerId}
      neighborId={neighborId}
      page={page.notExported}
      query={search}
      filters={filters}
      enabled={notExportedEnabled}>
    <RoutesFilteredProvider
      routeServerId={routeServerId}
      neighborId={neighborId}
      query={search}
      filters={filters}
      page={page.filtered}>
    <RoutesReceivedProvider 
      routeServerId={routeServerId}
      neighborId={neighborId}
      query={search}
      filters={filters}
      page={page.received}>{/* innermost used for api status */}

      <RoutesFiltersProvider>
        <RoutesPageContent />
      </RoutesFiltersProvider>

    </RoutesReceivedProvider>
    </RoutesFilteredProvider>
    </RoutesNotExportedProvider>

    </RoutesTableConfigProvider>
    </RouteDetailsProvider>
    </RelatedNeighborsProvider>
    </NeighborProvider>
  );
}

export default RoutesPage;
