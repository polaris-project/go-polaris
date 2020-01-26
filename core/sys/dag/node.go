// Package dag implements a persisted, partially readable merklized graph data structure used to store transactions.
package dag

import "crypto"

// Node represents a point in a graph.
type Node struct {
	// hash is a blake3 hash of the data contained inside the node
	hash crypto.Hash
}
