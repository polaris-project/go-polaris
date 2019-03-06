// Package validator represents a helper methods useful for validators in the Polaris network.
// Methods in the validator package are specified in terms of a validator interface, that of which is
// also implemented in the validator package.
package validator

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/polaris-project/go-polaris/config"
	"github.com/polaris-project/go-polaris/crypto"
	"github.com/polaris-project/go-polaris/types"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestNewBeaconDagValidator tests the functionality of the NewBeaconDagValidator() helper method.
func TestNewBeaconDagValidator(t *testing.T) {
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
}

// TestValidateTransaction tests the functionality of the ValidateTransaction() helper method.
func TestValidateTransaction(t *testing.T) {
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

	types.WorkingDagDB.Close() // Close dag db
}

/* END EXPORTED METHODS TESTS */
