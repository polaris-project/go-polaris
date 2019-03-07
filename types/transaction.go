// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"errors"
	"math/big"
	"time"

	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/crypto"
)

var (
	// ErrNoWorkingHost represents an error describing a p2p.WorkingHost value of nil.
	ErrNoWorkingHost = errors.New("no valid global host was found")
)

// Transaction is a data type representing a transfer of monetary value between addresses.
// A transactions does not necessarily imply the transfer of value between human peers, but also contracts.
type Transaction struct {
	AccountNonce uint64 `json:"nonce" gencodec:"required"` // Index in account transaction list

	Amount *big.Float `json:"amount" gencodec:"required"` // Transaction value

	Sender    *common.Address `json:"sender" gencodec:"required"`    // Transaction sender
	Recipient *common.Address `json:"recipient" gencodec:"required"` // Transaction recipient

	ParentTransactions []common.Hash `json:"parent" gencodec:"required"` // Parent hash

	GasPrice *big.Int `json:"gas_price" gencodec:"required"` // Gas price in units equivalent to 0.000000001 of a single unit
	GasLimit uint64   `json:"gas_limit" gencodec:"required"` // Value of gas price willing to pay for transaction

	Payload []byte `json:"payload" gencodec:"required"` // Data sent with transaction (i.e. contract bytecode, message, etc...)

	Signature *Signature `json:"signature" gencodec:"required"` // ECDSA transaction signature

	Timestamp time.Time `json:"timestamp" gencodec:"required"` // Transaction timestamp

	Hash common.Hash `json:"hash" gencodec:"required"` // Transaction hash
}

/* BEGIN EXPORTED METHODS */

// NewTransaction creates a new transaction with the given account nonce, value, sender, recipient, gas price, gas limit, and payload.
func NewTransaction(accountNonce uint64, amount *big.Float, sender, recipient *common.Address, parentTransactions []common.Hash, gasLimit uint64, gasPrice *big.Int, payload []byte) *Transaction {
	transaction := &Transaction{
		AccountNonce:       accountNonce,       // Set account nonce
		Amount:             amount,             // Set amount
		Sender:             sender,             // Set sender
		Recipient:          recipient,          // Set recipient
		ParentTransactions: parentTransactions, // Set parents
		GasPrice:           gasPrice,           // Set gas price
		GasLimit:           gasLimit,           // Set gas limit
		Payload:            payload,            // Set payload
		Signature:          nil,                // Set signature
		Timestamp:          time.Now().UTC(),   // Set timestamp
	}

	(*transaction).Hash = crypto.Sha3(transaction.Bytes()) // Set transaction hash

	return transaction // Return initialized transaction
}

// CalculateTotalValue calculates the total value of a transaction, including both its amount and total gas.
func (transaction *Transaction) CalculateTotalValue() *big.Float {
	return new(big.Float).Add(transaction.Amount, new(big.Float).SetInt(transaction.GasPrice.Mul(transaction.GasPrice, big.NewInt(int64(transaction.GasLimit))))) // Return total value
}

/* BEGIN EXPORTED METHODS */
