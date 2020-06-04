// Package config provides DAG configuration helper methods and structs.
// Most notably, config provides the DagConfig struct, that of which is used to specify
// supply allocations, the dag identifier, and other metadata.
package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/polaris-project/go-polaris/common"
)

// ErrCouldNotDeserializeConfig represents an invalidly serialized dag config--that of which cannot be correctly deserialized.
var ErrCouldNotDeserializeConfig = errors.New("could not deserialize dag config")

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

// DagConfigFromBytes deserializes a dag config from the given byte array b.
func DagConfigFromBytes(b []byte) *DagConfig {
	buffer := &DagConfig{} // Init buffer

	err := json.Unmarshal(b, buffer) // Unmarshal into buffer
	if err != nil {                  // Check for errors
		return &DagConfig{} // Return nil dag config
	}

	return buffer // Return deserialized config
}

// WriteToMemory writes the given dag config to persistent memory.
func (dagConfig *DagConfig) WriteToMemory() error {
	err := common.CreateDirIfDoesNotExist(common.ConfigDir) // Create config dir if necessary
	if err != nil {                                         // Check for errors
		return err // Return found error
	}

	err = ioutil.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/config_%s.json", common.ConfigDir, dagConfig.Identifier)), dagConfig.Bytes(), 0o644) // Write dag config to persistent memory

	if err != nil { // Check for errors
		return err // Return error
	}

	return nil // No error occurred, return nil
}

// ReadDagConfigFromMemory reads a dag config with the given identifier from persistent memory.
func ReadDagConfigFromMemory(identifier string) (*DagConfig, error) {
	data, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/config_%s.json", common.ConfigDir, identifier))) // Read file
	if err != nil {                                                                                                  // Check for errors
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
