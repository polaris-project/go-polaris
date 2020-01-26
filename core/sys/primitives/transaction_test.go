// Package primitives implements a series of basic types required by the network.
package primitives

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/polaris-project/go-polaris/accounts"
	"github.com/polaris-project/go-polaris/crypto"
)

// TestNewTransaction tests the functionality of the NewTransaction helper method.
func TestNewTransaction(t *testing.T) {
	// Generate an account that we'll send the transaction from
	sender, err := accounts.NewAccount()

	// Check for any errors
	if err != nil {
		// Panic
		t.Fatal(err)
	}

	// Generate an account that we'll receive the transcation from
	recipient, err := accounts.NewAccount()

	// Check for any errors
	if err != nil {
		// Panic
		t.Fatal(err)
	}

	// Initialize the test transaction
	tx := NewTransaction(big.NewInt(0), sender.Address(), recipient.Address(), big.NewInt(0), []crypto.Address{}, make(map[crypto.Hash]Receipt), []byte("test"))

	// Make sure that the transaction was actually initialized
	if tx.IsNil() {
		t.Fatal("transaction should not have a zero value")
	}
}

// TestTransactionHash tests the functionality of the transaction Hash helper method.
func TestTransactionHash(t *testing.T) {
	// Generate an account that we'll send the transaction from
	sender, err := accounts.NewAccount()

	// Check for any errors
	if err != nil {
		// Panic
		t.Fatal(err)
	}

	// Generate an account that we'll receive the transcation from
	recipient, err := accounts.NewAccount()

	// Check for any errors
	if err != nil {
		// Panic
		t.Fatal(err)
	}

	// Initialize the test transaction
	tx := NewTransaction(big.NewInt(0), sender.Address(), recipient.Address(), big.NewInt(0), nil, nil, []byte("test"))

	// Calculate the transaction's hash
	txHash := tx.Hash()

	// Make sure that the transaction was actually initialized
	if txHash.IsZero() {
		t.Fatal("transaction should not have a zero value")
	}
}

// TestDeserializeTransaction tests the functionality of the DeserializeTransaction helper method.
func TestDeserializeTransaction(t *testing.T) {
	// Generate an account that we'll send the transaction from
	sender, err := accounts.NewAccount()

	// Check for any errors
	if err != nil {
		// Panic
		t.Fatal(err)
	}

	// Generate an account that we'll receive the transcation from
	recipient, err := accounts.NewAccount()

	// Check for any errors
	if err != nil {
		// Panic
		t.Fatal(err)
	}

	// Initialize the test transaction
	tx := NewTransaction(big.NewInt(0), sender.Address(), recipient.Address(), big.NewInt(0), nil, nil, []byte("test"))

	// Convert the transaction to a slice of bytes so that we can try to deserialize it from them
	txBytes, err := tx.Serialize()

	// Check for errors
	if err != nil {
		// Panic
		t.Fatal(err)
	}

	// Deserialize the transaction
	deserializedTx, err := DeserializeTransaction(txBytes)

	// Check for errors
	if err != nil {
		// Panic
		t.Fatal(err)
	}

	// Calculate the hash of the deserialized transaction
	deserializedTxHash := deserializedTx.Hash()

	// These should be the same transaction, so they should have the same hashes
	if oldHash := tx.Hash(); deserializedTxHash != oldHash {
		// Panic
		t.Fatalf("found %s; expected %s", hex.EncodeToString(deserializedTxHash[:]), hex.EncodeToString(oldHash[:]))
	}
}

// TestTransactionIsNil tests the functionality of the transaction IsNil helper method.
func TestTransactionIsNil(t *testing.T) {
	// Make a zero-value transaction
	tx := NewTransaction(nil, crypto.Hash{}, crypto.Hash{}, nil, nil, nil, nil)

	// Check that the transaction has some fields set to zero or nil
	if !tx.IsNil() {
		t.Fatal("transaction should have been nil")
	}
}
