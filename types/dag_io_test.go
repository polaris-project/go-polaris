// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"bytes"
	"testing"

	"github.com/polaris-project/go-polaris/config"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestBytesDag tests the functionality of the dag Bytes() helper method.
func TestBytesDag(t *testing.T) {
	dagConfig := config.NewDagConfig(nil, "test_network", 1) // Initialize new dag config with test genesis file.

	dag, err := NewDag(dagConfig) // Initialize dag with dag config

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	t.Log(dag.Bytes()) // Log dag bytes

	WorkingDagDB.Close() // Close dag db
}

/* END EXPORTED METHODS TESTS */

/* BEGIN INTERNAL METHODS TESTS */

// TestReadDagDbHeaderFromMemory tests the functionality of the readDagDbHeaderFromMemory() helper method.
func TestReadDagDbHeaderFromMemory(t *testing.T) {
	dagConfig := config.NewDagConfig(nil, "test_network", 1) // Initialize new dag config with test genesis file.

	dag, err := NewDag(dagConfig) // Initialize dag with dag config

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = dag.WriteToMemory() // Write dag db header to persistent memory

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	readDag, err := readDagDbHeaderFromMemory(dag.DagConfig.Identifier) // Read dag db header

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	if !bytes.Equal(readDag.Bytes(), dag.Bytes()) { // Check dags not equivalent
		t.Fatal("dags should be equivalent") // Panic
	}

	WorkingDagDB.Close() // Close dag db
}

// TestWriteToMemoryDagDbHeader tests the functionality of the writeToMemory() dag db header helper.
func TestWriteToMemoryDagDbHeader(t *testing.T) {
	dagConfig := config.NewDagConfig(nil, "test_network", 1) // Initialize new dag config with test genesis file.

	dag, err := NewDag(dagConfig) // Initialize dag with dag config

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = dag.WriteToMemory() // Write dag db header to persistent memory

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	WorkingDagDB.Close() // Close working dag db
}

/* END INTERNAL METHODS TESTS */
