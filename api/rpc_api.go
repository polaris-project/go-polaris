// Package api contains all rpc and rest-related api helper methods, structs, etc...
package api

// RPCAPI represents an RPC API.
type RPCAPI struct {
	Network string `json:"network"` // API network

	SupportedProtocols []string `json:"protocols"` // Api protocols
}

/* BEGIN EXPORTED METHODS */

// NewRPCAPI initializes a new RPCAPI instance.
func NewRPCAPI(network string, protocols []string) *RPCAPI {
	// TODO: Generate certs
	return &RPCAPI{
		Network:            network,   // Set network
		SupportedProtocols: protocols, // Set protocols
	} // Return initialized API
}

/* END EXPORTED METHODS */
