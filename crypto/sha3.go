// Package crypto provides cryptography helper methods.
package crypto

import (
	"encoding/hex"

	"golang.org/x/crypto/sha3"
)

// Sha3 hashes a given message via sha3.
func Sha3(b []byte) []byte {
	hash := sha3.New256() // Init hasher

	hash.Write(b) // Write

	return hash.Sum(nil) // Return final hash
}

// Sha3String hashes a given message via sha3 and encodes the hashed message to a hex string.
func Sha3String(b []byte) string {
	b = Sha3(b) // Hash

	return hex.EncodeToString(b) // Return string
}

// Sha3n hashes a given message via sha3 n times.
func Sha3n(b []byte, n uint) []byte {
	for x := uint(0); x != n; x++ { // Hash n times
		b = Sha3(b) // Hash
	}

	return b // Return hashed
}

// Sha3nString hashes a given message via sha3 n times and encodes the hashed message to a hex string.
func Sha3nString(b []byte, n uint) string {
	b = Sha3n(b, n) // Hash

	return hex.EncodeToString(b) // Return string
}

// Sha3d hashes a given message twice via sha3.
func Sha3d(b []byte) []byte {
	return Sha3(Sha3(b)) // Return sha3d result
}

// Sha3dString hashes a given message via sha3d and encodes the hashed message to a hex string.
func Sha3dString(b []byte) string {
	b = Sha3d(b) // Hash

	return hex.EncodeToString(b) // Return string
}
