// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/polaris-project/go-polaris/crypto"
	"github.com/polaris-project/go-polaris/p2p"
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

// TestPublish tests the functionality of the Publish() helper method.
func TestPublish(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	address := crypto.AddressFromPrivateKey(privateKey) // Generate address

	transaction := NewTransaction(
		0,                      // Nonce
		big.NewFloat(0),        // Amount
		address,                // Sender
		nil,                    // Recipient
		nil,                    // Parents
		0,                      // Gas limit
		big.NewInt(0),          // Gas price
		[]byte("test payload"), // Payload
	) // Initialize a new transaction

	err = SignTransaction(transaction, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	ctx, cancel := context.WithCancel(context.Background()) // Get context

	defer cancel() // Cancel

	_, err = p2p.NewHost(ctx, 3861) // Initialize host

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = transaction.Publish(context.Background(), "test_network") // Publish transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

/* END EXPORTED METHODS TESTS */
