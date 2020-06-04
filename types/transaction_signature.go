// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"errors"
	"math/big"

	"github.com/polaris-project/go-polaris/common"
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
	MarshaledPublicKey []byte `json:"pub" gencodec:"required"` // Signature public key

	V []byte   `json:"v" gencodec:"required"` //  Signature message
	R *big.Int `json:"r" gencodec:"required"` // Signature retrieval
	S *big.Int `json:"s" gencodec:"required"` // Signature retrieval
}

/* BEGIN EXPORTED METHODS */

// SignTransaction signs a given transaction via ecdsa, and sets the transaction signature to the new signature.
// Returns a new signature composed of v, r, s values.
// If the transaction has already been signed, returns an ErrAlreadySigned error, as well as a nil signature pointer.
// If the transaction has no hash, returns an ErrNilHash error, as well as a nil signature pointer.
func SignTransaction(transaction *Transaction, privateKey *ecdsa.PrivateKey) error {
	if transaction.Hash.IsNil() { // Check no existing hash
		return ErrNilHash // Return no hash error
	}

	if transaction.Signature == nil { // Check no existing signature
		r, s, err := ecdsa.Sign(rand.Reader, privateKey, transaction.Hash.Bytes()) // Sign via ECDSA
		if err != nil {                                                            // Check for errors
			return err // Return found error
		}

		(*transaction).Signature = &Signature{
			MarshaledPublicKey: elliptic.Marshal(elliptic.P521(), privateKey.PublicKey.X, privateKey.PublicKey.Y), // Set marshaled public key
			V:                  transaction.Hash.Bytes(),                                                          // Set hash
			R:                  r,                                                                                 // Set R
			S:                  s,                                                                                 // Set S
		} // Set transaction signature

		(*transaction).Hash = common.NewHash(nil) // Set hash to nil

		(*transaction).Hash = crypto.Sha3(transaction.Bytes()) // Set transaction hash

		return nil // Return signature
	}

	return ErrAlreadySigned // Return already signed error
}

// SignMessage signs a given message hash via ecdsa, and returns a new signature
func SignMessage(messageHash common.Hash, privateKey *ecdsa.PrivateKey) (*Signature, error) {
	if messageHash.IsNil() { // Check nil hash
		return nil, ErrNilHash // Return no hash error
	}

	r, s, err := ecdsa.Sign(rand.Reader, privateKey, messageHash.Bytes()) // Sign via ECDSA
	if err != nil {                                                       // Check for errors
		return nil, err // Return found error
	}

	signature := &Signature{
		V: messageHash.Bytes(), // Set hash
		R: r,                   // Set R
		S: s,                   // Set S
	} // Set transaction signature

	return signature, nil // Return signature
}

// Verify checks that a given signature is valid, and returns whether or not the given signature is valid.
// If no signature exists at the given memory address, false is returned.
func (signature *Signature) Verify(address *common.Address) bool {
	x, y := elliptic.Unmarshal(elliptic.P521(), signature.MarshaledPublicKey) // Unmarshal public key

	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P521(), // Set curve
		X:     x,               // Set x
		Y:     y,               // Set y
	} // Recover public key

	if signature == nil { // Check no existent signature
		return false // No signature to verify
	}

	if *crypto.AddressFromPublicKey(publicKey) != *address { // Check invalid public key
		return false // Invalid
	}

	return ecdsa.Verify(publicKey, signature.V, signature.R, signature.S) // Verify signature contents
}

/* END EXPORTED METHODS */
