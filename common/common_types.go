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

/* BEGIN EXPORTED METHODS */

// NewHash converts a given hash, b, to a 32-byte mem-prefix-compliant hash.
func NewHash(b []byte) Hash {
	var hash Hash // Hash

	bCropped := b // Init cropped buffer

	if len(b) > len(hash) { // Check crop side
		bCropped = bCropped[len(bCropped)-HashLength:] // Crop
	}

	copy(hash[HashLength-len(bCropped):], bCropped) // Copy source

	return hash // Return hash value
}

// IsNil checks if a given hash is nil.
func (hash Hash) IsNil() bool {
	nilBytes := 0 // Init nil bytes buffer

	for _, byteVal := range hash[:] { // Iterate through hash
		if byteVal == 0 { // Check nil byte
			nilBytes++ // Increment nil bytes
		}
	}

	return nilBytes == HashLength
}

// Bytes converts a given hash to a byte array.
func (hash Hash) Bytes() []byte {
	return hash[:] // Return byte array value
}

/* END EXPORTED METHODS */
