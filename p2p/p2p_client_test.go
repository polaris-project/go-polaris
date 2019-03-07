package p2p

import (
	"context"
	"testing"

	inet "github.com/libp2p/go-libp2p-net"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestNewClient tests the functionality of the client initializer.
func TestNewClient(t *testing.T) {
	_, err := NewHost(context.Background(), 2831) // Initialize host

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	client := NewClient("test_network") // Initialize client

	if client == nil { // Check client is nil
		t.Fatal("client should not be nil") // Panic
	}
}

// TestStartServingStream tests the functionality of the StartServingStream() helper method.
func TestStartServingStream(t *testing.T) {
	_, err := NewHost(context.Background(), 2831) // Initialize host

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	client := NewClient("test_network") // Initialize client

	if client == nil { // Check client is nil
		t.Fatal("client should not be nil") // Panic
	}

	err = client.StartServingStream("test_stream", func(inet.Stream) {}) // Start serving stream

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

/* END EXPORTED METHODS TESTS */
