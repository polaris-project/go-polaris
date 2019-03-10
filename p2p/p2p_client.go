package p2p

import (
	"bytes"
	"context"
	"errors"
	"time"

	"github.com/juju/loggo"

	"github.com/polaris-project/go-polaris/common"

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

	localBestTransaction, _ := (*client.Validator).GetWorkingDag().GetBestTransaction() // Get local best transaction

	if localBestTransaction == nil { // Check no genesis
		genCtx, cancel := context.WithCancel(ctx) // Initialize context

		err = client.SyncGenesis(genCtx)

		if err != nil { // Check for errors
			cancel()

			return err // Return found error
		}

		cancel() // Cancel
	}

	if bytes.Equal(remoteBestTransaction.Hash.Bytes(), localBestTransaction.Hash.Bytes()) { // Check equivalent best last transaction hashes
		return nil // Nothing to sync!
	}

	return nil // No error occurred, return nil
}

// SyncGenesis syncs the local genesis transaction set for the working dag.
func (client *Client) SyncGenesis(ctx context.Context) error {
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

/*
	END TRANSACTION HELPERS
*/

/* END EXPORTED METHODS */
