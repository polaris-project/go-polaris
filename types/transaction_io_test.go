// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"math/big"
	"testing"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestBytesTransaction tests the functionality of the Bytes() transaction helper method.
func TestBytesTransaction(t *testing.T) {
	transaction := NewTransaction(
		0,                      // Nonce
		nil,                    // Sender
		nil,                    // Recipient
		big.NewInt(10),         // Amount
		1,                      // Gas limit
		big.NewInt(1000),       // Gas price
		[]byte("test payload"), // Payload
	) // Initialize a new transaction using the NewTransaction method

	t.Log(transaction.Bytes()) // Log transaction bytes & test the Bytes() method
}

// TestBytesTransaction tests the functionality of the String() transaction helper method.
func TestStringTransaction(t *testing.T) {
	transaction := NewTransaction(
		0,                      // Nonce
		nil,                    // Sender
		nil,                    // Recipient
		big.NewInt(10),         // Amount
		1,                      // Gas limit
		big.NewInt(1000),       // Gas price
		[]byte("test payload"), // Payload
	) // Initialize a new transaction using the NewTransaction method

	t.Log(transaction.String()) // Log transaction string & test the String() method
}

/* END EXPORTED METHODS TESTS */
