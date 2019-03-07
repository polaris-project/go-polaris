package p2p

import (
	"context"
	"errors"

	inet "github.com/libp2p/go-libp2p-net"
	protocol "github.com/libp2p/go-libp2p-protocol"
	"github.com/polaris-project/go-polaris/types"
	"github.com/polaris-project/go-polaris/validator"
)

var (
	// ErrNoWorkingHost represents an error describing a WorkingHost value of nil.
	ErrNoWorkingHost = errors.New("no working host")

	// ErrNilHash defines an error describing a situation in which a message has no hash.
	ErrNilHash = errors.New("hash not set")
)

// Client represents an active p2p peer, that of which is serving a list of available stream header protocol paths.
type Client struct {
	Network string `json:"network"` // Active network

	Validator *validator.Validator // Validator
}

/* BEGIN EXPORTED METHODS */

// NewClient initializes a new client
func NewClient(network string, validator *validator.Validator) *Client {
	return &Client{
		Network:   network,   // Set network
		Validator: validator, // Set validator
	}
}

// StartServingStream starts serving a stream on a given header protocol path.
func (client *Client) StartServingStream(streamHeaderProtocolPath string, handler func(inet.Stream)) error {
	if WorkingHost == nil { // Check no host
		return ErrNoWorkingHost // Return found error
	}

	WorkingHost.SetStreamHandler(protocol.ID(streamHeaderProtocolPath), handler) // Set handler

	return nil // No error occurred, return nil
}

// PublishTransaction publishes a given transaction.
func (client *Client) PublishTransaction(ctx context.Context, transaction *types.Transaction) error {
	if WorkingHost == nil { // Check no host
		return ErrNoWorkingHost // Return found error
	}

	if err := (*client.Validator).ValidateTransaction(transaction); err != nil { // Validate transaction
		return err // Return found error
	}

	context, cancel := context.WithCancel(ctx) // Get context

	defer cancel() // Cancel

	return BroadcastDht(context, WorkingHost, transaction.Bytes(), GetStreamHeaderProtocolPath(client.Network, PublishTransaction), client.Network) // Broadcast transaction
}

/* END EXPORTED METHODS */
