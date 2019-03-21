// Package api contains all rpc and rest-related api helper methods, structs, etc...
package api

import (
	"context"
	"net/http"
)

// RPCAPI represents an RPC API.
type RPCAPI struct {
	Network string `json:"network"` // API network

	SupportedProtocols []string `json:"protocols"` // API protocols

	URI string `json:"uri"` // API URI

	Server *http.Server `json:"server"` // Working server
}

/* BEGIN EXPORTED METHODS */

// NewRPCAPI initializes a new RPCAPI instance.
func NewRPCAPI(network string, uri string) (*RPCAPI, error) {
	err := generateCert("rpc", []string{"localhost", "127.0.0.1"}) // Generate tls certs

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return &RPCAPI{
		Network: network, // Set network
		URI:     uri,     // Set URI
	}, nil // Return initialized API
}

// NewRPCAPINoTLS initializes a new RPCAPI instance without enabling TLS.
func NewRPCAPINoTLS(network string, uri string) *RPCAPI {
	return &RPCAPI{
		Network: network, // Set network
		URI:     uri,     // Set URI
	} // Return initialized API
}

// GetAPIProtocol gets the working rpc api protocol.
func (rpcAPI *RPCAPI) GetAPIProtocol() string {
	return "RPC" // Return protocol
}

// GetAPIURI gets the current rpc api URI.
func (rpcAPI *RPCAPI) GetAPIURI() string {
	return rpcAPI.URI // Return URI
}

// GetSupportedFormats gets the formats supported by the current rpc api.
func (rpcAPI *RPCAPI) GetSupportedFormats() []string {
	return []string{"protobuf"} // Return format
}

// GetIsServing returns whether or not the api is currently serving.
func (rpcAPI *RPCAPI) GetIsServing() bool {
	return rpcAPI.Server != nil // Check is serving
}

// StartServing starts serving the API.
func (rpcAPI *RPCAPI) StartServing(ctx context.Context) error {
	// TODO: Serve

	return nil // No error occurred, return nil
}

/* END EXPORTED METHODS */
