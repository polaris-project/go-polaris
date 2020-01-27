// Package dag implements a persisted, partially readable merklized graph data structure used to store transactions.
package dag

import (
	"fmt"

	"github.com/boltdb/bolt"
	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/core/sys/primitives"
	"github.com/polaris-project/go-polaris/crypto"
)

// StateBucket is a byte-slice representation of the name of the database bucket in which
// state data is held.
var StateBucket = []byte("state")

// TransactionsBucket is a byte-slice representation of the name of the database bucket in which
// transaction data is held.
var TransactionsBucket = []byte("transactions")

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

// AddNode adds the given node to the DAG, returning a reference to said node.
func (d *Dag) AddNode(n Node) (*Node, error) {
	// Add the node to the list of nodes contained in the DAG
	d.Nodes = append(d.Nodes, n)

	// Get a reference to the node we just added to the DAG
	addedNode := &d.Nodes[len(d.Nodes)-1]

	// Establish a route to the node by its transaction hash
	d.HashRoutes[n.Hash()] = addedNode

	// Iterate through each of the parents defined by the transaction
	for _, parentHash := range n.Transaction.Parents {
		// Try to get a mutable reference to the node containing the parent transaction
		if parentNode, ok := d.HashRoutes[parentHash]; ok {
			// Register the node as a child of the parent node
			parentNode.Children = append(parentNode.Children, addedNode)

			// Update the parent -> children DAG mapping to reflect the new node that we added to the graph
			if childrenMapEntry, ok := d.Children[parentHash]; ok {
				// Register the node as a child in the Children mapping
				childrenMapEntry = append(childrenMapEntry, addedNode)
			}
		}
	}

	// Add the node to the database
	return addedNode, d.DB.Update(func(tx *bolt.Tx) error {
		// Get a reference to the bucket that we'll use to store transactions
		b := tx.Bucket(TransactionsBucket)

		// Serialize the transaction so that we can put it into the database
		serialized, err := n.Transaction.Serialize()

		// Check for any errors that might have arisen whilst serializing
		if err != nil {
			// Return the error
			return err
		}

		// Store the transaction's bytes in the database
		return b.Put(n.Hash(), serialized)
	})
}

// AddTransaction creates a new node, and fills it with the given transaction, which will
// be added to the dag. This new node is returned.
func (d *Dag) AddTransaction(tx primitives.Transaction) (*Node, error) {
	// Make a new node with the given transaction
	n := Node{tx, nil}

	// Add the node to the dag
	return d.AddNode(n)
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

	// Create the necessary buckets for the DAG database
	err = dag.createBuckets()

	// Check for any errors
	if err != nil {
		// Return the errors
		return Dag{}, err
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

		// Create the necessary buckets in the DAG's database
		return d, d.createBuckets()
	}
}

// LoadPartial generates a configuration option that forces the DAG to load only state headers from the database.
func LoadPartial() Option {
	// Generate the option and return it
	return func(d Dag) (Dag, error) {
		// Load a list of nodes contained inside the graph's database
		d.DB.View(func(tx *bolt.Tx) error {
			// Get a reference to the database's transactions bucket
			b := tx.Bucket(transactions)

			// Get a cursor for the transactions bucket, so we can start iterating over all the transactions
			// and add them to the database
			c := b.Cursor()

			// Iterate through each of the nodes placed into the database
			for k, v := c.First(); k != nil; k, v = c.Next() {

			}
		})
	}
}

// createBuckets creates each of the bolt buckets necessary for proper operation of the DAG--
// namely, the transaction and state buckets.
func (d *Dag) createBuckets() error {
	// Create each of the necessary buckets
	d.DB.Update(func(tx *bolt.Tx) error {
		// Create the transactions bucket
		_, err := tx.CreateBucketIfNotExists(TransactionsBucket)
		// Check for any errors that occurred whilst creating the bucket
		if err != nil {
			return err
		}

		// Create the state bucket
		_, err = tx.CreateBucketIfNotExists(StateBucket)

		// Return any errors that might have occurred
		return err
	})
}
