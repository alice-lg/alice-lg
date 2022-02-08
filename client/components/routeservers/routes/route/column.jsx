
/*
 * Routes Rendering Columns
 */
import _ from 'underscore'

import React from 'react'
import {connect} from 'react-redux'

import FilterReason
  from 'components/routeservers/communities/filter-reason'

import NoexportReason
  from 'components/routeservers/communities/noexport-reason'

import {ROUTES_RECEIVED,
        ROUTES_FILTERED,
        ROUTES_NOT_EXPORTED} from '../actions'


import {PrimaryIndicator,
        BlackholeIndicator,
        RpkiIndicator,
        RejectCandidateIndicator} from './flags'


// Helper: Lookup value in route path
export const _lookup = (r, path) => {
  return path.split(".").reduce((acc, elem) => acc[elem], r);
}


export const ColDefault = function(props) {
  return (
    <td>
      <span onClick={props.onClick}>{_lookup(props.route, props.column)}</span>
    </td>
  );
}

// Include filter and noexport reason in this column.
export const ColNetwork = function(props) {
  return (
    <td className="col-route-network">
      <span className="route-network" onClick={props.onClick}>
        {props.route.network} 
      </span>
      <FilterReason route={props.route} />
      <NoexportReason route={props.route} />
    </td>
  );
}

// Special AS Path Widget
export const ColAsPath = function(props) {
    let asns = _lookup(props.route, "bgp.as_path");
    if(!asns){
        asns = [];
    }
    const baseUrl = "https://irrexplorer.nlnog.net/asn/AS"

    let asnLinks = asns.map((asn, i) => {
      return (<a key={`${asn}_${i}`} href={baseUrl + asn} target="_blank">{asn} </a>);
    });

    return (
        <td>
          {asnLinks}
        </td>
    );
}


export const ColFlags = function(props) {
  return (
    <td className="col-route-flags">
      <span className="route-prefix-flags">
        <RpkiIndicator route={props.route} />
        <PrimaryIndicator route={props.route} />
        <BlackholeIndicator route={props.route}
                            blackholes={props.blackholes} />
        <RejectCandidateIndicator route={props.route} />
      </span>
    </td>
  );
}


// Meta component, decides what to render based on on 
// prop 'column'.
export default function(props) {
  const widgets = {
    "network": ColNetwork,
    "flags": ColFlags,
    "bgp.as_path": ColAsPath,

    "Flags": ColFlags,
    "ASPath": ColAsPath,
  };

  let Widget = widgets[props.column] || ColDefault;
  return (
    <Widget column={props.column} route={props.route}
            blackholes={props.blackholes}
            onClick={props.onClick} />
  );
}


