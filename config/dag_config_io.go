// Package config provides DAG configuration helper methods and structs.
// Most notably, config provides the DagConfig struct, that of which is used to specify
// supply allocations, the dag identifier, and other metadata.
package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/polaris-project/go-polaris/common"
)

/* BEGIN EXPORTED METHODS */

// String serializes a given dag config to a string via json.
func (dagConfig *DagConfig) String() string {
	marshaledVal, _ := json.MarshalIndent(*dagConfig, "", "  ") // Marshal JSON

	return string(marshaledVal) // Return marshalled JSON as a string
}

// Bytes serializes a given dag config to a byte array via json.
func (dagConfig *DagConfig) Bytes() []byte {
	marshaledVal, _ := json.MarshalIndent(*dagConfig, "", "  ") // Marshal JSON

	return marshaledVal // Return marshalled JSON
}

// WriteToMemory writes the given dag config to persistent memory.
func (dagConfig *DagConfig) WriteToMemory() error {
	err := common.CreateDirIfDoesNotExit(common.ConfigDir) // Create config dir if necessary

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = ioutil.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/config_%s.json", common.DataDir, dagConfig.Identifier)), dagConfig.Bytes(), 0644) // Write dag config to persistent memory

	if err != nil { // Check for errors
		return err // Return error
	}

	return nil // No error occurred, return nil
}

// ReadDagConfigFromMemory reads a dag config with the given identifier from persistent memory.
func ReadDagConfigFromMemory(identifier string) (*DagConfig, error) {
	data, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/%s", common.ConfigDir, identifier))) // Read file

	if err != nil { // Check for errors
		return &DagConfig{}, err // Return found error
	}

	buffer := &DagConfig{} // Initialize buffer

	err = json.Unmarshal(data, buffer) // Deserialize JSON into buffer.

	if err != nil { // Check for errors
		return &DagConfig{}, err // Return found error
	}

	return buffer, nil // No error occurred, return read dag config
}

/* END EXPORTED METHODS */
