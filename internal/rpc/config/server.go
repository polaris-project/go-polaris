package config

import (
	"context"

	"github.com/polaris-project/go-polaris/config"
	configProto "github.com/polaris-project/go-polaris/internal/proto/config"
)

// Server represents a Polaris RPC server.
type Server struct{}

/* BEGIN EXPORTED METHODS */

// NewDagConfig handles the NewDagConfig request method.
func (server *Server) NewDagConfig(ctx context.Context, request *configProto.GeneralRequest) (*configProto.GeneralResponse, error) {
	config, err := config.NewDagConfigFromGenesis(request.FilePath) // Initialize config

	if err != nil { // Check for errors
		return &configProto.GeneralResponse{}, err // Return found error
	}

	err = config.WriteToMemory() // Write config to memory

	if err != nil { // Check for errors
		return &configProto.GeneralResponse{}, err // Return found error
	}

	return &configProto.GeneralResponse{Message: config.String()}, nil // Return config
}

/* END EXPORTED METHODS */
