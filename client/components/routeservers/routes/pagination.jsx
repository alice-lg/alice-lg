
import React from 'react'
import {connect} from 'react-redux'

import {Link} from 'react-router'



const PageLink = function(props) {
  const linkPage = parseInt(props.page);
  const label = props.label || (linkPage + 1);

  if (props.disabled) {
    return <span>{label}</span>;
  }

  let pr = props.pageReceived;
  let pf = props.pageFiltered;
  let pn = props.pageNotExported;

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

  const search = `?pr=${pr}&pf=${pf}&pn=${pn}`;
  const hash   = `#${props.anchor}`;
  const linkTo = {
    pathname: props.routing.pathname,
    hash:     hash,
    search:   search,
  };


  return (
    <Link to={linkTo}>{label}</Link>
  );
}


class RoutesPaginatorView extends React.Component {

  render() {
    if (this.props.totalPages <= 1) {
      return null; // Nothing to paginate
    }


    const pageLinks = Array.from(Array(this.props.totalPages), (_, i) => {
      let className = "";
      if (i == this.props.page) {
        className = "active";
      }

      return (
        <li key={i} className={className}>
          <PageLink page={i}
                    routing={this.props.routing}
                    anchor={this.props.anchor}
                    pageReceived={this.props.pageReceived}
                    pageFiltered={this.props.pageFiltered}
                    pageNotExported={this.props.pageNotExported} />
        </li>
      );
    });

    let prevLinkClass = "";
    if (this.props.page == 0) {
      prevLinkClass = "disabled";
    }

    let nextLinkClass = "";
    if (this.props.page + 1 == this.props.totalPages) {
      nextLinkClass = "disabled";
    }

    return (
      <nav aria-label="Routes Pagination">
        <ul className="pagination">
          <li className={prevLinkClass}>
            <PageLink page={this.props.page - 1}
                      label="&laquo;"
                      disabled={this.props.page == 0}
                      routing={this.props.routing}
                      anchor={this.props.anchor}
                      pageReceived={this.props.pageReceived}
                      pageFiltered={this.props.pageFiltered}
                      pageNotExported={this.props.pageNotExported} />
          </li>
          {pageLinks}
          <li className={nextLinkClass}>
            <PageLink page={this.props.page + 1}
                      disabled={this.props.page + 1 == this.props.totalPages}
                      label="&raquo;"
                      routing={this.props.routing}
                      anchor={this.props.anchor}
                      pageReceived={this.props.pageReceived}
                      pageFiltered={this.props.pageFiltered}
                      pageNotExported={this.props.pageNotExported} />
          </li>
        </ul>
      </nav>
    );
  }
}


export const RoutesPaginator = connect(
  (state) => ({
      pageReceived:    state.routes.receivedPage,
      pageFiltered:    state.routes.filteredPage,
      pageNotExported: state.routes.notExportedPage,

      routing: state.routing.locationBeforeTransitions
  })
)(RoutesPaginatorView);


export class RoutesPaginationInfo extends React.Component {
  render() {
    const totalResults = this.props.totalResults;
    const perPage = this.props.pageSize;
    const start = this.props.page * perPage + 1;
    const end = Math.min(start + perPage - 1, totalResults);
    if (this.props.totalPages == 1) {
      let routes = "route";
      if (totalResults > 1) {
        routes = "routes";
      }

      return (
        <div className="routes-pagination-info pull-right">
          Showing <b>all</b> of <b>{totalResults}</b> {routes}
        </div>
      );
    }
    return (
      <div className="routes-pagination-info pull-right">
        Showing <b>{start} - {end}</b> of <b>{totalResults}</b> total routes
      </div>
     );
  }
}

