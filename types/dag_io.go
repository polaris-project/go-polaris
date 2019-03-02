// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/polaris-project/go-polaris/common"
)

/* BEGIN EXPORTED METHODS */

// Bytes serializes a given dag header to a byte array via JSON.
func (dag *Dag) Bytes() []byte {
	marshaledVal, _ := json.MarshalIndent(*dag, "", "  ") // Marshal

	return marshaledVal // Return bytes
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// readDagDbHeaderFromMemory attempts to read the dag db header with the given identifier from persistent memory.
func readDagDbHeaderFromMemory(identifier string) (*Dag, error) {
	data, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/db_header_%s.json", common.DbDir, identifier))) // Read header

	if err != nil { // Check for errors
		return &Dag{}, err // Return found error
	}

	buffer := &Dag{} // Initialize db header buffer

	err = json.Unmarshal(data, buffer) // Unmarshal JSON into buffer

	if err != nil { // Check for errors
		return &Dag{}, err // Return found error
	}

	return buffer, nil // No error occurred, retrun nil
}

// writeToMemory writes the dag header to persistent memory.
func (dag *Dag) writeToMemory() error {
	err := common.CreateDirIfDoesNotExit(common.DbDir) // Create db dir if necessary

	if err != nil { // Check for errors
		return err // Return found error
	}

	return ioutil.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/db_header_%s.json", common.DbDir, dag.DagConfig.Identifier)), dag.Bytes(), 0644) // Write dag header to persistent memory
}

/* END INTERNAL METHODS */
