// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"bytes"
	"math/big"
	"testing"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestTransactionFromBytes tests the functionality of the TransactionFromBytes() transaction helper method.
func TestTransactionFromBytes(t *testing.T) {
	transaction := NewTransaction(
		0,                      // Nonce
		big.NewFloat(10),       // Amount
		nil,                    // Sender
		nil,                    // Recipient
		nil,                    // Parents
		1,                      // Gas limit
		big.NewInt(1000),       // Gas price
		[]byte("test payload"), // Payload
	) // Initialize a new transaction using the NewTransaction method

	if !bytes.Equal(transaction.Bytes(), TransactionFromBytes(transaction.Bytes()).Bytes()) { // Check transactions not equal
		t.Fatal("deserialized transaction should be equivalent to source") // Panic
	}
}

// TestBytesTransaction tests the functionality of the Bytes() transaction helper method.
func TestBytesTransaction(t *testing.T) {
	transaction := NewTransaction(
		0,                      // Nonce
		big.NewFloat(10),       // Amount
		nil,                    // Sender
		nil,                    // Recipient
		nil,                    // Parents
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
		big.NewFloat(10),       // Amount
		nil,                    // Sender
		nil,                    // Recipient
		nil,                    // Parents
		1,                      // Gas limit
		big.NewInt(1000),       // Gas price
		[]byte("test payload"), // Payload
	) // Initialize a new transaction using the NewTransaction method

	t.Log(transaction.String()) // Log transaction string & test the String() method
}

/* END EXPORTED METHODS TESTS */
