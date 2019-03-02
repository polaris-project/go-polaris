// Package p2p provides common peer-to-peer communications helper methods and definitions.
package p2p

import (
	libp2pdht "github.com/libp2p/go-libp2p-kad-dht"
)

var (
	BootstrapNodes = []string{
		""
	}
	// WorkingDHT is the current global DHT instance.
	WorkingDHT libp2pdht.IpfsDHT
)
