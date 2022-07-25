
import { useRef
       , useEffect
       }
  from 'react';
import { Link
       }
  from 'react-router-dom';

import { useConfig }
  from 'app/context/config';
import { useQuery
       , useQueryLocation
       , PARAM_LOAD_NOT_EXPORTED
       , PARAM_PAGE_RECEIVED
       , PARAM_PAGE_FILTERED
       , PARAM_PAGE_NOT_EXPORTED
       }
  from 'app/context/query';
import { ROUTES_RECEIVED
       , ROUTES_FILTERED
       , ROUTES_NOT_EXPORTED

       , useRoutesReceived
       , useRoutesFiltered 
       , useRoutesNotExported 
       }
  from 'app/context/routes';

import { useScrollToAnchor }
  from 'app/components/navigation/hash';

import QuickLinks
  from 'app/components/routes/QuickLinks';
import EmptyResults
  from 'app/components/routes/EmptyResults';
import RoutesTable
  from 'app/components/routes/RoutesTable';
import RouteDetailsModal
  from 'app/components/routes/RouteDetailsModal';
import Paginator
  from 'app/components/pagination/Paginator';
import PaginationInfo
  from 'app/components/pagination/PaginationInfo';
import LoadingIndicator
  from 'app/components/spinners/LoadingIndicator';


const RoutesHeader = ({type}) => {
  const rtype = {
    [ROUTES_RECEIVED]: "accepted",
    [ROUTES_FILTERED]: "filtered",
    [ROUTES_NOT_EXPORTED]: "not exported"
  }[type];
  const cls = `card-header card-header-routes ${type}`;
  return (<p className={cls}>Routes {rtype}</p>);
};


const RoutesLoading = () => {
  return (
    <div className={`card routes-view`}>
      <LoadingIndicator />
    </div>
  );
};


const createRoutesSet = (type, useRoutes) => () => {
  const results = useRoutes();
  const pageKey = {
    [ROUTES_RECEIVED]: PARAM_PAGE_RECEIVED,
    [ROUTES_FILTERED]: PARAM_PAGE_FILTERED,
    [ROUTES_NOT_EXPORTED]: PARAM_PAGE_NOT_EXPORTED,
  }[type];
  const anchor = {
    [ROUTES_RECEIVED]: 'routes-received',
    [ROUTES_FILTERED]: 'routes-filtered',
    [ROUTES_NOT_EXPORTED]: 'routes-not-exported',
  }[type];

  if (!results.requested) {
    return null;
  }
  if (results.loading) {
    return <RoutesLoading />;
  }
  if (results.totalResults === 0) {
    return null; // Nothing to show here.
  }

  // Render the routes card
  return (
    <div className={`card routes-view ${type}`}>
      <div className="row">
        <div className="col-md-6 routes-header-container">
          <RoutesHeader type={type} />
        </div>
        <div className="col-md-6">
          <PaginationInfo results={results} />
        </div>
      </div>
      <RoutesTable results={results} />
      <center>
        <Paginator
          results={results}
          pageKey={pageKey}
          anchor={anchor} />
      </center>
    </div>
  );
};

const RoutesReceived = createRoutesSet(
  ROUTES_RECEIVED,
  useRoutesReceived,
);

const RoutesFiltered = createRoutesSet(
  ROUTES_FILTERED,
  useRoutesFiltered,
);

const RoutesNotExported = createRoutesSet(
  ROUTES_NOT_EXPORTED,
  useRoutesNotExported,
);

/**
 * Show a button to load routes not exported on demand.
 * IF config states loading routes shoud be done automatically
 * update the query parameter.
 */
const RoutesNotExportedRequest = () => {
  const { noexport } = useConfig();
  const { requested } = useRoutesNotExported();
  const [query, setQuery] = useQuery();
  const requestNotExported = useQueryLocation({
    ...query,
    [PARAM_LOAD_NOT_EXPORTED]: "1",
  });

  const onDemand = noexport?.load_on_demand;
  
  useEffect(() => {
    if (onDemand === false) {
      setQuery((q) => ({
        ...q,
        [PARAM_LOAD_NOT_EXPORTED]: "1",
      }));
    }
  }, [onDemand, setQuery]);

  if (requested) {
    return null;
  }
  return (
    <div className="card routes-view routes-not-exported">
      <div className="row">
        <div className="col-md-6">
          <RoutesHeader type={ROUTES_NOT_EXPORTED} />
        </div>
      </div>
      <p className="help">
        Due to the potentially high amount of routes not exported,
        they are only fetched on demand.
      </p>

      <Link to={requestNotExported}
        className="btn btn-block btn-danger">
         Load Routes Not Exported
      </Link>
    </div>
  );
}

/**
 * Show all routes
 */
const Routes = () => {
  const refReceived = useRef();
  const refFiltered = useRef();
  const refNotExported = useRef();

  // Scroll to anchor
  useScrollToAnchor({
    "#routes-received": refReceived,
    "#routes-filtered": refFiltered,
    "#routes-not-exported": refNotExported,
  });

  return (
    <div className="routes-view">

      <QuickLinks />
      <EmptyResults />

      <RouteDetailsModal />

      <div ref={refFiltered}>
        <RoutesFiltered />
      </div>

      <div ref={refReceived}>
        <RoutesReceived />
      </div>

      <div ref={refNotExported}>
        <RoutesNotExportedRequest />
        <RoutesNotExported />
      </div>

    </div>
  );
};

export default Routes;
