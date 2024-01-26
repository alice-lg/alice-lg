package api

import "errors"

// ErrTooManyRoutes is returned when the result set
// of a route query exceeds the maximum allowed number of routes.
var ErrTooManyRoutes = errors.New("too many routes")
