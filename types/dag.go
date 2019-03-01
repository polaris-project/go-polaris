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

	if WorkingDagDB.
}

/*
	BEGIN DB READING HELPER METHODS
*/



/*
	END DB READING HELPER METHODS
*/

/* END EXPORTED METHODS */
