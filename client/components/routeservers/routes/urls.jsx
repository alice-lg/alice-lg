
import {filtersUrlEncode} from 'components/filters/encoding'

export const makeLinkProps = function(props) {
  const linkPage = parseInt(props.page, 10);

  let pr = props.pageReceived;
  let pf = props.pageFiltered;
  let pn = props.pageNotExported;
  let ne = props.loadNotExported;

  // Numeric flags
  ne = ne ? 1 : 0;

  // This here can be surely more elegant.
  switch(props.anchor) {
    case "routes-received":
      pr = linkPage;
      break;
    case "routes-filtered":
      pf = linkPage;
      break;
    case "routes-not-exported":
      pn = linkPage;
      break;
  }

  let filtering = "";
  if (props.filtersApplied) {
    filtering = filtersUrlEncode(props.filtersApplied);
  }

  const query = props.routing.query.q || "";
  const search = `?ne=${ne}&pr=${pr}&pf=${pf}&pn=${pn}&q=${query}${filtering}`;

  let hash = null;
  if (props.anchor) {
    hash = `#${props.anchor}`;
  }

  const linkTo = {
    pathname: props.routing.pathname,
    hash:     hash,
    search:   search,
  };

  return linkTo;
}


export const makePeerLinkProps = function(rsId, protocolId) {
  const linkTo = {
    pathname: `/routeservers/${rsId}/protocols/${protocolId}/routes`,
  };

  return linkTo;
}

