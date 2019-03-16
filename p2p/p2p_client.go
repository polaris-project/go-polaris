package p2p

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"time"

	"github.com/boltdb/bolt"
	"github.com/juju/loggo"

	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/crypto"

	"github.com/polaris-project/go-polaris/types"
	"github.com/polaris-project/go-polaris/validator"
)

var (
	// ErrNoWorkingHost represents an error describing a WorkingHost value of nil.
	ErrNoWorkingHost = errors.New("no working host")

	// ErrNilHash defines an error describing a situation in which a message has no hash.
	ErrNilHash = errors.New("hash not set")

	// ErrNoAvailablePeers defines an error describing an available peer sampling set with a length of 0.
	ErrNoAvailablePeers = errors.New("no available peers")
)

var (
	// logger is the p2p package logger.
	logger = getLogger()
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
			logger.Errorf("intermittent sync errored (if private net, this is expected): %s", err.Error()) // Log error
		}
	}
}

// SyncDag syncs the working dag.
func (client *Client) SyncDag(ctx context.Context) error {
	logger.Infof("starting dag sync") // Log start dag sync

	if WorkingHost == nil { // Check no host
		return ErrNoWorkingHost // Return found error
	}

	bestLastTransactionHash, err := client.RequestBestTransactionHash(ctx, 64) // Request best tx hash

	if err != nil { // Check for errors
		return err // Return found error
	}

	remoteBestTransaction, err := client.RequestTransactionWithHash(ctx, bestLastTransactionHash, 16) // Get last transaction

	if err != nil { // Check for errors
		return err // Return found error
	}

	logger.Infof("dag sync: determined must sync up to transaction with hash %s", hex.EncodeToString(remoteBestTransaction.Hash.Bytes())) // Log must sync up to

	if (*client.Validator).GetWorkingDag().Genesis.IsNil() { // Check no genesis
		logger.Infof("couldn't find a valid genesis transaction; syncing") // Log sync genesis

		genCtx, cancel := context.WithCancel(ctx) // Initialize context

		err = client.SyncGenesis(genCtx) // Sync genesis

		if err != nil { // Check for errors
			cancel()

			return err // Return found error
		}

		logger.Infof("finished syncing gensis") // Log finish syncing genesis

		cancel() // Cancel
	}

	logger.Infof("syncing best transaction") // Log sync best transaction

	return client.SyncBestTransaction(ctx, remoteBestTransaction.Hash) // No error occurred, return nil
}

// SyncBestTransaction syncs the best local and remote transactions.
func (client *Client) SyncBestTransaction(ctx context.Context, remoteBestTransactionHash common.Hash) error {
	localBestTransaction, err := (*client.Validator).GetWorkingDag().GetBestTransaction() // Get local best transaction

	if err != nil { // Check for errors
		return err // Return found error
	}

	logger.Infof("syncing best transaction from local best transaction: %s and remote transaction: %s", hex.EncodeToString(localBestTransaction.Hash.Bytes()), hex.EncodeToString(remoteBestTransactionHash.Bytes())) // Log sync best transaction

	for !bytes.Equal(remoteBestTransactionHash.Bytes(), localBestTransaction.Hash.Bytes()) { // Do until valid best last transaction hash
		getChildrenCtx, cancel := context.WithCancel(ctx) // Initialize context

		logger.Infof("requesting children for current best transaction: %s", hex.EncodeToString(localBestTransaction.Hash.Bytes())) // Log request children

		childHashes, err := client.RequestTransactionChildren(getChildrenCtx, localBestTransaction.Hash, 16) // Get child hashes

		if err != nil { // Check for errors
			cancel() // Cancel

			return err // Return found error
		}

		logger.Infof("found %d children for current best transaction: %s", len(childHashes), hex.EncodeToString(localBestTransaction.Hash.Bytes())) // Log children

		cancel() // Cancel

		for _, childHash := range childHashes { // Iterate through child hashes
			requestTransactionCtx, cancel := context.WithCancel(ctx) // Initialize context

			destinationTransaction, err := client.RequestTransactionWithHash(requestTransactionCtx, childHash, 16) // Get transaction

			if err != nil { // Check for errors
				cancel() // Cancel

				return err // Return found error
			}

			cancel() // Cancel

			if err := (*client.Validator).ValidateTransaction(destinationTransaction); err == nil { // Check valid transaction
				logger.Infof("adding child: %s", hex.EncodeToString(destinationTransaction.Hash.Bytes())) // Log add children

				err = (*client.Validator).GetWorkingDag().AddTransaction(destinationTransaction) // Add transaction to local dag

				if err != nil { // Check for errors
					return err // Return found error
				}
			}
		}

		localBestTransaction, _ = (*client.Validator).GetWorkingDag().GetBestTransaction() // Get local best transaction
	}

	return nil // No error occurred, return nil
}

