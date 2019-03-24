// Package dag represents the dag RPC server.
package dag

import (
	"context"
	"encoding/hex"
	"strings"

	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/config"
	dagProto "github.com/polaris-project/go-polaris/internal/proto/dag"
	"github.com/polaris-project/go-polaris/types"
)

// Server represents a Polaris RPC server.
type Server struct{}

/* BEGIN EXPORTED METHODS */

// NewDag handles the NewDag request method.
func (server *Server) NewDag(ctx context.Context, request *dagProto.GeneralRequest) (*dagProto.GeneralResponse, error) {
	config, err := config.ReadDagConfigFromMemory(request.Network) // Read config

	if err != nil { // Check for errors
		return &dagProto.GeneralResponse{}, err // Return found error
	}

	dag, err := types.NewDag(config) // Initialize dag

	if err != nil { // Check for errors
		return &dagProto.GeneralResponse{}, err // Return found error
	}

	err = dag.WriteToMemory() // Write dag to persistent memory

	if err != nil { // Check for errors
		return &dagProto.GeneralResponse{}, err // Return found error
	}

	dag.Close() // Close dag

	return &dagProto.GeneralResponse{Message: string(dag.Bytes())}, nil // Return dag db header string
}

// MakeGenesis handles the MakeGenesis request method.
func (server *Server) MakeGenesis(ctx context.Context, request *dagProto.GeneralRequest) (*dagProto.GeneralResponse, error) {
	dag, err := types.OpenDag(request.Network) // Open dag

	if err != nil { // Check for errors
		return &dagProto.GeneralResponse{}, err // Return found error
	}

	genesisTransactions, err := dag.MakeGenesis() // Make genesis

	if err != nil { // Check for errors
		return &dagProto.GeneralResponse{}, err // Return found error
	}

	var genesisTransactionStrings []string // Init string value buffer

	for _, transaction := range genesisTransactions { // Iterate through genesis transactions
		genesisTransactionStrings = append(genesisTransactionStrings, hex.EncodeToString(transaction.Hash.Bytes())) // Append hash
	}

	return &dagProto.GeneralResponse{Message: strings.Join(genesisTransactionStrings, ", ")}, nil // Return hashes
}

// GetTransactionByHash handles the GetTransactionByHash request method.
func (server *Server) GetTransactionByHash(ctx context.Context, request *dagProto.GeneralRequest) (*dagProto.GeneralResponse, error) {
	dag, err := types.OpenDag(request.Network) // Open dag

	if err != nil { // Check for errors
		return &dagProto.GeneralResponse{}, err // Return found error
	}

	transactionHashBytes, err := hex.DecodeString(request.TransactionHash) // Decode hash hex value

	if err != nil { // Check for errors
		return &dagProto.GeneralResponse{}, err // Return found error
	}

	transaction, err := dag.GetTransactionByHash(common.NewHash(transactionHashBytes)) // Query tx

	if err != nil { // Check for errors
		return &dagProto.GeneralResponse{}, err // Return found error
	}

	return &dagProto.GeneralResponse{Message: transaction.String()}, nil // Return tx JSON string value
}

/* END EXPORTED METHODS */
