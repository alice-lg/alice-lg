
import { useQueryParams }
  from 'app/context/query';
import { useFiltersQuery }
  from 'app/context/filters';
import { RoutesSearchProvider }
  from 'app/context/search';
import { RouteDetailsProvider }
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

const SearchStatus = () => {
  return (
    <>Implement status</>
  );  
}

const SearchGlobalContent = () => {
  const isLoading = false;

  return (
    <div className="lookup-container">
      <PageHeader></PageHeader>
      <div className="row details-main">
        <div className="col-main col-lg-9 col-md-12">

        <SearchGlobalInput />

        <Routes />

        </div>
        <div className="col-lg-3 col-md-12 col-aside-details">
          <div className="card">
            <SearchStatus />
          </div>
          <WaitingCard isLoading={isLoading} />
    {/* <FiltersEditor /> */}
        </div>
      </div>
    </div>
  );
}


/**
 * Global search is similar to Routes, however
 * routes are loaded by the SearchGlobalProvider
 */
const SearchGlobalPage = () => {
  const query = useQueryParams({
    q: "",
    pf: "0",
    pr: "0",
  });
  const [filters] = useFiltersQuery();
  return (
    <RoutesSearchProvider
      filters={filters}
      query={query.q}
      pageFiltered={query.pf}
      pageReceived={query.pr}>

      <RoutesFiltersProvider>
      <RouteDetailsProvider>
        <SearchGlobalContent />
      </RouteDetailsProvider>
      </RoutesFiltersProvider>

    </RoutesSearchProvider>
  );
}

export default SearchGlobalPage;
