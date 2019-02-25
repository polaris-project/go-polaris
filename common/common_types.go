// Package common defines a set of commonly used helper methods and data types.
package common

const (
	// HashLength is the standardized length of a hash
	HashLength = 32

	// AddressLength is the standardized length of an address
	AddressLength = 20
)

// Hash represents the 32 byte output of sha3().
type Hash [HashLength]byte

// Address represents a 20 byte, hex-encoded ECDSA public key.
type Address [AddressLength]byte
