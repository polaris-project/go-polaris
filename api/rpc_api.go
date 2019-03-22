// Package api contains all rpc and rest-related api helper methods, structs, etc...
package api

import (
	"fmt"
	"path/filepath"

	"github.com/polaris-project/go-polaris/common"
	configProto "github.com/polaris-project/go-polaris/internal/proto/config"
	cryptoProto "github.com/polaris-project/go-polaris/internal/proto/crypto"
	configServer "github.com/polaris-project/go-polaris/internal/rpc/config"
	cryptoServer "github.com/polaris-project/go-polaris/internal/rpc/crypto"

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
	err := generateCert("rpc", []string{"localhost"}) // Generate cert

	if err != nil { // Check for errors
		return err // Return found error
	}

	configHandler := configProto.NewConfigServer(&configServer.Server{}, nil) // Get handler
	cryptoHandler := cryptoProto.NewCryptoServer(&cryptoServer.Server{}, nil) // Get handler

	mux := http.NewServeMux() // Init mux

	mux.Handle(configProto.ConfigPathPrefix, configHandler) // Set route handler
	mux.Handle(cryptoProto.CryptoPathPrefix, cryptoHandler) // Set route handler

	return http.ListenAndServeTLS(rpcAPI.URI, filepath.FromSlash(fmt.Sprintf("%s/rpcCert.pem", common.CertificatesDir)), filepath.FromSlash(fmt.Sprintf("%s/rpcKey.pem", common.CertificatesDir)), mux) // Start serving
}

/* END EXPORTED METHODS */
