// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/polaris-project/go-polaris/config"
	"github.com/polaris-project/go-polaris/crypto"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestNewDag tests the functionality of the NewDag() method.
func TestNewDag(t *testing.T) {
	os.RemoveAll(filepath.FromSlash("data/db/test_network.db")) // Remove existing db

	dagConfig := config.NewDagConfig(nil, "test_network", 1) // Initialize new dag config with test genesis file.

	_, err := NewDag(dagConfig) // Initialize dag with dag config

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	WorkingDagDB.Close() // Close dag db

	os.RemoveAll(filepath.FromSlash("data/db/test_network.db")) // Remove existing db
}

// TestGetTransactionByHash tests the functionality of the GetTransactionByHash() helper method.
func TestGetTransactionByHash(t *testing.T) {
	os.RemoveAll(filepath.FromSlash("data/db/test_network.db")) // Remove existing db

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
		big.NewFloat(0),                          // Amount
		crypto.AddressFromPrivateKey(privateKey), // Sender
		crypto.AddressFromPrivateKey(privateKey), // Recipient
		nil,                                      // Parents
		1,                                        // Gas limit
		big.NewInt(1000),                         // Gas price
		[]byte("test payload"),                   // Payload
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

	os.RemoveAll(filepath.FromSlash("data/db/test_network.db")) // Remove existing db
}

// TestGetTransactionByHash tests the functionality of the GetTransactionByAddress() helper method.
func TestGetTransactionByAddress(t *testing.T) {
	os.RemoveAll(filepath.FromSlash("data/db/test_network.db")) // Remove existing db

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
		big.NewFloat(0),                          // Amount
		crypto.AddressFromPrivateKey(privateKey), // Sender
		crypto.AddressFromPrivateKey(privateKey), // Recipient
		nil,                                      // Parents
		1,                                        // Gas limit
		big.NewInt(1000),                         // Gas price
		[]byte("test payload"),                   // Payload
	) // Create new transaction

	err = SignTransaction(transaction, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = dag.AddTransaction(transaction) // Add transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	transactions, err := dag.GetTransactionsByAddress(crypto.AddressFromPrivateKey(privateKey)) // Get transactions related to sender

	if err != nil { // Check for errors
		t.Fatal(err) // panic
	}

	if len(transactions) != 1 { // Check invalid tx set
		t.Fatalf("should have found 1 related transaction; found %d", len(transactions)) // Log invalid filter
	}

	WorkingDagDB.Close() // Close working dag db

	os.RemoveAll(filepath.FromSlash("data/db/test_network.db")) // Remove existing db
}

// TestCalculateAddressBalance tests the functionality of the CalculateAddressBalance() helper method.
func TestCalculateAddressBalance(t *testing.T) {
	os.RemoveAll(filepath.FromSlash("data/db/test_network.db")) // Remove existing db

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
		big.NewFloat(1),                          // Amount
		crypto.AddressFromPrivateKey(privateKey), // Sender
		nil,                                      // Recipient
		nil,                                      // Parents
		1,                                        // Gas limit
		big.NewInt(1000),                         // Gas price
		[]byte("test payload"),                   // Payload
	) // Create new transaction

	err = SignTransaction(transaction, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = dag.forceAddTransaction(transaction) // Add transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	balance, err := dag.CalculateAddressBalance(crypto.AddressFromPrivateKey(privateKey)) // Calculate balance

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	if balance.Cmp(big.NewFloat(-1001.0)) != 0 { // Check invalid balance
		t.Fatal("invalid balance calculation") // Panic
	}

	WorkingDagDB.Close() // Close working dag db

	os.RemoveAll(filepath.FromSlash("data/db/test_network.db")) // Remove existing db
}

/* END EXPORTED METHODS TESTS */
