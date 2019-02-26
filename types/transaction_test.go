package types

import (
	"math/big"
	"testing"
)

func TestNewTransactions(t *testing.T) {

	// Create a new transaction using the NewTransaction method
	transaction := NewTransaction(
		0,                      // Nonce
		big.NewInt(10),         // Amount
		nil,                    // Sender
		nil,                    // Recipient
		1,                      // Gas limit
		big.NewInt(1000),       // Gas price
		[]byte("test payload"), // Payload
	)

	t.Log(transaction) // Log the serialized transaction

}
