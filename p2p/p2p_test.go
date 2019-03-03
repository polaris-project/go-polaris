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

	BootstrapNodes = []string{
		"/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
		"/ip4/104.236.179.241/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",
	} // Set bootstrap nodes

	_, err := NewHost(ctx, 2831) // Initialize libp2p host with context and nat manager

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

// TestBroadcastDht tests the functionality of the BroadcastDht() helper method.
func TestBroadcastDht(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background()) // Get context

	defer cancel() // Cancel

	BootstrapNodes = []string{
		"/ip4/104.131.131.82/tcp/4001/ipfs/QmaCpDMGvV2BGHeYERUEnRQAwe3N8SzbUtfsmvsqQLuvuJ",
		"/ip4/104.236.179.241/tcp/4001/ipfs/QmSoLPppuBtQSGwKDZT2M73ULpjvfd3aZ6ha4oFGL1KrGM",
	} // Set bootstrap nodes

	host, err := NewHost(ctx, 2831) // Initialize libp2p host with context and nat manager

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = BroadcastDht(ctx, host, []byte("test"), "/test/1.0.0", "test_network") // Broadcast

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

/* END EXPORTED METHODS TESTS */
