// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"testing"

	"github.com/polaris-project/go-polaris/config"
)

/* BEGIN EXPORTED METHODS */

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

/* END EXPORTED METHODS */
