package p2p

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	inet "github.com/libp2p/go-libp2p-net"
	"github.com/polaris-project/go-polaris/config"
	"github.com/polaris-project/go-polaris/types"
	"github.com/polaris-project/go-polaris/validator"
)

/* BEGIN EXPORTED METHODS TESTS */

// TestStartServingStream tests the functionality of the StartServingStream() helper method.
func TestStartServingStream(t *testing.T) {
	os.RemoveAll(filepath.FromSlash("data/db/test_network.db")) // Remove existing db

	dagConfig := config.NewDagConfig(nil, "test_network", 1) // Initialize new dag config with test genesis file.

	dag, err := types.NewDag(dagConfig) // Initialize dag with dag config

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	validator := validator.Validator(validator.NewBeaconDagValidator(dagConfig, dag)) // Initialize validator

	client := NewClient("test_network", &validator) // Initialize client

	if client == nil { // Check client is nil
		t.Fatal("client should not be nil") // Panic
	}

	ctx, cancel := context.WithCancel(context.Background()) // Initialize context

	defer cancel() // Cancel

	_, err = NewHost(ctx, 3000) // Initialize host

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = client.StartServingStream("test_stream", func(inet.Stream) {}) // Start serving stream

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	os.RemoveAll(filepath.FromSlash("data/db/test_network.db")) // Remove existing db
}

/* END EXPORTED METHODS TESTS */
