// Package primitives implements a series of basic types required by the network.
package primitives

import (
	"math/big"

	"github.com/polaris-project/go-polaris/crypto"
	"github.com/tendermint/go-amino"
)

// Transaction represents an atomic transfer between two users of the polaris network.
type Transaction struct {
	// Nonce is the index of the transaction in the sender's history of transactions
	Nonce *big.Int

	// Sender is the address of the sender of the transaction
	Sender crypto.Address

	// Recipient is the address of the recipient of the transaction
	Recipient crypto.Address

	// Value is the value of the transaction, or the number of pulsars sent
	Value *big.Int

	// Parents is a list of hashes representing each of the transactions that the transaction was spawned off of
	Parents []crypto.Hash

	// Receipts represents the set of transaction outputs for each of the transaction's parents
	Receipts map[crypto.Hash]Receipt

	// Payload is any miscellaneous data sent along with the transaction
	Payload []byte
}

// NewTransaction initializes a new Transaction with the given nonce, sender, recipient, value, and payload.
func NewTransaction(nonce *big.Int, sender crypto.Address, recipient crypto.Address, value *big.Int, parents []crypto.Hash, receipts map[crypto.Hash]Receipt, payload []byte) Transaction {
	// Initialize and return the transaction
	return Transaction{
		nonce,
		sender,
		recipient,
		value,
		parents,
		receipts,
		payload,
	}
}

// Hash returns a hash of a binary-serialized version of the transaction.
func (t *Transaction) Hash() crypto.Hash {
	// Serialize the transaction as a slice of bytes so we can hash it
	bytes, err := t.Serialize()

	// If we weren't able to serialize the transaction, return an empty hash
	if err != nil {
		return crypto.Hash{}
	}

	// Return the hash
	return crypto.HashBlake3(bytes)
}

// DeserializeTransaction deserializes a transaction from the given slice of bytes.
func DeserializeTransaction(b []byte) (Transaction, error) {
	// Declare a tx buffer that we'll deserialize the bytes into
	var tx Transaction

	// Return the deserialized transaction, as well as any errors
	return tx, amino.UnmarshalBinaryBare(b, &tx)
}

// Serialize serializes the transaction using the amino binary format.
func (t *Transaction) Serialize() ([]byte, error) {
	// Serialize the transaction, and return any errors that occurred
	return amino.MarshalBinaryBare(t)
}

// IsNil checks if the transaction's fields have been initialized to zero. Note: the payload field is omitted in this calculation.
func (t *Transaction) IsNil() bool {
	return t.Nonce == nil || t.Sender.IsZero() || t.Recipient.IsZero() || t.Value == nil || t.Parents == nil || t.Receipts == nil
}
