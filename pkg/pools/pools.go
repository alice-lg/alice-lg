// Package pools provides deduplication pools for strings
// and lists of ints and strings.
package pools

import "log"

// Default pools: These pools are defined globally
// and are defined per intended usage

// Neighbors stores neighbor IDs
var Neighbors *String

// Networks4 stores network ip v4 addresses
var Networks4 *String

// Networks6 stores network ip v6 addresses
var Networks6 *String

// Interfaces stores interfaces like: eth0, bond0 etc...
var Interfaces *String

// Gateways4 store ip v4 gateway addresses
var Gateways4 *String

// Gateways6 store ip v6 gateway addresses
var Gateways6 *String

// Origins is a store for 'IGP'
var Origins *String

// ASPaths stores lists of ASNs
var ASPaths *IntList

// Types stores a list of types (['BGP', 'univ'])
var Types *StringList

// Initialize global pools
func init() {
	log.Println("initializing memory pools")

	Neighbors = NewString()
	Networks4 = NewString()
	Networks6 = NewString()
	Interfaces = NewString()
	Gateways4 = NewString()
	Gateways6 = NewString()
	Origins = NewString()
	ASPaths = NewIntList()
	Types = NewStringList()
}
