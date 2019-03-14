// Package p2p provides common peer-to-peer communications helper methods and definitions.
package p2p

import (
	"bufio"
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/libp2p/go-libp2p/p2p/protocol/ping"

	"github.com/libp2p/go-libp2p"
	crypto "github.com/libp2p/go-libp2p-crypto"
	host "github.com/libp2p/go-libp2p-host"
	dht "github.com/libp2p/go-libp2p-kad-dht"
	peer "github.com/libp2p/go-libp2p-peer"
	peerstore "github.com/libp2p/go-libp2p-peerstore"
	protocol "github.com/libp2p/go-libp2p-protocol"
	routed "github.com/libp2p/go-libp2p/p2p/host/routed"
	multiaddr "github.com/multiformats/go-multiaddr"
	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/config"
)

// Stream header protocol definitions
const (
	PublishTransaction StreamHeaderProtocol = iota

	RequestConfig

	RequestBestTransaction

	RequestTransaction

	RequestGenesisHash

	RequestChildHashes
)

var (
	// StreamHeaderProtocolNames represents all stream header protocol names.
	StreamHeaderProtocolNames = []string{
		"pub_transaction",
		"req_config",
		"req_best_transaction",
		"req_transaction",
		"req_genesis_hash",
		"req_transaction_children_hashes",
	}

	// BootstrapNodes represents all default bootstrap nodes on the given network.
	BootstrapNodes = []string{
		"/ip4/108.41.124.60/tcp/3030/ipfs/QmWy8fZPX4hnTmXtFzgUTa8ZGceHhdhUEj3wonj1r3bMEG",
	}

	// WorkingHost is the current global routed host.
	WorkingHost *routed.RoutedHost

	// NodePort is the current node port
	NodePort = 3030

	// ErrTimedOut is an error definition representing a timeout.
	ErrTimedOut = errors.New("timed out")
)

// StreamHeaderProtocol represents the stream protocol type enum.
type StreamHeaderProtocol int

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

	host, err := libp2p.New(ctx, libp2p.NATPortMap(), libp2p.ListenAddrStrings("/ip4/0.0.0.0/tcp/"+strconv.Itoa(port), "/ip6/::1/tcp/"+strconv.Itoa(port)), libp2p.Identity(privateKey)) // Initialize libp2p host

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	dht, err := BootstrapDht(ctx, host) // Bootstrap dht

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	routedHost := routed.Wrap(host, dht) // Initialize routed host

	WorkingHost = routedHost // Set routed host

	return WorkingHost, nil // Return working routed host
}

// GetBestBootstrapAddress attempts to fetch the best bootstrap node.
func GetBestBootstrapAddress(ctx context.Context, host *routed.RoutedHost) string {
	for _, bootstrapAddress := range BootstrapNodes { // Iterate through bootstrap nodes
		multiaddr, err := multiaddr.NewMultiaddr(bootstrapAddress) // Parse address

		if err != nil { // Check for errors
			continue // Continue
		}

		peerID, err := peer.IDB58Decode(strings.Split(bootstrapAddress, "ipfs/")[1]) // Get peer ID

		if err != nil { // Check for errors
			continue // Continue
		}

		host.Peerstore().AddAddr(peerID, multiaddr, 10*time.Second) // Add bootstrap peer

		peerInfo, err := peerstore.InfoFromP2pAddr(multiaddr) // Get peer info

		if err != nil { // Check for errors
			continue // Continue
		}

		bootstrapCheckCtx, cancel := context.WithCancel(ctx) // Get context

		err = host.Connect(bootstrapCheckCtx, *peerInfo) // Connect to peer

		if err != nil { // Check for errors
			cancel() // Cancel
			continue // Continue
		}

		_, err = ping.Ping(bootstrapCheckCtx, host, peerID) // Attempt to ping

		if err == nil { // Check no errors
			cancel()                // Cancel
			return bootstrapAddress // Return bootstrap address
		}

		cancel() // Cancel
	}

	return "localhost" // Return localhost
}