// SyncGenesis syncs the local genesis transaction set for the working dag.
func (client *Client) SyncGenesis(ctx context.Context) error {
	if WorkingHost == nil { // Check no host
		return ErrNoWorkingHost // Return found error
	}

	getGenHashCtx, cancel := context.WithCancel(ctx) // Get context

	genesisHashes, err := BroadcastDhtResult(getGenHashCtx, WorkingHost, types.GenesisHashRequest, GetStreamHeaderProtocolPath(client.Network, RequestGenesisHash), client.Network, 128) // Get genesis transaction hashes

	if err != nil { // Check for errors
		cancel() // Cancel

		return err // Return found error
	}

	cancel() // Cancel

	occurrences := make(map[common.Hash]int64) // Occurrences of each transaction hash

	bestGenesisHash := genesisHashes[0] // Init best genesis hash buffer

	for _, genesisHash := range genesisHashes { // Iterate through genesis hashes
		if bytes.Equal(genesisHash, make([]byte, len(genesisHash))) { // Check is nil
			continue // Continue
		}

		occurrences[common.NewHash(genesisHash)]++ // Increment occurrences of genesis hash

		if occurrences[common.NewHash(genesisHash)] > occurrences[common.NewHash(bestGenesisHash)] { // Check better genesis hash
			bestGenesisHash = genesisHash // Set best genesis hash
		}
	}

	getGenCtx, cancel := context.WithCancel(ctx) // Get context

	genesisTransaction, err := client.RequestTransactionWithHash(getGenCtx, common.NewHash(bestGenesisHash), 16) // Get genesis transaction

	if err != nil { // Check for errors
		cancel() // cancel

		return err // Return found error
	}

	cancel() // Cancel

	if !genesisTransaction.Hash.IsNil() { // Ensure is not nil genesis
		_, err = (*client.Validator).GetWorkingDag().GetTransactionByHash(genesisTransaction.Hash) // Get transaction by hash

		if err == nil { // Check tx already exists
			return validator.ErrDuplicateTransaction // Return found error
		}

		return types.WorkingDagDB.Update(func(tx *bolt.Tx) error {
			_, err := tx.CreateBucketIfNotExists([]byte("transaction-bucket")) // Create tx bucket if it doesn't already exist

			if err != nil { // Check for errors
				return err // Return found error
			}

			workingTransactionBucket := tx.Bucket([]byte("transaction-bucket")) // Get transaction bucket

			err = workingTransactionBucket.Put(genesisTransaction.Hash.Bytes(), genesisTransaction.Bytes()) // Put transaction

			if err != nil { // Check for errors
				return err // Return found error
			}

			(*client.Validator).GetWorkingDag().Genesis = genesisTransaction.Hash // Set genesis hash

			return (*client.Validator).GetWorkingDag().WriteToMemory() // Write genesis to db header
		}) // Write genesis transaction
	}

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

		if currentTransaction.Hash.IsNil() { // Check hash is nil
			continue // Continue
		}

		occurrences[currentTransaction.Hash]++ // Increment occurrences

		if occurrences[currentTransaction.Hash] > occurrences[bestTransaction.Hash] { // Check better than last transaction
			*bestTransaction = *currentTransaction // Set best transaction
		}
	}

	return bestTransaction, nil // Return best transaction
}

