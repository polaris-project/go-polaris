// Package validator represents a helper methods useful for validators in the Polaris network.
// Methods in the validator package are specified in terms of a validator interface, that of which is
// also implemented in the validator package.
package validator

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/config"
	"github.com/polaris-project/go-polaris/crypto"
	"github.com/polaris-project/go-polaris/types"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestNewBeaconDagValidator tests the functionality of the NewBeaconDagValidator() helper method.
func TestNewBeaconDagValidator(t *testing.T) {
	os.RemoveAll(filepath.FromSlash("data/db/test_network.db")) // Remove existing db

	dagConfig := config.NewDagConfig(nil, "test_network", 1) // Initialize new dag config with test genesis file.

	dag, err := types.NewDag(dagConfig) // Initialize dag with dag config

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	validator := NewBeaconDagValidator(dagConfig, dag) // Initialize validator

	if validator == nil { // Check validator is nil
		t.Fatal("validator should not be nil") // Panic
	}

	types.WorkingDagDB.Close() // Close dag db

	os.RemoveAll(filepath.FromSlash("data/db/test_network.db")) // Remove existing db
}

// TestValidateTransaction tests the functionality of the ValidateTransaction() helper method.
func TestValidateTransaction(t *testing.T) {
	os.RemoveAll(filepath.FromSlash("data/db/test_network.db")) // Remove existing db

	dagConfig := config.NewDagConfig(nil, "test_network", 1) // Initialize new dag config with test genesis file.

	dag, err := types.NewDag(dagConfig) // Initialize dag with dag config

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	address := crypto.AddressFromPrivateKey(privateKey) // Generate address

	transaction := types.NewTransaction(
		0,                      // Nonce
		big.NewFloat(0),        // Amount
		address,                // Sender
		nil,                    // Recipient
		nil,                    // Parents
		0,                      // Gas limit
		big.NewInt(0),          // Gas price
		[]byte("test payload"), // Payload
	) // Initialize a new transaction

	err = types.SignTransaction(transaction, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	validator := NewBeaconDagValidator(dagConfig, dag) // Initialize validator

	if validator == nil { // Check validator is nil
		t.Fatal("validator should not be nil") // Panic
	}

	if err := validator.ValidateTransaction(transaction); err != nil { // Validate
		t.Fatalf("tx should be valid; got %s error", err.Error()) // Panic
	}

	err = dag.AddTransaction(transaction) // Add transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	child := types.NewTransaction(
		1,                               // Nonce
		big.NewFloat(0),                 // Amount
		address,                         // Sender
		nil,                             // Recipient
		[]common.Hash{transaction.Hash}, // Parents
		0,                               // Gas limit
		big.NewInt(0),                   // Gas price
		[]byte("test payload"),          // Payload
	) // Create child transaction

	sibling := types.NewTransaction(
		2,                               // Nonce
		big.NewFloat(0),                 // Amount
		address,                         // Sender
		nil,                             // Recipient
		[]common.Hash{transaction.Hash}, // Parents
		0,                               // Gas limit
		big.NewInt(0),                   // Gas price
		[]byte("test payload"),          // Payload
	) // Create child transaction

	err = types.SignTransaction(child, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = types.SignTransaction(sibling, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	if err := validator.ValidateTransaction(child); err != nil { // Validate
		t.Fatalf("tx should be valid; got %s error", err.Error()) // Panic
	}

	err = dag.AddTransaction(child) // Add child transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	if err := validator.ValidateTransaction(sibling); err != nil { // Validate
		t.Fatalf("tx should be valid; got %s error", err.Error()) // Panic
	}

	types.WorkingDagDB.Close() // Close dag db

	os.RemoveAll(filepath.FromSlash("data/db/test_network.db")) // Remove existing db
}

/* END EXPORTED METHODS TESTS */
