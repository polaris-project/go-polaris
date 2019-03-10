// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/polaris-project/go-polaris/common"
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

// TestGetTransactionByAddress tests the functionality of the GetTransactionByAddress() helper method.
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

// TestGetTransactionsBySender tests the functionality of the GetTransactionBySender() helper method.
func TestGetTransactionsBySender(t *testing.T) {
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

	transactions, err := dag.GetTransactionsBySender(crypto.AddressFromPrivateKey(privateKey)) // Get transactions from sender

	if err != nil { // Check for errors
		t.Fatal(err) // panic
	}

	if len(transactions) != 1 { // Check invalid tx set
		t.Fatalf("should have found 1 related transaction; found %d", len(transactions)) // Log invalid filter
	}

	WorkingDagDB.Close() // Close working dag db

	os.RemoveAll(filepath.FromSlash("data/db/test_network.db")) // Remove existing db
}

// TestGetTransactionChildren tests the functionality of the GetTransactionChildren() helper method.
func TestGetTransactionChildren(t *testing.T) {
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

	child := NewTransaction(
		1,                                        // Nonce
		big.NewFloat(0),                          // Amount
		crypto.AddressFromPrivateKey(privateKey), // Sender
		crypto.AddressFromPrivateKey(privateKey), // Recipient
		[]common.Hash{transaction.Hash},          // Set parent hash
		1,                                        // Gas limit
		big.NewInt(1000),                         // Gas price
		[]byte("test payload"),                   // Payload
	) // Create new child transaction

	err = SignTransaction(child, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = dag.AddTransaction(child) // Add transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	children, err := dag.GetTransactionChildren(transaction.Hash) // Get children

	if err != nil { // Check for errors
		t.Fatal(err) // panic
	}

	if len(children) != 1 { // Check invalid tx set
		t.Fatalf("should have found 1 child transaction; found %d", len(children)) // Log invalid filter
	}

	WorkingDagDB.Close() // Close working dag db

	os.RemoveAll(filepath.FromSlash("data/db/test_network.db")) // Remove existing db
}

// TestGetBestTransaction tests the functionality of the GetBestTransaction() helper method.
func TestGetBestTransaction(t *testing.T) {
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

	dag.Genesis = transaction.Hash // Set genesis hash

	child := NewTransaction(
		1,                                        // Nonce
		big.NewFloat(0),                          // Amount
		crypto.AddressFromPrivateKey(privateKey), // Sender
		crypto.AddressFromPrivateKey(privateKey), // Recipient
		[]common.Hash{transaction.Hash},          // Set parent hash
		1,                                        // Gas limit
		big.NewInt(1000),                         // Gas price
		[]byte("test payload"),                   // Payload
	) // Create new child transaction

	err = SignTransaction(child, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = dag.AddTransaction(child) // Add transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	bestTransaction, err := dag.GetBestTransaction() // Get best transaction

	if err != nil { // Check for errors
		t.Fatal(err) // panic
	}

	if !bytes.Equal(bestTransaction.Hash.Bytes(), child.Hash.Bytes()) { // Check invalid best tx
		t.Fatalf("invalid best transaction; found %s, but wanted %s", hex.EncodeToString(bestTransaction.Hash.Bytes()), child.Hash.Bytes()) // Log invalid best transaction
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
