// Package config provides DAG configuration helper methods and structs.
// Most notably, config provides the DagConfig struct, that of which is used to specify
// supply allocations, the dag identifier, and other metadata.
package config

import "testing"

/* BEGIN EXPORTED METHODS TESTS */

// TestNewDagConfig tests the functionality of the NewDagConfig() helper method.
func TestNewDagConfig(t *testing.T) {
	dagConfig, err := NewDagConfig("test_genesis.json") // Initialize new dag config with test genesis file.

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(dagConfig) // Log success
}

/* END EXPORTED METHODS TESTS */
