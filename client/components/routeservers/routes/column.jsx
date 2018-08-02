
/*
 * Routes Rendering Columns
 */

import React from 'react'


// Helper:
export const PrimaryIndicator = function(props) {
  if (props.route.details && props.route.primary) {
    return(
      <span className="primary-route is-primary-route">&gt;
        <div>Best Route</div>
      </span>
    );
  }

  // Default
  return (
    <span className="primary-route not-primary-route"></span>
  )
}


export const ColDefault = function(props) {
  return (
    <td>
      <span onClick={props.onClick}>{_lookup(props.route, props.column)}</span>
    </td>
  )
}

// Include filter and noexport reason in this column.
export const ColNetwork = function(props) {
  return (
    <td className="col-route-network">
      <span className="route-network" onClick={props.onClick}>
        <PrimaryIndicator route={props.route} />
        {props.route.network}
      </span>
      {props.displayReasons == ROUTES_FILTERED && <FilterReason route={props.route} />}
      {props.displayReasons == ROUTES_NOT_EXPORTED && <NoexportReason route={props.route} />}
    </td>
  );
}

// Special AS Path Widget
export const ColAsPath = function(props) {
    const asns = _lookup(props.route, "bgp.as_path");
    const baseUrl = "http://irrexplorer.nlnog.net/search/"

    let asnLinks = asns.map((asn, i) => {
      return (<a key={`${asn}_${i}`} href={baseUrl + asn} target="_blank">{asn} </a>);
    });

    return (
        <td>
          {asnLinks}
        </td>
    );
}



// Meta component, decides what to render based on on 
// prop 'column'.
export default RouteColumn = function(props) {
  const widgets = {
    "network": ColNetwork,
    "bgp.as_path": ColAsPath,

    "ASPath": ColAsPath,
  };

  let Widget = widgets[props.column] || ColDefault;
  return (
    <Widget column={props.column} route={props.route}
            displayReasons={props.displayReasons}
            onClick={props.onClick} />
  );
}


