
export const makeSearchQueryProps = function(query) {
  return {
    pathname: '/search',
    search: `?q=${query}`
  }
}

