// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/polaris-project/go-polaris/crypto"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestSignTransaction tests the functionality of the SignTransaction() method.
func TestSignTransaction(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate ecdsa private key

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	transaction := &Transaction{
		AccountNonce: 0,                // Set nonce
		Sender:       nil,              // Set sender
		Recipient:    nil,              // Set recipient
		GasPrice:     big.NewInt(1000), // Set gas price
		Payload:      []byte("test"),   // Set payload
		Signature:    nil,              // Set signature
	} // Initialize transaction

	transaction.Hash = crypto.Sha3(transaction.Bytes()) // Set hash

	_, err = SignTransaction(transaction, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	if transaction.Signature == nil { // Check was not signed
		t.Fatal("transaction not signed") // Panic
	}
}

// TestVerifySignature tests the functionality of the VerifySignature() method.
func TestVerifySignature(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate ecdsa private key

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	transaction := &Transaction{
		AccountNonce: 0,                // Set nonce
		Sender:       nil,              // Set sender
		Recipient:    nil,              // Set recipient
		GasPrice:     big.NewInt(1000), // Set gas price
		Payload:      []byte("test"),   // Set payload
		Signature:    nil,              // Set signature
	} // Initialize transaction

	transaction.Hash = crypto.Sha3(transaction.Bytes()) // Set hash

	_, err = SignTransaction(transaction, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	if !transaction.Signature.Verify(&privateKey.PublicKey) { // Check that signature is valid
		t.Fatal("signature should be valid") // Panic
	}
}

/* END EXPORTED METHODS TESTS */
