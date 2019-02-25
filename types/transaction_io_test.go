package types

import (
	"math/big"
	"testing"

	"github.com/polaris-project/go-polaris/crypto"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestBytesTransaction tests the functionality of the Bytes() transaction helper method.
func TestBytesTransaction(t *testing.T) {
	transaction := &Transaction{
		AccountNonce: 0,                // Set nonce
		Sender:       nil,              // Set sender
		Recipient:    nil,              // Set recipient
		GasPrice:     big.NewInt(1000), // Set gas price
		Payload:      []byte("test"),   // Set payload
		Signature:    nil,              // Set signature
	}

	transaction.Hash = crypto.Sha3(transaction.Bytes()) // Hash transaction

	t.Log(transaction.Bytes()) // Log transaction bytes
}

/* END EXPORTED METHODS TESTS */
