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

	Network uint64 `json:"network"` // Dag version (e.g. 0 => mainnet, 1 => testnet, etc...)
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

	alloc := make(map[string]float64) // Init alloc map

	for key, value := range readJSON["alloc"].(map[string]interface{}) { // Iterate through genesis addresses
		alloc[key] = value.(float64) // Set alloc for address
	}

	return &DagConfig{
		Network:    uint64(readJSON["network"].(float64)), // Set network
		Identifier: readJSON["identifier"].(string),       // Set ID
		Alloc:      alloc,                                 // Set supply allocation
	}, nil
}

/* END EXPORTED METHODS */
