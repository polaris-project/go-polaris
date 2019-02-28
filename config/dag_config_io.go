// Package config provides DAG configuration helper methods and structs.
// Most notably, config provides the DagConfig struct, that of which is used to specify
// supply allocations, the dag identifier, and other metadata.
package config

import "encoding/json"

/* BEGIN EXPORTED METHODS */

// String serializes a given dag config to a string via json.
func (dagConfig *DagConfig) String() string {
	marshaledVal, _ := json.MarshalIndent(*dagConfig, "", "  ") // Marshal JSON

	return string(marshaledVal) // Return marshalled JSON as a string
}

/* END EXPORTED METHODS */
