// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"errors"
	"fmt"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/config"
)

var (
	// WorkingDagDB represents the current opened dag database.
	WorkingDagDB *bolt.DB

	// ErrDagAlreadyExists represents a situation in which a user attempted to overwrite an existing DAG.
	ErrDagAlreadyExists = errors.New("dag already exists")
)

/* BEGIN EXPORTED METHODS */

// NewDag creates a new dag with the given config, and writes the dag db to memory.
// The newly opened dag db is stored in the WorkingDagDB variable.
// If a dag with the given identifier already exists, an ErrDagAlreadyExists error is returned.
func NewDag(config *config.DagConfig) error {
	err := config.WriteToMemory() // Write dag config to persistent memory

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = common.CreateDirIfDoesNotExit(common.DbDir) // Make database directory

	if err != nil { // Check for errors
		return err // Return found error
	}

	dagDB, err := bolt.Open(filepath.FromSlash(fmt.Sprintf("%s/%s.db", common.DbDir, config.Identifier)), 0644, &bolt.Options{Timeout: 5 * time.Second}) // Open DB with timeout

	if err != nil { // Check for errors
		return err // Return found error
	}

	WorkingDagDB = dagDB // Set dag DB

	return nil // No error occurred, return nil
}

/* END EXPORTED METHODS */
