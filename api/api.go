// Package api contains all rpc and rest-related api helper methods, structs, etc...
package api

// API defines a standardized API spec, that of which can be implemented via REST, RPC, etc...
type API interface {
	GetAPIProtocol() string // Get API protocol (e.g. REST_HTTP, RPC, etc...)

	GetAPIURI() string // Get API uri (e.g. https://localhost:8080/)

	GetSupportedFormats() []string // Get supported formats

	GetIsServing() bool // Check is serving

	StartServingAPI() error // Start serving API
}
