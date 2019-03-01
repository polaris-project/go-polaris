// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"testing"

	"github.com/polaris-project/go-polaris/config"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestNewDag tests the functionality of the NewDag() method.
func TestNewDag(t *testing.T) {
	dagConfig := config.NewDagConfig(nil, "test_network", 1) // Initialize new dag config with test genesis file.

	err := NewDag(dagConfig) // Initialize dag with dag config

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

/* END EXPORTED METHODS TESTS */
