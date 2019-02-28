// Package config provides DAG configuration helper methods and structs.
// Most notably, config provides the DagConfig struct, that of which is used to specify
// supply allocations, the dag identifier, and other metadata.
package config

import (
	"encoding/json"
	"io/ioutil"
)

// DagConfig represents a DAG configuration.
type DagConfig struct {
	Alloc map[string]float64 `json:"alloc"` // Account balances at genesis

	Identifier string `json:"identifier"` // Dag/network name (e.g. "mainnet_beta", "mainnet_alpha")

	NetworkID uint64 `json:"network"` // Dag version (e.g. 0 => mainnet, 1 => testnet, etc...)
}

/* BEGIN EXPORTED METHODS */

// NewDagConfig generates a new DagConfig from the given genesis.json file.
func NewDagConfig(genesisFilePath string) (*DagConfig, error) {
	rawJSON, err := ioutil.ReadFile(genesisFilePath) // Read genesis file

	if err != nil { // Check for errors
		return &DagConfig{}, err // Return found error
	}

	var readJSON map[string]interface{} // Init buffer

	err = json.Unmarshal(rawJSON, &readJSON) // Unmarshal to buffer

	if err != nil { // Check for errors
		return &DagConfig{}, err // Return found error
	}

	return &DagConfig{
		NetworkID:  readJSON["networkID"].(uint64),         // Set network ID
		Identifier: readJSON["identifier"].(string),        // Set ID
		Alloc:      readJSON["alloc"].(map[string]float64), // Set supply allocation
	}, nil
}

/* END EXPORTED METHODS */
