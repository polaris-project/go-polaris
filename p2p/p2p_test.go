// Package p2p provides common peer-to-peer communications helper methods and definitions.
package p2p

import (
	"context"
	"testing"

	"github.com/libp2p/go-libp2p"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestBootstrapDht tests the functionality of the BootstrapDht() helper method.
func TestBootstrapDht(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background()) // Get context

	defer cancel() // Cancel

	host, err := libp2p.New(ctx, libp2p.NATPortMap()) // Initialize libp2p host with context and nat manager

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = BootstrapDht(ctx, host) // Bootstrap

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

/* END EXPORTED METHODS TESTS */
