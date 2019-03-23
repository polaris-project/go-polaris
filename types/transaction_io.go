// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/polaris-project/go-polaris/common"
)

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

// WriteToMemory writes a given transaction to persistent memory in the mempool.
func (transaction *Transaction) WriteToMemory() error {
	err := common.CreateDirIfDoesNotExit(common.MempoolDir) // Create mempool dir if necessary

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = ioutil.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/transaction_%s.json", common.MempoolDir, hex.EncodeToString(transaction.Hash.Bytes()))), transaction.Bytes(), 0644) // Write transaction to persistent memory

	if err != nil { // Check for errors
		return err // Return error
	}

	return nil // No error occurred, return nil
}

// ReadTransactionFromMemory reads a given transaction (specified by hash) from persistent memory.
func ReadTransactionFromMemory(hash common.Hash) (*Transaction, error) {
	data, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/transaction_%s.json", common.MempoolDir, hex.EncodeToString(hash.Bytes())))) // Read transaction

	if err != nil { // Check for errors
		return &Transaction{}, err // Return found error
	}

	buffer := &Transaction{} // Initialize buffer

	err = json.Unmarshal(data, buffer) // Deserialize JSON into buffer.

	if err != nil { // Check for errors
		return &Transaction{}, err // Return found error
	}

	return buffer, nil // No error occurred, return read transaction
}

/* END EXPORTED METHODS */