// RequestTransactionChildren requests the children of a transaction from a sampling set of nPeers size.
func (client *Client) RequestTransactionChildren(ctx context.Context, parentHash common.Hash, nPeers int) ([]common.Hash, error) {
	if WorkingHost == nil { // Check no host
		return []common.Hash{}, ErrNoWorkingHost // Return error
	}

	childHashesAllResponses, err := BroadcastDhtResult(ctx, WorkingHost, parentHash.Bytes(), GetStreamHeaderProtocolPath(client.Network, RequestChildHashes), client.Network, nPeers) // Request child hashes

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	occurrences := make(map[common.Hash]int64) // Occurrences of each transaction hash

	bestChildHashSet := []common.Hash{}      // Init best hash set buffer
	bestChildHashSetHashSum := common.Hash{} // Init best hash set hash sum buffer

	for _, childHashes := range childHashesAllResponses { // Iterate through children
		hashes := bytes.Split(childHashes, []byte("end_hash")) // Separate hashes by delimiter

		var castedHashes []common.Hash // Init casted buffer

		for _, childHash := range hashes { // Iterate through hashes
			castedHashes = append(castedHashes, common.NewHash(childHash)) // Append casted hash
		}

		if bytes.Equal(bytes.Join(hashes, []byte{}), make([]byte, len(hashes[0])*len(hashes))) { // Check is nil
			continue // Continue
		}

		occurrences[crypto.Sha3(bytes.Join(hashes, []byte{}))]++ // Increment occurrences

		if occurrences[crypto.Sha3(bytes.Join(hashes, []byte{}))] > occurrences[bestChildHashSetHashSum] { // Check new best hash set
			bestChildHashSet = castedHashes                                     // Set best hash set
			bestChildHashSetHashSum = crypto.Sha3(bytes.Join(hashes, []byte{})) // Set best hash set hash sum
		}
	}

	return bestChildHashSet, nil // No error occurred, return children
}

// RequestBestTransactionHash returns the average best tx hash between nPeers.
func (client *Client) RequestBestTransactionHash(ctx context.Context, nPeers int) (common.Hash, error) {
	lastTransactionHashes, err := BroadcastDhtResult(ctx, WorkingHost, types.BestTransactionRequest, GetStreamHeaderProtocolPath(client.Network, RequestBestTransaction), client.Network, nPeers) // Get last transaction hashes

	if err != nil { // Check for errors
		return common.Hash{}, err // Return found error
	}

	occurrences := make(map[common.Hash]int64) // Occurrences of each transaction hash

	if len(lastTransactionHashes) == 0 { // Check no peers
		return common.Hash{}, ErrNoAvailablePeers // Return error
	}

	bestLastTransactionHash := lastTransactionHashes[0] // Init last transaction buffer

	for _, lastTransactionHash := range lastTransactionHashes { // Iterate through last transaction hashes
		if bytes.Equal(lastTransactionHash, make([]byte, len(lastTransactionHash))) { // Check nil
			continue // Continue
		}

		occurrences[common.NewHash(lastTransactionHash)]++ // Increment occurrences of transaction hash

		if occurrences[common.NewHash(lastTransactionHash)] > occurrences[common.NewHash(bestLastTransactionHash)] { // Check better last hash
			bestLastTransactionHash = lastTransactionHash // Set best last transaction hash
		}
	}

	return common.NewHash(bestLastTransactionHash), nil // Return best hash
}

/*
	END TRANSACTION HELPERS
*/

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// getLogger gets the p2p package logger, and sets the levels of said logger.
func getLogger() loggo.Logger {
	logger := loggo.GetLogger("p2p") // Get logger

	loggo.ConfigureLoggers("p2p=INFO") // Configure loggers

	return logger // Return logger
}

/* END INTERNAL METHODS */
