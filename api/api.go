// Package api contains all rpc and rest-related api helper methods, structs, etc...
package api

// API defines a standardized API spec, that of which can be implemented via REST, RPC, etc...
type API interface {
	GetApiProtocol() string // Get API protocol (e.g. REST JSON, RPC JSON, etc...)

	GetApiURI() string // Get API uri (e.g. https://localhost:8080/)

	StartServingAPI() error // Start serving API
}
