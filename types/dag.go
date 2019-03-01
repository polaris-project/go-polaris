// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/config"
)

var (
	transactionBucket = []byte("transaction-bucket")
)

var (
	// WorkingDagDB represents the current opened dag database.
	WorkingDagDB *bolt.DB

	// ErrDagAlreadyExists represents an error describing
	// the attempted overwriting of an existing DAG.
	ErrDagAlreadyExists = errors.New("dag already exists")

	// ErrDagDbNotOpened represents an error describing the attempted appending of a data set
	// to the working dag, despite the fact that the dag db has not yet been opened.
	ErrDagDbNotOpened = errors.New("dag db has not been opened")

	// ErrNilTransaction represents an error describing a transaction pointer of nil value.
	ErrNilTransaction = errors.New("transaction pointer is nil")

	// ErrNilTransactionAtHash represents an error describing a transaction pointer of nil value discovered through the
	// querying of the working db for an invalid hash.
	ErrNilTransactionAtHash = errors.New("no transaction exists in the dag with given hash")
)

// Dag is a simple struct used to abstract db reading and writing methods.
type Dag struct {
	DagConfig *config.DagConfig `json:"config"` // Dag config
}

/* BEGIN EXPORTED METHODS */

// NewDag creates a new dag with the given config, and writes the dag db to memory.
// The newly opened dag db is stored in the WorkingDagDB variable.
func NewDag(config *config.DagConfig) (*Dag, error) {
	err := config.WriteToMemory() // Write dag config to persistent memory

	if err != nil { // Check for errors
		return &Dag{}, err // Return found error
	}

	err = common.CreateDirIfDoesNotExit(common.DbDir) // Make database directory

	if err != nil { // Check for errors
		return &Dag{}, err // Return found error
	}

	if _, err := os.Stat(filepath.FromSlash(fmt.Sprintf("%s/%s.db", common.DbDir, config.Identifier))); err == nil { // Check db already exists
		return &Dag{}, ErrDagAlreadyExists // Return found error
	}

	dagDB, err := bolt.Open(filepath.FromSlash(fmt.Sprintf("%s/%s.db", common.DbDir, config.Identifier)), 0644, &bolt.Options{Timeout: 5 * time.Second}) // Open DB with timeout

	if err != nil { // Check for errors
		return &Dag{}, err // Return found error
	}

	WorkingDagDB = dagDB // Set dag DB

	return &Dag{
		DagConfig: config, // Set config
	}, nil // Return dag
}

// AddTransaction appends a given transaction to the working dag.
// Returns an ErrDagDbNotOpened error if the working dag db is nil (has been not opened).
// Return an ErrNilTransaction error if the given transaction pointer is nil.
// Returns an ErrDuplicateTransaction error if the transaction already exists in the working dag db.
// Returns an ErrTransactionNotSigned error if the transaction has not been signed.
// Return an ErrSignatureInvalid error if the transaction's signature is invalid.
func (dag *Dag) AddTransaction(transaction *Transaction) error {
	if WorkingDagDB == nil { // Check dag db not opened
		return ErrDagDbNotOpened // Return found error
	}

	if transaction == nil { // Check nil pointer
		return ErrNilTransaction // Return error
	}

	return nil // No error occurred, return nil
}

/*
	BEGIN DB READING HELPER METHODS
*/

// GetTransactionByHash attempts to query the working dag db by the given transaction hash.
// If no transaction exists at this hash, an
func (dag *Dag) GetTransactionByHash(transactionHash common.Hash) (*Transaction, error) {
	var txBytes []byte // Init buffer

	err := WorkingDagDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(transactionBucket) // Get tx bucket

		txBytes = bucket.Get(transactionHash.Bytes()) // Get tx at hash

		if txBytes == nil { // Check no transaction at hash
			return ErrNilTransactionAtHash // Return error
		}

		return nil // No error occurred, return nil
	})

	if err != nil { // Check for errors
		return &Transaction{}, err // Return found error
	}

	return TransactionFromBytes(txBytes), nil // Return deserialized tx
}

/*
	END DB READING HELPER METHODS
*/

/* END EXPORTED METHODS */