// BootstrapConfig bootstraps a dag config to the list of bootstrap nodes.
func BootstrapConfig(ctx context.Context, host *routed.RoutedHost, bootstrapAddress string, network string) (*config.DagConfig, error) {
	peerID, err := peer.IDB58Decode(strings.Split(bootstrapAddress, "ipfs/")[1]) // Get peer ID

	if err != nil { // Check for errors
		return &config.DagConfig{}, err // Return found error
	}

	readCtx, cancel := context.WithCancel(ctx) // Get context

	stream, err := (*host).NewStream(readCtx, peerID, protocol.ID(GetStreamHeaderProtocolPath(network, RequestConfig))) // Initialize new stream

	if err != nil { // Check for errors
		cancel() // Cancel

		return &config.DagConfig{}, err // Return found error
	}

	reader := bufio.NewReader(stream) // Initialize reader from stream

	var dagConfigBytes bytes.Buffer // Initialize dag config bytes buffer

	readStartTime := time.Now() // Get start time

	finished := false // Init finished bool

	finishedChan := &finished // Get finished ref

	for dagConfigBytes.Bytes() == nil || len(dagConfigBytes.Bytes()) == 0 { // Read while nil
		go func() {
			if *finishedChan != true { // Check not finished
				io.Copy(&dagConfigBytes, reader) // Non-blocking read
			}
		}()

		if time.Now().Sub(readStartTime) > 10*time.Second { // Check for timeout
			cancel() // Cancel

			return &config.DagConfig{}, ErrTimedOut // Return found error
		}
	}

	*finishedChan = true // Set finished

	deserializedConfig := config.DagConfigFromBytes(dagConfigBytes.Bytes()) // Deserialize

	if deserializedConfig == nil { // Check nil
		cancel() // Cancel

		return &config.DagConfig{}, config.ErrCouldNotDeserializeConfig // Return error
	}

	cancel() // Cancel

	return deserializedConfig, nil // Return deserialized dag config
}

// BootstrapDht bootstraps the WorkingDHT to the list of bootstrap nodes.
func BootstrapDht(ctx context.Context, host host.Host) (*dht.IpfsDHT, error) {
	dht, err := dht.New(ctx, host) // Initialize DHT with host and context

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	err = dht.Bootstrap(ctx) // Bootstrap

	if err != nil { // Check for errors
		return nil, err // Return found error
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

	return dht, nil // No error occurred, return nil
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

// BroadcastDhtResult send a given message to all nodes in a dht, and returns the result from each node.
func BroadcastDhtResult(ctx context.Context, host *routed.RoutedHost, message []byte, streamProtocol string, dagIdentifier string, nPeers int) ([][]byte, error) {
	peers := host.Peerstore().Peers() // Get peers

	results := [][]byte{} // Init results buffer

	for x, peer := range peers { // Iterate through peers
		if x >= nPeers { // Check has sent to enough peers
			break // Break
		}

		if peer == (*host).ID() { // Check not same node
			continue // Continue
		}

		stream, err := (*host).NewStream(ctx, peer, protocol.ID(streamProtocol)) // Connect

		if err != nil { // Check for errors
			continue // Continue
		}

		readWriter := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream)) // Initialize reader/writer

		_, err = readWriter.Write(message) // Write message

		if err != nil { // Check for errors
			continue // Continue
		}

		var responseBytes []byte // Initialize response bytes buffer

		for readByte, err := readWriter.ReadByte(); err != nil; { // Read until EOF
			responseBytes = append(responseBytes, readByte) // Append read byte
		}

		results = append(results, responseBytes) // Append response
	}

	return results, nil // No error occurred, return response
}

// GetStreamHeaderProtocolPath attempts to determine the libp2p stream header protocol URI from a given stream protocol and network.
func GetStreamHeaderProtocolPath(network string, streamProtocol StreamHeaderProtocol) string {
	return fmt.Sprintf("/%s/%s", network, StreamHeaderProtocolNames[streamProtocol]) // Return URI
}

/* END EXPORTED METHODS */
