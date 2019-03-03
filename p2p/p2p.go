// Package p2p provides common peer-to-peer communications helper methods and definitions.
package p2p

import (
	"bufio"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/libp2p/go-libp2p-protocol"
	routed "github.com/libp2p/go-libp2p/p2p/host/routed"
	multiaddr "github.com/multiformats/go-multiaddr"
	"github.com/polaris-project/go-polaris/common"
)

var (
	// BootstrapNodes represents all default bootstrap nodes on the given network.
	BootstrapNodes = []string{
		"/ip4/108.41.124.60/tcp/53956/ipfs/QmWy8fZPX4hnTmXtFzgUTa8ZGceHhdhUEj3wonj1r3bMEG",
	}

	// WorkingDHT is the current global DHT instance.
	WorkingDHT *dht.IpfsDHT
)

/* BEGIN EXPORTED METHODS */

// NewHost initializes a new libp2p host with the given context.
func NewHost(ctx context.Context, port int) (*routed.RoutedHost, error) {
	peerIdentity, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	if _, err := os.Stat(filepath.FromSlash(fmt.Sprintf("%s/identity.pem", common.PeerIdentityDir))); err == nil { // Check existing p2p identity
		data, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/identity.pem", common.PeerIdentityDir))) // Read identity

		if err != nil { // Check for errors
			return nil, err // Return found error
		}

		block, _ := pem.Decode(data) // Decode pem

		peerIdentity, err = x509.ParseECPrivateKey(block.Bytes) // Parse private key pem block

		if err != nil { // Check for errors
			return nil, err // Return found error
		}
	} else { // No existing p2p identity
		x509Encoded, err := x509.MarshalECPrivateKey(peerIdentity) // Marshal identity

		if err != nil { // Check for errors
			return nil, err // Return found error
		}

		pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded}) // Encode to pem

		err = common.CreateDirIfDoesNotExit(common.PeerIdentityDir) // Create identity dir if it doesn't already exist

		if err != nil { // Check for errors
			return nil, err // Return found error
		}

		err = ioutil.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/identity.pem", common.PeerIdentityDir)), pemEncoded, 0644) // Write identity

		if err != nil { // Check for errors
			return nil, err // Return found error
		}
	}

	privateKey, _, err := crypto.ECDSAKeyPairFromKey(peerIdentity) // Get privateKey key

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	host, err := libp2p.New(ctx, libp2p.NATPortMap(), libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/"+strconv.Itoa(port), "/ip6/::1/tcp/"+strconv.Itoa(port)), libp2p.Identity(privateKey)) // Initialize libp2p host

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	err = BootstrapDht(ctx, host) // Bootstrap dht

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	return routed.Wrap(host, WorkingDHT), nil // Initialize routed libp2p host
}

// BootstrapDht bootstraps the WorkingDHT to the list of bootstrap nodes.
func BootstrapDht(ctx context.Context, host host.Host) error {
	var err error // Init error buffer

	WorkingDHT, err = dht.New(ctx, host) // Initialize DHT with host and context

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = WorkingDHT.Bootstrap(ctx) // Bootstrap

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

// BroadcastDht attempts to send a given message to all nodes in a dht at a given endpoint.
func BroadcastDht(ctx context.Context, host *routed.RoutedHost, message []byte, streamProtocol string, dagIdentifier string) error {
	peers := host.Peerstore().Peers() // Get peers

	for _, peer := range peers { // Iterate through peers
		if peer == (*host).ID() { // Check not same node
			continue // Continue
		}

		stream, err := (*host).NewStream(ctx, peer, protocol.ID(streamProtocol)) // Connect

		if err != nil { // Check for errors
			continue // Continue
		}

		writer := bufio.NewWriter(stream) // Initialize writer

		_, err = writer.Write(message) // Write message

		if err != nil { // Check for errors
			continue // Continue
		}
	}

	return nil // No error occurred, return nil
}

/* END EXPORTED METHODS */
