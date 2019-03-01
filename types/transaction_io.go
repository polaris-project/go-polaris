// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import "encoding/json"

/* BEGIN EXPORTED METHODS */

// TransactionFromBytes deserializes a transaction from a given byte array.
func TransactionFromBytes(b []byte) *Transaction {
	buffer := &Transaction{} // Initialize tx buffer

	err := json.Unmarshal(b, buffer) // Unmarshal

	if err != nil { // Check for errors
		return &Transaction{}
	}

	return buffer // Return deserialized transaction
}

// Bytes serializes a given transaction to a byte array via json.
func (transaction *Transaction) Bytes() []byte {
	marshaledVal, _ := json.MarshalIndent(*transaction, "", "  ") // Marshal JSON

	return marshaledVal // Return marshaled value
}

// String serializes a given transaction to a string via json.
func (transaction *Transaction) String() string {
	marshaledVal, _ := json.MarshalIndent(*transaction, "", "  ") // Marshal JSON

	return string(marshaledVal) // Returned the marshalled JSON as a string
}

/* END EXPORTED METHODS */
