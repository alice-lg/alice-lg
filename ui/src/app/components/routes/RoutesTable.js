
import { useCallback }
  from 'react';


import { useConfig }
  from 'app/context/config';

import FilterReason
  from 'app/components/routes/FilterReason';
import NoExportReason 
  from 'app/components/routes/NoExportReason';


// Helper: Lookup value in route path
export const getAttr = (r, path) => {
  return path.split(".").reduce((acc, elem) => acc[elem], r);
}


export const ColDefault = ({onClick, route, column}) => {
  return (
    <td>
      <span onClick={onClick}>{getAttr(route, column)}</span>
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

/*
export const ColFlags = ({route}) => {
  return (
    <td className="col-route-flags">
      <span className="route-prefix-flags">
        <RpkiIndicator route={route} />
        <PrimaryIndicator route={route} />
        <BlackholeIndicator route={route}
                            blackholes={blackholes} />
        <RejectCandidateIndicator route={route} />
      </span>
    </td>
  );
}
*/

const RouteColumn = ({onClick, column, route}) => {
  const widgets = {
    "network": ColNetwork,
    //"flags": ColFlags,
    "bgp.as_path": ColAsPath,

   // "Flags": ColFlags,
    "ASPath": ColAsPath,
  };

  const handleClick = useCallback(() => onClick(route), [route, onClick]);

  let Widget = widgets[column] || ColDefault;
  return (
    <Widget
      column={column}
      route={route}
      onClick={handleClick} />
  );
};



const RoutesTable = ({results}) => {
  const config = useConfig();

  const columns = config.routes_columns;
  const columnsOrder = config.routes_columns_order;

  const { routes } = results;

  const showAttributesModal = useCallback((route) => {
    console.log("show attributes:", route);
  }, []);

  if(!routes.length === 0) {
    return null;
  }

  const rows = routes.map((r, i) => (
    <tr key={i}>
      {columnsOrder.map((col) => (
        <RouteColumn
          key={col}
          onClick={showAttributesModal}
          column={col}
          route={r} />
      ))}
    </tr>
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
