// Package p2p provides common peer-to-peer communications helper methods and definitions.
package p2p

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"
	"testing"

	protocol "github.com/libp2p/go-libp2p-protocol"
	"github.com/polaris-project/go-polaris/crypto"
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

// TestPublish tests the functionality of the Publish() helper method.
func TestPublish(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	address := crypto.AddressFromPrivateKey(privateKey) // Generate address

	transaction := NewTransaction(
		0,                      // Nonce
		big.NewFloat(0),        // Amount
		address,                // Sender
		nil,                    // Recipient
		nil,                    // Parents
		0,                      // Gas limit
		big.NewInt(0),          // Gas price
		[]byte("test payload"), // Payload
	) // Initialize a new transaction

	err = SignTransaction(transaction, privateKey) // Sign transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	ctx, cancel := context.WithCancel(context.Background()) // Get context

	defer cancel() // Cancel

	_, err = NewHost(ctx, 3861) // Initialize host

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}

	err = transaction.Publish(context.Background(), "test_network") // Publish transaction

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

// TestGetStreamHeaderProtocolPath tests the functionality of the GetStreamHeaderProtocol() helper method.
func TestGetStreamHeaderProtocolPath(t *testing.T) {
	streamHeaderProtocolPath := GetStreamHeaderProtocolPath("test_network", PublishTransaction) // Get stream header protocol URI

	t.Log(protocol.ID(streamHeaderProtocolPath)) // Get libp2p representation
}

/* END EXPORTED METHODS TESTS */
