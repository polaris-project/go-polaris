// Package validator represents a helper methods useful for validators in the Polaris network.
// Methods in the validator package are specified in terms of a validator interface, that of which is
// also implemented in the validator package.
package validator

import (
	"bytes"
	"errors"

	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/config"
	"github.com/polaris-project/go-polaris/crypto"
	"github.com/polaris-project/go-polaris/types"
)

var (
	// ErrInvalidTransactionHash is an error definition representing a transcation hash of invalid value.
	ErrInvalidTransactionHash = errors.New("transaction hash is invalid")

	//ErrInvalidTransactionTimestamp is an error definition representing a transaction timestamp of invalid value.
	ErrInvalidTransactionTimestamp = errors.New("invalid transaction timestamp")

	// ErrInvalidTransactionSignature is an error definition representing a transaction signature of invalid value.
	ErrInvalidTransactionSignature = errors.New("invalid transaction signature")

	// ErrInsufficientSenderBalance is an error definition representing a sender balance of insufficient value.
	ErrInsufficientSenderBalance = errors.New("insufficient sender balance")

	// ErrDuplicateTransaction is an error definition representing a transaction of duplicate value in the working dag.
	ErrDuplicateTransaction = errors.New("transcation already exists in the working dag (duplicate)")

	// ErrInvalidTransactionDepth is an error definition representing a transaction of invalid depth value.
	ErrInvalidTransactionDepth = errors.New("invalid transaction depth (not best transaction)")

	// ErrInvalidNonce is an error definition representing a transaction of invalid nonce value.
	ErrInvalidNonce = errors.New("invalid transaction nonce")
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
// Each validation issue is returned as an error.
func (validator *BeaconDagValidator) ValidateTransaction(transaction *types.Transaction) error {
	if !validator.ValidateTransactionHash(transaction) { // Check invalid hash
		return ErrInvalidTransactionHash // Invalid hash
	}

	if !validator.ValidateTransactionTimestamp(transaction) { // Check invalid timestamp
		return ErrInvalidTransactionTimestamp // Invalid timestamp
	}

	if !validator.ValidateTransactionSignature(transaction) { // Check invalid signature
		return ErrInvalidTransactionSignature // Invalid signature
	}

	if !validator.ValidateTransactionSenderBalance(transaction) { // Check invalid value
		return ErrInsufficientSenderBalance // Invalid value
	}

	if !validator.ValidateTransactionIsNotDuplicate(transaction) { // Check duplicate
		return ErrDuplicateTransaction // Duplicate
	}

	if !validator.ValidateTransactionDepth(transaction) { // Check valid depth
		return ErrInvalidTransactionDepth // Invalid depth
	}

	if !validator.ValidateTransactionNonce(transaction) { // Check valid nonce
		return ErrInvalidNonce // Invalid nonce
	}

	return nil // Transaction is valid
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

// ValidateTransactionIsNotDuplicate checks that a given transaction does not already exist in the working dag.
func (validator *BeaconDagValidator) ValidateTransactionIsNotDuplicate(transaction *types.Transaction) bool {
	_, err := validator.WorkingDag.GetTransactionByHash(transaction.Hash) // Attempt to get tx by hash

	if err == nil { // Check transaction exists
		return false // Transaction is duplicate
	}

	return true // Transaction is unique
}

// ValidateTransactionDepth checks that a given transaction's parent hash is a member of the last edge.
func (validator *BeaconDagValidator) ValidateTransactionDepth(transaction *types.Transaction) bool {
	for _, parentHash := range transaction.ParentTransactions { // Iterate through parent hashes
		children, err := validator.WorkingDag.GetTransactionChildren(parentHash) // Get children of transaction

		if err != nil { // Check for errors
			return false // Invalid
		}

		for _, child := range children { // Iterate through children
			currentChildren, err := validator.WorkingDag.GetTransactionChildren(child.Hash) // Get children of current child

			if err != nil { // Check for errors
				return false // Invalid
			}

			if len(currentChildren) != 0 { // Check child has children
				return false // Invalid depth
			}
		}
	}

	return true // Valid
}

// ValidateTransactionNonce checks that a given transaction's nonce is equivalent to the sending account's last nonce + 1.
func (validator *BeaconDagValidator) ValidateTransactionNonce(transaction *types.Transaction) bool {
	senderTransactions, err := validator.WorkingDag.GetTransactionsBySender(transaction.Sender) // Get sender txs

	if err != nil { // Check for errors
		return false // Invalid
	}

	if len(senderTransactions) == 0 { // Check is genesis
		if transaction.AccountNonce != 0 { // Check nonce is not 0
			return false // Invalid nonce
		}

		return true // Valid nonce
	}

	lastNonce := uint64(0) // Init nonce buffer

	for _, currentTransaction := range senderTransactions { // Iterate through sender txs
		if currentTransaction.AccountNonce > lastNonce { // Check greater than last nonce
			lastNonce = currentTransaction.AccountNonce // Set last nonce
		}
	}

	if transaction.AccountNonce != lastNonce+1 { // Check invalid nonce
		return false // Invalid nonce
	}

	return true // Valid nonce
}

/* END EXPORTED METHODS */
