// Package accounts defines a set of ECDSA private-public keypair management utilities and helper methods.
package accounts

import "testing"

/* BEGIN EXPORTED METHODS */

// TestWriteToMemory tests the functionality of the WriteToMemory() helper method.
func TestWriteToMemory(t *testing.T) {
	account, err := NewAccount() // Initialize new account

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = account.WriteToMemory() // Write account to persistent memory

	if err != nil { // Check for errors
		t.Fatal(err) // panic
	}
}

// TestReadAccountFromMemory tests the functionality of the ReadAccountFromMemory() helper method.
func TestReadAccountFromMemory(t *testing.T) {
	account, err := NewAccount() // Initialize new account

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = account.WriteToMemory() // Write account to persistent memory

	if err != nil { // Check for errors
		t.Fatal(err) // panic
	}

	_, err = ReadAccountFromMemory(account.Address()) // Read account from memory

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

/* END EXPORTED METHODS */
