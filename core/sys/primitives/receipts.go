// Package primitives implements a series of basic types required by the network.
package primitives

import "github.com/polaris-project/go-polaris/crypto"

// LogKey represents a key for which a particular value is stored in a transaction's logs.
type LogKey [16]byte

// Receipt represents the result of a state transition executed on the polaris runtime.
type Receipt struct {
	// Logs are the output of some past transaction
	Logs map[LogKey][]byte

	// StateRoot is the hash of the state root corresponding to the receipt
	StateRoot crypto.Hash
}
