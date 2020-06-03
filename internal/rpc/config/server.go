// Package config represents the config RPC server.
package config

import (
	"context"
	"io/ioutil"
	"strings"

	"github.com/polaris-project/go-polaris/common"

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

// GetAllConfigs handles the GetAllConfigs request method.
func (server *Server) GetAllConfigs(ctx context.Context, request *configProto.GeneralRequest) (*configProto.GeneralResponse, error) {
	files, err := ioutil.ReadDir(common.ConfigDir) // Read all files in config dir

	if err != nil { // Check for errors
		return &configProto.GeneralResponse{}, err // Return found error
	}

	networkNames := []string{} // Init network name buffer

	for _, file := range files { // Iterate through files
		networkNames = append(networkNames, strings.Split(strings.Split(file.Name(), "config_")[1], ".json")[0]) // Append network name
	}

	return &configProto.GeneralResponse{Message: strings.Join(networkNames, ", ")}, nil // Return network names
}

// GetConfig handles the GetConfig request method.
func (server *Server) GetConfig(ctx context.Context, request *configProto.GeneralRequest) (*configProto.GeneralResponse, error) {
	config, err := config.ReadDagConfigFromMemory(request.Network) // Read config from persistent memory

	if err != nil { // Check for errors
		return &configProto.GeneralResponse{}, err // Return found error
	}

	return &configProto.GeneralResponse{Message: config.String()}, nil // Return config
}

/* END EXPORTED METHODS */
