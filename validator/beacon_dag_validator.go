// Package validator represents a helper methods useful for validators in the Polaris network.
// Methods in the validator package are specified in terms of a validator interface, that of which is
// also implemented in the validator package.
package validator

import (
	"bytes"

	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/config"
	"github.com/polaris-project/go-polaris/crypto"
	"github.com/polaris-project/go-polaris/types"
)

// BeaconDagValidator represents a main dag validator.
type BeaconDagValidator struct {
	Config *config.DagConfig `json:"config"` // Config represents the beacon dag config

	WorkingDag *types.Dag `json:"dag"` // Working validator dag
}

/* BEGIN EXPORTED METHODS */

// NewBeaconDagValidator initializes a new beacon dag with a given config and working dag.
func NewBeaconDagValidator(config *config.DagConfig, workingDag *types.Dag) *BeaconDagValidator {
	return &BeaconDagValidator{
		Config:     config,     // Set config
		WorkingDag: workingDag, // Set working dag
	}
}

// ValidateTransaction validates the given transaction, transaction via the standard beacon dag validator.
func (validator *BeaconDagValidator) ValidateTransaction(transaction *types.Transaction) bool {
	if !validator.ValidateTransactionHash(transaction) { // Check invalid hash
		return false // Invalid hash
	}

	if !validator.ValidateTransactionTimestamp(transaction) { // Check invalid timestamp
		return false // Invalid timestamp
	}

	if !validator.ValidateTransactionSignature(transaction) { // Check invalid signature
		return false // Invalid signature
	}

	if !validator.ValidateTransactionSenderBalance(transaction) { // Check invalid value
		return false // Invalid value
	}

	return true // Transaction is valid
}

// ValidateTransactionHash checks that a given transaction's hash is equivalent to the calculated hash of that given transaction.
func (validator *BeaconDagValidator) ValidateTransactionHash(transaction *types.Transaction) bool {
	if transaction.Hash.IsNil() { // Check transaction doesn't have transaction
		return false // No valid hash
	}

	unsignedTx := *transaction // Get unsigned

	unsignedTx.Hash = common.NewHash(nil) // Set hash to nil

	unsignedTx.Signature = transaction.Signature // Reset signature

	return bytes.Equal(transaction.Hash.Bytes(), crypto.Sha3(unsignedTx.Bytes()).Bytes()) // Return hashes equivalent
}

// ValidateTransactionTimestamp validates the given transaction's timestamp against that of its parents.
// If the timestamp of any one of the given transaction's parents is after the given transaction's timestamp, false is returned.
// If any one of the transaction's parent transactions cannot be found in the working dag, false is returned.
func (validator *BeaconDagValidator) ValidateTransactionTimestamp(transaction *types.Transaction) bool {
	for _, parentHash := range transaction.ParentTransactions { // Iterate through parent hashes
		parentTransaction, err := validator.WorkingDag.GetTransactionByHash(parentHash) // Get parent transaction pointer

		if err != nil { // Check for errors
			return false // Invalid parent
		}

		if parentTransaction.Timestamp.After(transaction.Timestamp) {
			return false // Invalid timestamp
		}
	}

	return true // Valid timestamp
}

// ValidateTransactionSignature validates the given transaction's signature against the transaction sender's public key.
// If the transaction's signature is nil, false is returned.
func (validator *BeaconDagValidator) ValidateTransactionSignature(transaction *types.Transaction) bool {
	if transaction.Signature == nil { // Check has no signature
		return false // Nil signature
	}

	return transaction.Signature.Verify(transaction.Sender) // Return signature validity
}

// ValidateTransactionSenderBalance checks that a given transaction's sender has a balance greater than or equal to the transaction's total value (including gas costs).
func (validator *BeaconDagValidator) ValidateTransactionSenderBalance(transaction *types.Transaction) bool {
	balance, err := validator.WorkingDag.CalculateAddressBalance(transaction.Sender) // Calculate balance

	if err != nil { // Check for errors
		return false // Invalid
	}

	return balance.Cmp(transaction.CalculateTotalValue()) == 0 || balance.Cmp(transaction.CalculateTotalValue()) == 1 // Return sender balance adequate
}

/* END EXPORTED METHODS */
