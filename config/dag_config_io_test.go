// Package config provides DAG configuration helper methods and structs.
// Most notably, config provides the DagConfig struct, that of which is used to specify
// supply allocations, the dag identifier, and other metadata.
package config

import (
	"bytes"
	"testing"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestString tests the functionality of the dag config string() helper method.
func TestString(t *testing.T) {
	dagConfig, err := NewDagConfigFromGenesis("test_genesis.json") // Initialize new dag config with test genesis file.

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	if dagConfig.String() == "" { // Check nil string
		t.Fatal("dag config string value should not be nil") // Panic
	}

	t.Log(dagConfig.String()) // Log success
}

// TestBytes tests the functionality fo the dag config bytes() helper method.
func TestBytes(t *testing.T) {
	dagConfig, err := NewDagConfigFromGenesis("test_genesis.json") // Initialize new dag config with test genesis file.

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	if dagConfig.Bytes() == nil { // Check nil byte val
		t.Fatal("dag config bytes value should not be nil") // Panic
	}
}

// TestWriteToMemory tests the functionality of outbound dag config I/O.
func TestWriteToMemory(t *testing.T) {
	dagConfig, err := NewDagConfigFromGenesis("test_genesis.json") // Initialize new dag config with test genesis file.

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = dagConfig.WriteToMemory() // Write dag config to persistent memory

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

// TestReadDagConfigFromMemory tests the functionality of inbound dag config I/O.
func TestReadDagConfigFromMemory(t *testing.T) {
	dagConfig, err := NewDagConfigFromGenesis("test_genesis.json") // Initialize new dag config with test genesis file.

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = dagConfig.WriteToMemory() // Write dag config to persistent memory

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	readDagConfig, err := ReadDagConfigFromMemory(dagConfig.Identifier) // Read dag config

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	if !bytes.Equal(dagConfig.Bytes(), readDagConfig.Bytes()) { // Check dag configs not equivalent
		t.Fatal("dag configs should be equivalent") // Panic
	}
}

/* END EXPORTED METHODS TESTS */
