// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/polaris-project/go-polaris/config"
	"github.com/polaris-project/go-polaris/crypto"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestNewDag tests the functionality of the NewDag() method.
func TestNewDag(t *testing.T) {
	dagConfig := config.NewDagConfig(nil, "test_network", 1) // Initialize new dag config with test genesis file.

	_, err := NewDag(dagConfig) // Initialize dag with dag config

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	WorkingDagDB.Close() // Close dag db
}

// TestGetTransactionByHash tests the functionality of the GetTransactionByHash() helper method.
func TestGetTransactionByHash(t *testing.T) {
	dagConfig := config.NewDagConfig(nil, "test_network", 1) // Initialize new dag config with test genesis file.

	dag, err := NewDag(dagConfig) // Initialize dag with dag config

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate ecdsa private key

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	transaction := NewTransaction(
		0,                                        // Nonce
		big.NewInt(0),                            // Amount
		crypto.AddressFromPrivateKey(privateKey), // Sender
		crypto.AddressFromPrivateKey(privateKey), // Recipient
		nil,                    // Parents
		1,                      // Gas limit
		big.NewInt(1000),       // Gas price
		[]byte("test payload"), // Payload
	) // Create new transaction

	err = SignTransaction(transaction, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = dag.AddTransaction(transaction) // Add transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	_, err = dag.GetTransactionByHash(transaction.Hash) // Get tx

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	WorkingDagDB.Close() // Close working dag db
}

/* END EXPORTED METHODS TESTS */
