// Package common defines a set of commonly used helper methods and data types.
package common

import (
	"testing"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestNewHash tests the functionality of the NewHash() helper method.
func TestNewHash(t *testing.T) {
	b := []byte("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80") // Hash test

	t.Log(NewHash(b)) // Log hash value
}

// TestIsNil tests the functionality of the IsNil() helper method.
func TestIsNil(t *testing.T) {
	hash := NewHash([]byte("36f028580bb02cc8272a9a020f4200e346e276ae664e45ee80745574e2f5ab80")) // Hash test
	nilHash := NewHash([]byte(""))                                                              // Hash nil

	if hash.IsNil() == true { // Check invalid determination
		t.Fatal("invalid determination: hash isn't nil") // Panic
	}

	if nilHash.IsNil() == false { // Check invalid determination
		t.Fatal("invalid determination: hash should be nil") // Panic
	}
}

// TestNewAddress tests the functionality of the NewAddress() helper method.
func TestNewAddress(t *testing.T) {
	address := NewAddress([]byte("99994467913f1743a5b2ca672fcfb4b7dd268d3f5d36f5e61079c8a2706c9112")) // Create new address

	t.Log(address) // Log address
}

/* END EXPORTED METHODS TESTS */
