// Package primitives implements a series of basic types required by the network.
package primitives

import (
	"math/big"
	"testing"

	"github.com/polaris-project/go-polaris/accounts"
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
	tx := NewTransaction(big.NewInt(0), sender.Address(), recipient.Address(), big.NewInt(0), []byte("test"))

	// Make sure that the transaction was actually initialized
	if tx.IsZero() {
		t.Fatal("transaction should not have a zero value")
	}
}
