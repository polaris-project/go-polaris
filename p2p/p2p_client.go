package p2p

import (
	"bufio"
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

// StartServingStreams attempts to start serving all necessary streams
func (client *Client) StartServingStreams(network string) error {
	err := client.StartServingStream(GetStreamHeaderProtocolPath(network, PublishTransaction), client.HandleReceiveTransaction) // Register tx handler

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = client.StartServingStream(GetStreamHeaderProtocolPath(network, RequestConfig), client.HandleReceiveConfigRequest) // Register config request handler

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// StartServingStream starts serving a stream on a given header protocol path.
func (client *Client) StartServingStream(streamHeaderProtocolPath string, handler func(inet.Stream)) error {
	if WorkingHost == nil { // Check no host
		return ErrNoWorkingHost // Return found error
	}

	WorkingHost.SetStreamHandler(protocol.ID(streamHeaderProtocolPath), handler) // Set handler

	return nil // No error occurred, return nil
}

// SyncDag syncs the working dag.
func (client *Client) SyncDag() error {
	// lastTransactionHashes :=
	return nil // No error occurred, return nil
}

/*
	BEGIN TRANSACTION HELPERS
*/

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

// HandleReceiveTransaction handles a new stream sending a transaction.
func (client *Client) HandleReceiveTransaction(stream inet.Stream) {
	reader := bufio.NewReader(stream) // Initialize reader from stream

	var transactionBytes []byte // Initialize transaction bytes buffer

	for readBytes, err := reader.ReadByte(); err != nil; { // Read until EOF
		transactionBytes = append(transactionBytes, readBytes) // Append read bytes
	}

	transaction := types.TransactionFromBytes(transactionBytes) // Deserialize transaction

	if err := (*client.Validator).ValidateTransaction(transaction); err == nil { // Check transaction valid
		(*client.Validator).GetWorkingDag().AddTransaction(transaction) // Add transaction to working dag
	}
}

// HandleReceiveBestTransactionRequest handle a new stream requesting for the best transaction hash.
func (client *Client) HandleReceiveBestTransactionRequest(stream inet.Stream) {
	writer := bufio.NewWriter(stream) // Initialize writer from stream

	bestTransaction, _ := (*client.Validator).GetWorkingDag().GetBestTransaction() // Get best transaction

	writer.Write(bestTransaction.Bytes()) // Write best transaction
}

/*
	END TRANSACTION HELPERS
*/

/*
	BEGIN CONFIG HELPERS
*/

// HandleReceiveConfigRequest handles a new stream requesting a
func (client *Client) HandleReceiveConfigRequest(stream inet.Stream) {
	writer := bufio.NewWriter(stream) // Initialize writer

	writer.Write((*client.Validator).GetWorkingConfig().Bytes()) // Write config bytes
}

/*
	END CONFIG HELPERS
*/

/* END EXPORTED METHODS */
