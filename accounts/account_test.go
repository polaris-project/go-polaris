// Package accounts implements types and methods for aiding in the generation and serialization of polaris accounts.
package accounts

import "testing"

// TestNewAccount tests the functionality of the NewAccount helper method.
func TestNewAccount(t *testing.T) {
	// Generate the account
	_, err := NewAccount()

	// Check for any errors that were generated whilst creating the account
	if err != nil {
		// Print the error
		t.Fatal(err)
	}
}

// TestAccountAddress tests the functionality of the account Address helper method.
func TestAccountAddress(t *testing.T) {
	// Generate the account
	acc, err := NewAccount()

	// Check for any errors that arose whilst creating the account
	if err != nil {
		// Panic
		t.Fatal(err)
	}

	// Get the account's address
	address := acc.Address()

	// Make sure that an actual address was generated, not just fluff
	if address.IsZero() {
		// Panic
		t.Fatal("should have generated a real address; found nil")
	}
}
