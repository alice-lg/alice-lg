
/*
 * Routing helpers: Make query link with props (q and sorting)
 */

import {urlEscape} from 'components/utils/query'

export const makeQueryLinkProps = function(routing, query, sort, order) {
  return {
    pathname: routing.pathname,
    search: `?s=${sort}&o=${order}&q=${urlEscape(query)}`
  };
}

