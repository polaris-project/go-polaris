// Package dag implements a persisted, partially readable merklized graph data structure used to store transactions.
package dag

import (
	"math/big"
	"testing"

	"github.com/polaris-project/go-polaris/core/sys/primitives"
	"github.com/polaris-project/go-polaris/crypto"
)

// TestNodeHash tests the functionality of the node Hash helper method.
func TestNodeHash(t *testing.T) {
	// Make a zero-value transaction, just for testing
	tx := primitives.NewTransaction(big.NewInt(0), crypto.Hash{}, crypto.Hash{}, big.NewInt(0), nil, nil, nil)

	// Make a node with the zero-value transaction as its contents
	node := Node{tx, nil}

	// Get the hash of the node
	nodeHash := node.Hash()

	// Ensure that the node's hash is known
	if nodeHash.IsZero() {
		// Panic
		t.Fatal("node should have a valid hash")
	}
}
