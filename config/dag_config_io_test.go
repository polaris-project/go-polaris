// Package config provides DAG configuration helper methods and structs.
// Most notably, config provides the DagConfig struct, that of which is used to specify
// supply allocations, the dag identifier, and other metadata.
package config

import "testing"

/* BEGIN EXPORTED METHODS TESTS */

// TestString tests the functionality of the dag config string() helper method.
func TestString(t *testing.T) {
	dagConfig, err := NewDagConfig("test_genesis.json") // Initialize new dag config with test genesis file.

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	if dagConfig.String() == "" { // Check nil string
		t.Fatal("dag config string value should not be nil") // Panic
	}

	t.Log(dagConfig.String()) // Log success
}

/* END EXPORTED METHODS TESTS */
