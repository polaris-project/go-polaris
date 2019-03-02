// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"math/big"
	"testing"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestNewTransactions tests the functionality of the NewTransaction method.
func TestNewTransactions(t *testing.T) {
	transaction := NewTransaction(
		0,                      // Nonce
		big.NewFloat(10),       // Amount
		nil,                    // Sender
		nil,                    // Recipient
		nil,                    // Parents
		1,                      // Gas limit
		big.NewInt(1000),       // Gas price
		[]byte("test payload"), // Payload
	) // Create a new transaction using the NewTransaction method

	t.Log(transaction) // Log the initialized transaction
}

/* END EXPORTED METHODS TESTS */
