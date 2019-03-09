// Package main is the main polaris mode entry point.
package main

import (
	"context"
	"flag"
	"path/filepath"

	"github.com/polaris-project/go-polaris/types"

	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/p2p"
	"github.com/polaris-project/go-polaris/validator"
)

var (
	dataDirFlag          = flag.String("data-dir", common.DataDir, "performs all node I/O operations in a given data directory")                            // Init data dir flag
	nodePortFlag         = flag.Int("node-port", p2p.NodePort, "run p2p host on given port")                                                                // Init node port flag
	networkFlag          = flag.String("network", "main_net", "run node on given network")                                                                  // Init network flag
	bootstrapNodeAddress = flag.String("bootstrap-address", p2p.BootstrapNodes[0], "manually prefer a given bootstrap node for all dht-related operations") // Init bootstrap node flag
)

// Main starts all necessary Polaris services, and parses command line arguments.
func main() {
	flag.Parse() // Parse flags

	setUserParams() // Set common params
}

// setUserParams sets the default values in the common, p2p package.
func setUserParams() {
	common.DataDir = filepath.FromSlash(*dataDirFlag) // Set data dir

	p2p.NodePort = *nodePortFlag // Set node port
}

// startNode creates a new libp2p host, and connects to the bootstrapped dht.
func startNode() error {
	ctx, cancel := context.WithCancel(context.Background()) // Init context

	defer cancel() // Cancel context

	host, err := p2p.NewHost(ctx, p2p.NodePort) // Initialize host

	if err != nil { // Check for errors
		return err // Return found error
	}

	if *bootstrapNodeAddress == p2p.BootstrapNodes[0] { // Check bootstrap node addr has not been set
		*bootstrapNodeAddress = p2p.GetBestBootstrapAddress(context.Background(), host) // Get best bootstrap node
	}

	config, err := p2p.BootstrapConfig(ctx, host, *bootstrapNodeAddress, *networkFlag) // Bootstrap dag config

	if err != nil { // Check for errors
		return err // Return found error
	}

	dag, err := types.NewDag(config) // Init dag

	if err != nil { // Check for errors
		return err // Return found error
	}

	validator := validator.Validator(validator.NewBeaconDagValidator(config, dag)) // Initialize validator

	client := p2p.NewClient(*networkFlag, &validator) // Initialize client

	return nil // No error occurred, return nil
}
