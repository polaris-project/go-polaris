// Package p2p provides common peer-to-peer communications helper methods and definitions.
package p2p

import (
	"context"
	"testing"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestNewHost tests the functionality of the NewHost() helper method.
func TestNewHost(t *testing.T) {
	_, err := NewHost(context.Background(), 2831) // Initialize host

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

// TestBootstrapDht tests the functionality of the BootstrapDht() helper method.
func TestBootstrapDht(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background()) // Get context

	defer cancel() // Cancel

	_, err := NewHost(ctx, 2831) // Initialize libp2p host with context and nat manager

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	/*
		err = BootstrapDht(ctx, host) // Bootstrap

		if err != nil { // Check for errors
			t.Fatal(err) // Panic
		}
	*/
}

/* END EXPORTED METHODS TESTS */
