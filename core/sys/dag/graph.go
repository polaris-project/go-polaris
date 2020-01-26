// Package dag implements a persisted, partially readable merklized graph data structure used to store transactions.
package dag

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/crypto"
)

// Option represents a generic functional configuration option applied to an instance of the Dag type.
type Option func(d Dag) (Dag, error)

// Dag is the global history of transactions from the beginning to end of time.
type Dag struct {
	// Nodes is a list of nodes stored in the DAG
	Nodes []Node

	// HashRoutes represents a mapping of hashes to their corresponding node indexes
	HashRoutes map[crypto.Hash]*Node

	// Children represents a mapping of parent hashes to a slice of their corresponding children
	Children map[crypto.Hash][]*Node

	// DB represents a connector to the local Dag persistence service.
	DB *bolt.DB
}

// NewDag initializes a new dag with the provided configuration options.
func NewDag(opts ...Option) (Dag, error) {
	// Open a bolt database that we can use to partially load the dag from
	db, err := bolt.Open(fmt.Sprintf("%s/dag.db", common.DefaultDbDir), 0600, nil)

	// Check for any errors
	if err != nil {
		return Dag{}, err // Return the error
	}

	// Make the DAG instance with the opened database, and otherwise nil values
	dag := Dag{
		Nodes:      []Node{},
		HashRoutes: make(map[crypto.Hash]*Node),
		Children:   make(map[crypto.Hash][]*Node),
		DB:         db,
	}

	// Iterate through each of the options in the provided options slice
	for _, opt := range opts {
		// Apply the option to the dag
		dag, err = opt(dag)

		// Check for any errors that occurred whilst applying the option
		if err != nil {
			// Return the error
			return dag, err
		}
	}

	// Return the opened DAG instance
	return dag, nil
}

// WithDataDir generates a configuration option for the Dag type whereby the Dag
// is initialized or opened in the given data directory.
func WithDataDir(dir string) Option {
	// Make the configuration option, and return it
	return func(d Dag) (Dag, error) {
		// Open the database
		db, err := bolt.Open(fmt.Sprintf("%s/db/database.db", dir), 0600, nil)

		// Check for any errors that occurred whilst opening the database
		if err != nil {
			// Don't apply the option, just return the original dag instance
			return d, err
		}

		// Set the dag's DB to the new DB
		d.DB = db

		// Return the new DAG
		return d, nil
	}
}

// LoadPartial generates a configuration option that forces the DAG to load only state headers from the database.
func LoadPartial() Option {
	// Generate the option and return it
	return func(d Dag) (Dag, error) {

	}
}
