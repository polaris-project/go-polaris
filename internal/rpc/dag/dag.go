// Package dag represents the dag RPC server.
package dag

import (
	"context"

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

/* END EXPORTED METHODS */
