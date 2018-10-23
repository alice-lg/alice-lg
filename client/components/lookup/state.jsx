

/*
 * Manage state
 */

import {
  filtersUrlEncode
} from 'components/filters/filter-encoding'

import {
  FILTER_GROUP_SOURCES,
  FILTER_GROUP_ASNS,
  FILTER_GROUP_COMMUNITIES,
  FILTER_GROUP_EXT_COMMUNITIES,
  FILTER_GROUP_LARGE_COMMUNITIES,
} from 'components/filters/filter-groups'

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

  let pagination = "";
  if (pr) {
    pagination += `pr=${pr}`
  }
  if (pf) {
    pagination += `pf=${pf}`
  }

  let filtering = "";
  if (props.filtersApplied) {
    filtering = filtersUrlEncode(props.filtersApplied);
  }

  const query = props.routing.query.q || "";

  const search = `?${pagination}&q=${query}${filtering}`;
  let hash = "";
  if (props.anchor) {
    hash += `#routes-${props.anchor}`;
  }

  const linkTo = {
    pathname: props.routing.pathname,
    hash:     hash,
    search:   search,
  };

  return linkTo;
}

