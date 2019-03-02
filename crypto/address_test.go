// Package crypto provides cryptography helper methods.
package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"testing"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestAddressFromPrivateKey tests the functionality of the AddressFromPrivateKey() helper method.
func TestAddressFromPrivateKey(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate ecdsa private key

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	address := AddressFromPrivateKey(privateKey) // Derive address

	t.Log(hex.EncodeToString(address.Bytes())) // Log success
}

// TestAddressFromPublicKey tests the functionality of the AddressFromPublicKey() helper method.
func TestAddressFromPublicKey(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate ecdsa private key

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	address := AddressFromPublicKey(&privateKey.PublicKey) // Derive address

	t.Log(address) // Log success
}

/* END EXPORTED METHODS TESTS */
