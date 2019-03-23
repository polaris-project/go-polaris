// Package transaction represents the transaction RPC server.
package transaction

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/polaris-project/go-polaris/accounts"
	"github.com/polaris-project/go-polaris/common"
	transactionProto "github.com/polaris-project/go-polaris/internal/proto/transaction"
	"github.com/polaris-project/go-polaris/p2p"
	"github.com/polaris-project/go-polaris/types"
)

var (
	// ErrNilHashRequest defines an error describing a TransactionHash length of 0.
	ErrNilHashRequest = errors.New("request did not contain a valid transaction hash")

	// ErrInvalidHashRequest defines an error describing a MessageHash length >||< common.HashLength.
	ErrInvalidHashRequest = errors.New("request did not contain a valid message hash")
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

	recipientBytes, err := hex.DecodeString(request.Address2) // Decode recipient address hex-encoded string value

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

	account, err := accounts.ReadAccountFromMemory(common.NewAddress(transaction.Sender.Bytes())) // Open account

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	err = types.SignTransaction(transaction, account.PrivateKey()) // Sign transaction

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	return &transactionProto.GeneralResponse{Message: transaction.Signature.String()}, nil // Return signature
}

// Publish handles the Publish request method.
func (server *Server) Publish(ctx context.Context, request *transactionProto.GeneralRequest) (*transactionProto.GeneralResponse, error) {
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

	publishContext, cancel := context.WithCancel(ctx) // Get context

	defer cancel() // Cancel

	err = p2p.WorkingClient.PublishTransaction(publishContext, transaction) // Publish transaction

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	return &transactionProto.GeneralResponse{Message: fmt.Sprintf("published transaction %s successfully", hex.EncodeToString(transaction.Hash.Bytes()))}, nil // Return success
}

// SignMessage handles the SignMessage request method.
func (server *Server) SignMessage(ctx context.Context, request *transactionProto.GeneralRequest) (*transactionProto.GeneralResponse, error) {
	senderBytes, err := hex.DecodeString(request.Address) // Decode sender address hex-encoded string value

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	account, err := accounts.ReadAccountFromMemory(common.NewAddress(senderBytes)) // Open account

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	if len(request.Payload) != common.HashLength { // Check invalid message hash
		return &transactionProto.GeneralResponse{}, ErrInvalidHashRequest
	}

	signature, err := types.SignMessage(common.NewHash(request.Payload), account.PrivateKey()) // Sign message

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	return &transactionProto.GeneralResponse{Message: signature.String()}, nil // Retrun signature
}

// Verify handles the Verify request method.
func (server *Server) Verify(ctx context.Context, request *transactionProto.GeneralRequest) (*transactionProto.GeneralResponse, error) {
	transactionHashBytes, err := hex.DecodeString(request.TransactionHash[0]) // Get transaction hash byte value

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	transaction, err := types.ReadTransactionFromMemory(common.NewHash(transactionHashBytes)) // Read transaction

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	return &transactionProto.GeneralResponse{Message: strconv.FormatBool(transaction.Signature.Verify(transaction.Sender))}, nil // Return is valid
}

// String handles the String request method.
func (server *Server) String(ctx context.Context, request *transactionProto.GeneralRequest) (*transactionProto.GeneralResponse, error) {
	transactionHashBytes, err := hex.DecodeString(request.TransactionHash[0]) // Get transaction hash byte value

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	transaction, err := types.ReadTransactionFromMemory(common.NewHash(transactionHashBytes)) // Read transaction

	if err != nil { // Check for errors
		return &transactionProto.GeneralResponse{}, err // Return found error
	}

	return &transactionProto.GeneralResponse{Message: transaction.String()}, nil // Return tx string value
}

/* END EXPORTED METHODS */
