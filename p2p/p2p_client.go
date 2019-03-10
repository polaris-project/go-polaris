package p2p

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/juju/loggo"

	"github.com/polaris-project/go-polaris/common"

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

var (
	// logger is the p2p package logger.
	logger = loggo.GetLogger("p2p")
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

	err = client.StartServingStream(GetStreamHeaderProtocolPath(network, RequestBestTransaction), client.HandleReceiveBestTransactionRequest) // Register best tx request handler

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = client.StartServingStream(GetStreamHeaderProtocolPath(network, RequestTransaction), client.HandleReceiveTransactionRequest) // Register transaction request handler

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

// StartIntermittentSync syncs the dag with a given context and duration.
func (client *Client) StartIntermittentSync(ctx context.Context, duration time.Duration) {
	for range time.Tick(duration) { // Sync every duration seconds
		err := client.SyncDag(ctx) // Sync dag

		if err != nil { // Check for errors
			logger.Errorf("intermittent sync errored: %s", err.Error()) // Log error
		}
	}
}

// SyncDag syncs the working dag.
func (client *Client) SyncDag(ctx context.Context) error {
	if WorkingHost == nil { // Check no host
		return ErrNoWorkingHost // Return found error
	}

	lastTransactionHashes, err := BroadcastDhtResult(ctx, WorkingHost, types.BestTransactionRequest, GetStreamHeaderProtocolPath(client.Network, RequestBestTransaction), client.Network, 64) // Get last transaction hashes

	if err != nil { // Check for errors
		return err // Return found error
	}

	occurrences := make(map[common.Hash]int64) // Occurrences of each transaction hash

	bestLastTransactionHash := lastTransactionHashes[0] // Init last transaction buffer

	for _, lastTransactionHash := range lastTransactionHashes { // Iterate through last transaction hashes
		occurrences[common.NewHash(lastTransactionHash)]++ // Increment occurrences of transaction hash

		if occurrences[common.NewHash(lastTransactionHash)] > occurrences[common.NewHash(bestLastTransactionHash)] { // Check better last hash
			bestLastTransactionHash = lastTransactionHash // Set best last transaction hash
		}
	}

	remoteBestTransaction, err := client.RequestTransactionWithHash(ctx, common.NewHash(bestLastTransactionHash), 16) // Get last transaction

	if err != nil { // Check for errors
		return err // Return found error
	}

	localBestTransaction, err := (*client.Validator).GetWorkingDag().GetBestTransaction() // Get local best transaction

	if err != nil { // Check for errors
		return err // Return found error
	}

	if bytes.Equal(remoteBestTransaction.Hash.Bytes(), localBestTransaction.Hash.Bytes()) { // Check equivalent best last transaction hashes
		return nil // Nothing to sync!
	}

	// TODO: Actually sync.

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

	return BroadcastDht(ctx, WorkingHost, transaction.Bytes(), GetStreamHeaderProtocolPath(client.Network, PublishTransaction), client.Network) // Broadcast transaction
}

// RequestTransactionWithHash requests a given transaction with a given hash from the network.
// Returns best response from peer sampling set nPeers.
func (client *Client) RequestTransactionWithHash(ctx context.Context, hash common.Hash, nPeers int) (*types.Transaction, error) {
	transactionBytes, err := BroadcastDhtResult(ctx, WorkingHost, hash.Bytes(), GetStreamHeaderProtocolPath(client.Network, RequestTransaction), client.Network, nPeers) // Request transaction

	if err != nil { // Check for errors
		return &types.Transaction{}, err // Return found error
	}

	occurrences := make(map[common.Hash]int64) // Occurrences of each transaction hash

	bestTransaction := types.TransactionFromBytes(transactionBytes[0]) // Init best transaction buffer

	for _, currentTransactionBytes := range transactionBytes { // Iterate through transaction bytes
		currentTransaction := types.TransactionFromBytes(currentTransactionBytes) // Deserialize

		occurrences[currentTransaction.Hash]++ // Increment occurrences

		if occurrences[currentTransaction.Hash] > occurrences[bestTransaction.Hash] { // Check better than last transaction
			*bestTransaction = *currentTransaction // Set best transaction
		}
	}

	return bestTransaction, nil // Return best transaction
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

// HandleReceiveTransactionRequest handles a new stream requesting transaction metadata with a given hash.
func (client *Client) HandleReceiveTransactionRequest(stream inet.Stream) {
	readWriter := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream)) // Init reader/writer for stream

	var targetHashBytes []byte // Initialize tx hash bytes buffer

	for readBytes, err := readWriter.ReadByte(); err != nil; { // Read until EOF
		targetHashBytes = append(targetHashBytes, readBytes) // Append read bytes
	}

	transaction, _ := (*client.Validator).GetWorkingDag().GetTransactionByHash(common.NewHash(targetHashBytes)) // Get transaction with hash

	readWriter.Write(transaction.Bytes()) // Write transaction bytes
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
