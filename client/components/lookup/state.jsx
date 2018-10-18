

/*
 * Manage state
 */

import {
  filtersUrlEncode
} from './filter-encoding'


/* 
 * Maybe this can be customized and injected into 
 * the PageLink component.
 */
export function  makeLinkProps(props) {
  const linkPage = parseInt(props.page, 10);

  let pr = props.pageReceived;
  let pf = props.pageFiltered;

  // This here can be surely more elegant.
  switch(props.anchor) {
    case "received":
      pr = linkPage;
      break;
    case "filtered":
      pf = linkPage;
      break;
  }

  let filtering = "";
  if (props.filtersApplied) {
    filtering = filtersUrlEncode(props.filtersApplied);
  }

  const query = props.routing.query.q || "";

  const search = `?pr=${pr}&pf=${pf}&q=${query}${filtering}`;
  const hash   = `#routes-${props.anchor}`;
  const linkTo = {
    pathname: props.routing.pathname,
    hash:     hash,
    search:   search,
  };

  return linkTo;
}

