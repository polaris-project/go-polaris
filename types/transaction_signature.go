// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"crypto/ecdsa"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/polaris-project/go-polaris/crypto"
)

var (
	// ErrAlreadySigned defines an error describing a situation in which a message has already been signed.
	ErrAlreadySigned = errors.New("already signed")

	// ErrNilHash defines an error describing a situation in which a message has no hash.
	ErrNilHash = errors.New("hash not set")
)

// Signature is a data type representing a verifiable ECDSA signature--that of which
// is not necessarily a transaction signature.
type Signature struct {
	V []byte   `json:"v" gencodec:"required"` //  Signature message
	R *big.Int `json:"r" gencodec:"required"` // Signature retrieval
	S *big.Int `json:"s" gencodec:"required"` // Signature retrieval
}

/* BEGIN EXPORTED METHODS */

// SignTransaction signs a given transaction via ecdsa, and sets the transaction signature to the new signature.
// Returns a new signature composed of v, r, s values.
// If the transaction has already been signed, returns an ErrAlreadySigned error, as well as a nil signature pointer.
// If the transaction has no hash, returns an ErrNilHash error, as well as a nil signature pointer.
func SignTransaction(transaction *Transaction, privateKey *ecdsa.PrivateKey) (*Signature, error) {
	if transaction.Hash.IsNil() { // Check no existing hash
		return &Signature{}, ErrNilHash // Return no hash error
	}

	if transaction.Signature == nil { // Check no existing signature
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, transaction.Hash.Bytes()) // Sign via ECDSA

		if err != nil { // Check for errors
			return &Signature{}, err // Return found error
		}

		(*transaction).Signature = &Signature{
			V: transaction.Hash.Bytes(), // Set hash
			R: r,                        // Set R
			S: s,                        // Set S
		} // Set transaction signature

		(*transaction).Hash = crypto.Sha3(transaction.Bytes()) // Set transaction hash

		return transaction.Signature, nil // Return signature
	}

	return &Signature{}, ErrAlreadySigned // Return already signed error
}

/* END EXPORTED METHODS */
