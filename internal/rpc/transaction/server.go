// Package transaction represents the transaction RPC server.
package transaction

import (
	"context"
	"encoding/hex"
	"errors"
	"math/big"

	"github.com/polaris-project/go-polaris/common"
	transactionProto "github.com/polaris-project/go-polaris/internal/proto/transaction"
	"github.com/polaris-project/go-polaris/types"
)

var (
	// ErrNilHashRequest defines an error describing a TransactionHash length of 0.
	ErrNilHashRequest = errors.New("request did not contain a valid transaction hash")
)

// Server represents a Polaris RPC server.
type Server struct{}

/* BEGIN EXPORTED METHODS */

// NewTransaction handles the NewTransaction request method.
func (server *Server) NewTransaction(ctx context.Context, request *transactionProto.GeneralRequest) (*transactionProto.GeneralResponse, error) {
	senderBytes, err := hex.DecodeString(request.Address) // Decode sender address hex-encoded string value

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	recipientBytes, err := hex.DecodeString(request.Address) // Decode recipient address hex-encoded string value

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	amount, _, err := big.ParseFloat(string(request.Amount), 10, 18, big.ToNearestEven) // Parse amount value

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	var parentHashes []common.Hash // Init parent hash buffer

	for _, parentHashString := range request.TransactionHash { // Iterate through parent hashes
		parentHashBytes, err := hex.DecodeString(parentHashString) // Decode hash string value

		if err != nil { // Check for errors
			return &transactionProto.GeneralResponse{}, err // Return found error
		}

		parentHashes = append(parentHashes, common.NewHash(parentHashBytes)) // Append hash
	}

	transaction := types.NewTransaction(request.Nonce, amount, common.NewAddress(senderBytes), common.NewAddress(recipientBytes), parentHashes, request.GasLimit, new(big.Int).SetBytes(request.GasPrice), request.Payload) // Initialize transaction

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	err = transaction.WriteToMemory() // Write transaction to mempool

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	return &transactionProto.GeneralResponse{Message: hex.EncodeToString(transaction.Hash.Bytes())}, nil // Return transaction hash string value
}

// CalculateTotalValue handles the CalculateTotalValue request method.
func (server *Server) CalculateTotalValue(ctx context.Context, request *transactionProto.GeneralRequest) (*transactionProto.GeneralResponse, error) {
	if len(request.TransactionHash) == 0 { // Check nothing to read
		return &transactionProto.GeneralResponse{}, ErrNilHashRequest // Return error
	}

	transactionHashBytes, err := hex.DecodeString(request.TransactionHash[0]) // Get transaction hash byte value

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	transaction, err := types.ReadTransactionFromMemory(common.NewHash(transactionHashBytes)) // Read transaction

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	return &transactionProto.GeneralResponse{Message: transaction.CalculateTotalValue().String()}, nil // Return total value
}

// SignTransaction handles the SignTransaction request method.
func (server *Server) SignTransaction(ctx context.Context, request *transactionProto.GeneralRequest) (*transactionProto.GeneralResponse, error) {
	if len(request.TransactionHash) == 0 { // Check nothing to read
		return &transactionProto.GeneralResponse{}, ErrNilHashRequest // Return error
	}

	transactionHashBytes, err := hex.DecodeString(request.TransactionHash[0]) // Get transaction hash byte value

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	transaction, err := types.ReadTransactionFromMemory(common.NewHash(transactionHashBytes)) // Read transaction

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	err = types.SignTransaction(transaction)
}

/* END EXPORTED METHODS */
