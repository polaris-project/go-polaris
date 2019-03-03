// Package p2p provides common peer-to-peer communications helper methods and definitions.
package p2p

import (
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
	multiaddr "github.com/multiformats/go-multiaddr"
	"github.com/polaris-project/go-polaris/common"
)

var (
	// BootstrapNodes represents all default bootstrap nodes on the given network.
	BootstrapNodes = []string{
		"/ip4/108.41.124.60/tcp/53956/ipfs/",
	}

	// WorkingDHT is the current global DHT instance.
	WorkingDHT *dht.IpfsDHT
)

/* BEGIN EXPORTED METHODS */

// NewHost initializes a new libp2p host with the given context.
func NewHost(ctx context.Context, port int) (host.Host, error) {
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

	return libp2p.New(ctx, libp2p.NATPortMap(), libp2p.ListenAddrStrings("/ip4/127.0.0.1/tcp/"+strconv.Itoa(port), "/ip6/::1/tcp/"+strconv.Itoa(port)), libp2p.Identity(privateKey)) // Initialize libp2p host
}

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
