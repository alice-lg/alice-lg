
import {urlEscape} from 'components/utils/query'

export const makeSearchQueryProps = function(query) {
  query = urlEscape(query); 
  return {
    pathname: '/search',
    search: `?q=${query}`
  }
}

