// Package accounts defines a set of ECDSA private-public keypair management utilities and helper methods.
package accounts

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"testing"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestNewAccount tests the functionality of the NewAccount() helper method.
func TestNewAccount(t *testing.T) {
	_, err := NewAccount() // Initialize new account

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

// TestAccountFromKey tests the functionality of the AccountFromKey() helper method.
func TestAccountFromKey(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	account := AccountFromKey(privateKey) // Initialize account from private key

	if account == nil { // Check for nil account
		t.Fatal("account should not be nil") // Panic
	}
}

/* END EXPORTED METHODS TESTS */
