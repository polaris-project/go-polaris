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

// TestNewTransactions tests the functionality of the NewTransaction method.
func TestNewTransactions(t *testing.T) {
	transaction := NewTransaction(
		0,                      // Nonce
		big.NewFloat(10),       // Amount
		nil,                    // Sender
		nil,                    // Recipient
		nil,                    // Parents
		1,                      // Gas limit
		big.NewInt(1000),       // Gas price
		[]byte("test payload"), // Payload
	) // Create a new transaction using the NewTransaction method

	t.Log(transaction) // Log the initialized transaction
}

// TestCalculateTotalValue tests the functionality of the CalculateTotalValue() helper method.
func TestCalculateTotalValue(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	address := crypto.AddressFromPrivateKey(privateKey) // Generate address

	transaction := NewTransaction(
		0,                      // Nonce
		big.NewFloat(5),        // Amount
		address,                // Sender
		nil,                    // Recipient
		nil,                    // Parents
		1000,                   // Gas limit
		big.NewInt(1),          // Gas price
		[]byte("test payload"), // Payload
	) // Initialize a new transaction

	if transaction.CalculateTotalValue().Cmp(big.NewFloat(5+1000)) != 0 { // Check invalid value calculation
		t.Fatal("invalid total value calculation") // Panic
	}
}

/* END EXPORTED METHODS TESTS */
