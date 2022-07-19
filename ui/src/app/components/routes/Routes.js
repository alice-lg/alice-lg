
import { useRef }
  from 'react';

import { ROUTES_RECEIVED
       , ROUTES_FILTERED
       , ROUTES_NOT_EXPORTED

       , useRoutesReceived
       , useRoutesFiltered 
       , useRoutesNotExported 
       }
  from 'app/context/routes';

import QuickLinks
  from 'app/components/routes/QuickLinks';
import EmptyResults
  from 'app/components/routes/EmptyResults';
import Paginator
  from 'app/components/pagination/Paginator';
import PaginationInfo
  from 'app/components/pagination/PaginationInfo';


export const RoutesHeader = ({type}) => {
  const rtype = {
    [ROUTES_RECEIVED]: "accepted",
    [ROUTES_FILTERED]: "filtered",
    [ROUTES_NOT_EXPORTED]: "not exported"
  }[type];
  const cls = `card-header card-header-routes ${type}`;
  return (<p className={cls}>Routes {rtype}</p>);
};


const createRoutesSet = (type, useRoutes) => () => {
  const results = useRoutes();
  const pageKey = {
    [ROUTES_RECEIVED]: 'pr',
    [ROUTES_FILTERED]: 'pf',
    [ROUTES_NOT_EXPORTED]: 'pn',
  }[type];
  const anchor = {
    [ROUTES_RECEIVED]: 'routes-received',
    [ROUTES_FILTERED]: 'routes-filtered',
    [ROUTES_NOT_EXPORTED]: 'routes-not-exported',
  }[type];

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
        [RoutesTable]
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


/**
 * Show all routes
 */
const Routes = () => {
  const refReceived = useRef();
  const refFiltered = useRef();
  const refNotExported = useRef();

  return (
    <div className="routes-view">

      <QuickLinks />
      <EmptyResults />

      <div ref={refReceived}>
        <RoutesReceived />
      </div>

      <div ref={refReceived}>
      </div>

      <div ref={refNotExported}>
      </div>

    </div>
  );
};

export default Routes;
