// Package common defines a set of commonly used helper methods and data types.
package common

import (
	"testing"

	"github.com/polaris-project/go-polaris/crypto"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestNewHash tests the functionality of the NewHash helper method.
func TestNewHash(t *testing.T) {
	b := crypto.Sha3([]byte("test")) // Hash test

	t.Log(NewHash(b)) // Log hash value
}

/* END EXPORTED METHODS TESTS */
