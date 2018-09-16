
/*
 * Routing helpers: Make query link with props (q and sorting)
 */

export const makeQueryLinkProps = function(routing, query, sort, order) {
  return {
    pathname: routing.pathname,
    search: `?s=${sort}&o=${order}&q=${query}`
  };
}

