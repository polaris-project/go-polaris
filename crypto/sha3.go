// Package crypto provides cryptography helper methods.
package crypto

import (
	"encoding/hex"

	"github.com/polaris-project/go-polaris/common"
	"golang.org/x/crypto/sha3"
)

/* BEGIN EXPORTED METHODS */

// Sha3 hashes a given message via sha3.
func Sha3(b []byte) common.Hash {
	hash := sha3.New256() // Init hasher

	hash.Write(b) // Write

	return common.NewHash(hash.Sum(nil)) // Return final hash
}

// Sha3String hashes a given message via sha3 and encodes the hashed message to a hex string.
func Sha3String(b []byte) string {
	b = Sha3(b).Bytes() // Hash

	return hex.EncodeToString(b) // Return string
}

// Sha3n hashes a given message via sha3 n times.
func Sha3n(b []byte, n uint) common.Hash {
	hashSource := b // Init editable buffer

	for x := uint(0); x != n; x++ { // Hash n times
		hashSource = Sha3(hashSource).Bytes() // Hash
	}

	return common.NewHash(hashSource) // Return hashed
}

// Sha3nString hashes a given message via sha3 n times and encodes the hashed message to a hex string.
func Sha3nString(b []byte, n uint) string {
	return hex.EncodeToString(Sha3n(b, n).Bytes()) // Return string hash
}

// Sha3d hashes a given message twice via sha3.
func Sha3d(b []byte) common.Hash {
	return Sha3(Sha3(b).Bytes()) // Return sha3d result
}

// Sha3dString hashes a given message via sha3d and encodes the hashed message to a hex string.
func Sha3dString(b []byte) string {
	return hex.EncodeToString(Sha3d(b).Bytes()) // Return string hash
}

/* END EXPORTED METHODS */
