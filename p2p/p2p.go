// Package p2p provides common peer-to-peer communications helper methods and definitions.
package p2p

import (
	"context"

	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	multiaddr "github.com/multiformats/go-multiaddr"
)

var (
	// BootstrapNodes represents all default bootstrap nodes on the given network.
	BootstrapNodes = []string{
		"/ip4/108.41.124.60/tcp/3333",
	}

	// WorkingDHT is the current global DHT instance.
	WorkingDHT *dht.IpfsDHT
)

/* BEGIN EXPORTED METHODS */

// BootstrapDht bootstraps the WorkingDHT to the list of bootstrap nodes.
func BootstrapDht(ctx context.Context, host host.Host) error {
	var err error // Init error buffer

	WorkingDHT, err = dht.New(ctx, host) // Initialize DHT with host and context

	if err != nil { // Check for errors
		return err // Return found error
	}

	for _, addr := range BootstrapNodes { // Iterate through bootstrap nodes
		address, err := multiaddr.NewMultiaddr(addr) // Parse multi address

		if err != nil { // Check for errors
			continue // Continue to next peer
		}

		peerInfo, err := peerstore.InfoFromP2pAddr(address) // Get peer info

		if err != nil { // Check for errors
			continue // Continue to next peer
		}

		err = host.Connect(ctx, *peerInfo) // Connect to discovered peer

		if err != nil { // Check for errors
			continue // Continue to next peer
		}
	}

	return nil // No error occurred, return nil
}

/* END EXPORTED METHODS */
