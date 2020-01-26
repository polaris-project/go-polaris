// Package dag implements a persisted, partially readable merklized graph data structure used to store transactions.
package dag

import (
	"github.com/polaris-project/go-polaris/core/sys/primitives"
	"github.com/polaris-project/go-polaris/crypto"
)

// Node represents a point in a graph.
type Node struct {
	// transaction is the transaction stored inside the node
	Transaction primitives.Transaction

	// Children is a slice containing each of the children of the node
	Children []*Node
}

// Hash gets the hash of a particular node.
func (n *Node) Hash() crypto.Hash {
	// Return the hash of the node's transaction
	return n.Transaction.Hash()
}
