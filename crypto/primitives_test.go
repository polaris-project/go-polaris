// Package crypto defines a set of helper methods for repetitive cryptographic operations used by the Polaris protocol.
package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"testing"

	"github.com/icrowley/fake"
	"lukechampine.com/blake3"
)

// validTestHashStringRepresentation is the only valid representation of the data hashed in the TestHashString method below.
const validTestHashStringRepresentation = "c0d8220bcf2330d7d265169c8d083a8f67c233e6b72d2a8f6b58967fe93c1179"

/* BEGIN TESTING FOR THE ADDRESS TYPE */

// TestAddressFromPrivateKey tests the functionality of the AddressFromPrivateKey helper method.
func TestAddressFromPrivateKey(t *testing.T) {
	// Generate a private key for testing
	_, pk, err := ed25519.GenerateKey(rand.Reader)

	// Check for any errors that occurred whilst generating the private key
	if err != nil {
		// Panic
		t.Fatal(err)
	}

	// Generate an address from the private key
	address := AddressFromPrivateKey(&pk)

	// Ensure that the generated address isn't a zero-value address.
	if address.IsZero() {
		// Panic with the given error
		t.Fatal(err)
	}
}

/* END TESTING FOR THE ADDRESS TYPE */

/* BEGIN TESTING FOR THE HASH TYPE */

// TestHashBlake3 tests the functionality of the HashBlake3 helper method.
func TestHashBlake3(t *testing.T) {
	// Generate a random word to test the hashing helper method on
	data := []byte(fake.Word())

	// Hash the generated data
	found := HashBlake3(data)

	// This shouldn't ever fail, unless we've changed the hashing algorithm. Should the hashing algorithm be changed, this failing test
	// will alert us that we need to make this same change in MANY other files, not just this one.
	if expected := blake3.Sum256(data); found != expected {
		// Make sure this test fails with proper reason
		t.Fatalf("found %s; expected %s", hex.EncodeToString(found[:]), hex.EncodeToString(expected[:]))
	}
}

// TestHashFromString tests the functionality of the provided string to hash conversion helper method.
func TestHashFromString(t *testing.T) {
	// Generate a hash of some some fake word to test our conversion helper method on
	data := HashBlake3([]byte(fake.Word()))

	// Try to convert the hash into a string, and back into a hash. This should always succeed.
	if _, err := HashFromString(data.String()); err != nil {
		// Log the error
		t.Fatal(err)
	}

	// Get the hexadecimal representation of the hash, but exclude the last item. Conversion should fail.
	invalidData := data.String()[:31]

	// Try to convert the invalid hash string back into a valid hash. This should ALWAYS fail.
	if _, err := HashFromString(invalidData); err == nil {
		// Panic
		t.Fatalf("%s shouldn't be able to be converted into a valid hash", invalidData)
	}
}

// TestHashString tests the functionality of the Hash String conversion helper method.
func TestHashString(t *testing.T) {
	// Generate a hash of some some fake word to test our conversion helper method on
	data := HashBlake3([]byte("some test data"))

	// Make sure that the provided hash string representation matches the pre-determined valid representation string
	if data.String() != validTestHashStringRepresentation {
		// Panic
		t.Fatalf("found %s, expected %s", data.String(), validTestHashStringRepresentation)
	}
}

/* END TESTING FOR THE HASH TYPE */
