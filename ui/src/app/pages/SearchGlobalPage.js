
import { useConfig 
       , RoutesTableConfigProvider
       }
  from 'app/context/config';
import { usePageQuery }
  from 'app/context/pagination';
import { useFiltersQuery }
  from 'app/context/filters';
import { useSearchQuery
       , RoutesSearchProvider }
  from 'app/context/search';
import { RouteDetailsProvider
       , useRoutesLoading
       }
  from 'app/context/routes';
import { RoutesFiltersProvider }
  from 'app/context/filters';

import PageHeader
  from 'app/components/page/Header';
import SearchGlobalInput
  from 'app/components/search/SearchGlobalInput';
import WaitingCard
  from 'app/components/spinners/WaitingCard';
import FiltersEditor
  from 'app/components/filters/FiltersEditor';
import Routes 
  from 'app/components/routes/Routes';
import { CacheStatus }
  from 'app/components/status/Status';

const SearchStatus = () => {
  return (
    <table className="routeserver-status">
      <tbody>
        <CacheStatus />
      </tbody>
    </table>
  );  
}

const SearchGlobalContent = () => {
  const isLoading = useRoutesLoading();
  const search = useSearchQuery();
  const hasQuery = search !== "";

  return (
    <div className="lookup-container">
      <PageHeader></PageHeader>
      <div className="row details-main">
        <div className="col-main col-lg-9 col-md-12">

        <SearchGlobalInput />

        { hasQuery && <Routes /> }

        </div>
        { hasQuery &&
          <div className="col-lg-3 col-md-12 col-aside-details">
            <div className="card">
              <SearchStatus />
            </div>
            <WaitingCard isLoading={isLoading} />
            <FiltersEditor />
          </div> }
      </div>
    </div>
  );
}


/**
 * Global search is similar to Routes, however
 * routes are loaded by the SearchGlobalProvider
 */
const SearchGlobalPage = () => {
  const config = useConfig();

  const page = usePageQuery();
  const search = useSearchQuery();
  const [filters] = useFiltersQuery();

  return (
    <RoutesTableConfigProvider
      columns={config.lookup_columns}
      columnsOrder={config.lookup_columns_order}>

    <RoutesSearchProvider
      filters={filters}
      query={search}
      pageFiltered={page.filtered}
      pageReceived={page.received}>

      <RoutesFiltersProvider>
      <RouteDetailsProvider>
        <SearchGlobalContent />
      </RouteDetailsProvider>
      </RoutesFiltersProvider>

    </RoutesSearchProvider>
    </RoutesTableConfigProvider>
  );
}

export default SearchGlobalPage;
