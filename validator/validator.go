// Package validator represents a helper methods useful for validators in the Polaris network.
// Methods in the validator package are specified in terms of a validator interface, that of which is
// also implemented in the validator package.
package validator

import "github.com/polaris-project/go-polaris/types"

// Validator represents any generic validator.
type Validator interface {
	ValidateTransaction(transaction *types.Transaction) bool // Validate a given transaction

	ValidateTransactionHash(transaction *types.Transaction) bool // Validate a given transaction's hash

	ValidateTransactionTimestamp(transaction *types.Transaction) bool // Validate a given transaction's timestamp

	ValidateTransactionSignature(transaction *types.Transaction) bool // Validate a given transaction's signature

	ValidateTransactionSenderBalance(transaction *types.Transaction) bool // Validate a given transaction's sender has

	ValidateTransactionIsNotDuplicate(transaction *types.Transaction) bool // Validate that a given transaction does not already exist in the working dag

	ValidationProtocol() string // Get the current validator's validation protocol
}
