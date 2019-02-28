// Package crypto provides cryptography helper methods.
package crypto

import "testing"

/* BEGIN EXPORTED METHODS TESTS */

// TestSha3 tests the functionality of the sha3() helper method.
func TestSha3(t *testing.T) {
	if Sha3([]byte("test")).IsNil() { // Check nil hash
		t.Fatal("hash should not be nil") // Panic
	}
}

// TestSha3String tests the functionality of the Sha3String() helper method.
func TestSha3String(t *testing.T) {
	if Sha3String([]byte("test")) == "" { // Check nil hash
		t.Fatal("hash should not be nil") // Panic
	}
}

// TestSha3n tests the functionality of the Sha3n() helper method.
func TestSha3n(t *testing.T) {
	if Sha3n([]byte("test"), 2).IsNil() { // Check nil hash
		t.Fatal("hash should not be nil") // Panic
	}
}

// TestSha3nString tests the functionality of the Sha3nString() helper method.
func TestSha3nString(t *testing.T) {
	if Sha3nString([]byte("test"), 2) == "" { // Check nil hash
		t.Fatal("hash should not be nil") // Panic
	}
}

// TestSha3d tests the functionality of the Sha3d() helper method.
func TestSha3d(t *testing.T) {
	if Sha3d([]byte("test")).IsNil() { // Check nil hash
		t.Fatal("hash should not be nil") // Panic
	}
}

// TestSha3dString tests the functionality of the Sha3dString() helper method.
func TestSha3dString(t *testing.T) {
	if Sha3dString([]byte("test")) == "" { // Check nil hash
		t.Fatal("hash should not be nil") // Panic
	}
}

/* END EXPORTED METHODS TESTS */
