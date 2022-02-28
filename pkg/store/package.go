// Package store provides an interface for keeping
// routes and neighbor information without querying
// a route server. The refresh happens in a configurable
// interval.
//
// There a currently two implementations: A postgres and
// and in-memory backend.
package store
