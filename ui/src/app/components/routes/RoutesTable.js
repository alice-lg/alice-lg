
import { useCallback
       , useMemo
       }
  from 'react';


import { useRoutesTableConfig
       }
  from 'app/context/config';
import { useSetRouteDetails }
  from 'app/context/routes';


import FilterReason
  from 'app/components/routes/FilterReason';
import NoExportReason 
  from 'app/components/routes/NoExportReason';
import PrimaryIndicator
  from 'app/components/routes/flags/PrimaryIndicator';
import RpkiIndicator
  from 'app/components/routes/flags/RpkiIndicator';
import BlackholeIndicator
  from 'app/components/routes/flags/BlackholeIndicator';
import RejectCandidateIndicator
  from 'app/components/routes/flags/RejectCandidateIndicator';

// Helper: Lookup value in route path
export const getAttr = (r, path) => {
  return path.split(".").reduce((acc, elem) => acc[elem], r);
}

// Linking: Create link targes as a function of the route
// Link to the route server
const linkRouteServer = (route) => 
  `/routeservers/${route?.routeserver?.id}`;

// Create a link to the routes page of a neighbor
const linkNeighborRoutes = (route) => {
  const rs = route?.routeserver?.id;
  const neighbor = route?.neighbor_id;
  return `/routeservers/${rs}/neighbors/${neighbor}/routes`;
}

// Default column: Show the attribute and bind the
// onClick attribute.
export const ColDefault = ({onClick, route, column}) => {
  return (
    <td>
      <span onClick={onClick}>{getAttr(route, column)}</span>
    </td>
  );
}

// ColLink provides a cell with a linkable target.
// The attribute `to` is a function of the `route`
// attribute, returning the url.
export const ColLink = ({to, route, column}) => {
  const href = to(route);
  return (
    <td>
      <a href={href} target="_blank" rel="noreferrer">
       {getAttr(route, column)}
      </a>
    </td>
  );
}

// Include filter and noexport reason in this column.
export const ColNetwork = ({onClick, route}) => {
  return (
    <td className="col-route-network">
      <span className="route-network" onClick={onClick}>
        {route.network} 
      </span>
      <FilterReason route={route} />
      <NoExportReason route={route} />
    </td>
  );
}

// Special AS Path Widget
export const ColAsPath = ({route}) => {
    let asns = getAttr(route, "bgp.as_path");
    if(!asns){
      asns = [];
    }
    const baseUrl = "https://irrexplorer.nlnog.net/asn/AS"
    let asnLinks = asns.map((asn, i) => {
      return (<a key={`${asn}_${i}`} href={baseUrl + asn} target="_blank" rel="noreferrer">{asn} </a>);
    });

    return (
      <td>
        {asnLinks}
      </td>
    );
}

export const ColFlags = ({route}) => {
  return (
    <td className="col-route-flags">
      <span className="route-prefix-flags">
        <RpkiIndicator route={route} />
        <PrimaryIndicator route={route} />
        <BlackholeIndicator route={route} />
        <RejectCandidateIndicator route={route} />
      </span>
    </td>
  );
}


export const ColRouteServer = ({route, column}) =>
  <ColLink to={linkRouteServer} route={route}  column={column} />;


export const ColNeighbor = ({route, column}) =>
  <ColLink to={linkNeighborRoutes} route={route} column={column} />;


const RouteColumn = ({onClick, column, route}) => {
  const cells = {
    "network": ColNetwork,
    "flags": ColFlags,
    "bgp.as_path": ColAsPath,

    "Flags": ColFlags,
    "ASPath": ColAsPath,

    "routeserver.name": ColRouteServer,
    "neighbor.description": ColNeighbor,
  };

  let Cell = cells[column] || ColDefault;
  return (
    <Cell
      column={column}
      route={route}
      onClick={onClick} />
  );
};


/**
 * RoutesRow renders a memoized row
 */
const RoutesRow = ({columns, route, onClick}) => {
  return useMemo(() => {
    const callback = () => onClick(route);
    const cols = columns.map((col) => (
      <RouteColumn
        key={col}
        onClick={callback}
        column={col}
        route={route} />
    ));
    return (
      <tr>{cols}</tr>
    );
  }, [columns, route, onClick]);
}


const RoutesTable = ({results}) => {
  const setRouteDetails = useSetRouteDetails();
  const { columns, columnsOrder } = useRoutesTableConfig();

  const { routes } = results;

  const showAttributesModal = useCallback((route) => {
    setRouteDetails(route);
  }, [setRouteDetails]);

  if(!routes.length === 0) {
    return null;
  }

  const rows = routes.map((r, i) => (
    <RoutesRow
      key={i}
      columns={columnsOrder}
      onClick={showAttributesModal}
      route={r} />
  ));

  return (
    <table className="table table-striped table-routes">
      <thead>
        <tr>
          {columnsOrder.map(col => <th key={col}>{columns[col]}</th>)}
        </tr>
      </thead>
      <tbody>
        {rows}
      </tbody>
    </table>
  );
}

export default RoutesTable;
