// Package main is the main polaris mode entry point.
package main

import (
	"context"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"time"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"

	"github.com/polaris-project/go-polaris/config"
	"github.com/polaris-project/go-polaris/types"

	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/p2p"
	"github.com/polaris-project/go-polaris/validator"
)

var (
	// errNoBootstrap defines an invalid bootstrap value error.
	errNoBootstrap = errors.New("bootstrap failed: was expecting a bootstrap peer address, got 'localhost' (must be able to bootstrap dag config if no config exists locally)")
)

var (
	dataDirFlag              = flag.String("data-dir", common.DataDir, "performs all node I/O operations in a given data directory")                            // Init data dir flag
	nodePortFlag             = flag.Int("node-port", p2p.NodePort, "run p2p host on given port")                                                                // Init node port flag
	networkFlag              = flag.String("network", "main_net", "run node on given network")                                                                  // Init network flag
	bootstrapNodeAddressFlag = flag.String("bootstrap-address", p2p.BootstrapNodes[0], "manually prefer a given bootstrap node for all dht-related operations") // Init bootstrap node flag
	silencedFlag             = flag.Bool("silence", false, "silence logs")                                                                                      // Init silence logs flag
	disableColoredOutputFlag = flag.Bool("no-colors", false, "disable colored output")                                                                          // Init disable colored output flag

	logger = loggo.GetLogger("") // Get logger

	intermittentSyncContext, cancel = context.WithCancel(context.Background()) // Get background sync context
)

// Main starts all necessary Polaris services, and parses command line arguments.
func main() {
	flag.Parse() // Parse flags

	setUserParams() // Set common params

	defer cancel() // Cancel

	err := startNode() // Start node

	if err != nil { // Check for errors
		logger.Criticalf("main panicked: %s", err.Error()) // Log pending panic

		os.Exit(1) // Panic
	}
}

// setUserParams sets the default values in the common, p2p package.
func setUserParams() {
	common.DataDir = filepath.FromSlash(*dataDirFlag) // Set data dir

	p2p.NodePort = *nodePortFlag // Set node port

	if !*disableColoredOutputFlag { // Check can log colored output
		loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr)) // Enabled colored output
	}

	if *silencedFlag { // Check should silence
		loggo.ResetLogging() // Silence
	}
}

// startNode creates a new libp2p host, and connects to the bootstrapped dht.
func startNode() error {
	ctx, cancel := context.WithCancel(context.Background()) // Init context

	defer cancel() // Cancel context

	host, err := p2p.NewHost(ctx, p2p.NodePort) // Initialize host

	if err != nil { // Check for errors
		return err // Return found error
	}

	if *bootstrapNodeAddressFlag == p2p.BootstrapNodes[0] { // Check bootstrap node addr has not been set
		*bootstrapNodeAddressFlag = p2p.GetBestBootstrapAddress(context.Background(), host) // Get best bootstrap node
	}

	dagConfig, err := config.ReadDagConfigFromMemory(*networkFlag) // Read config

	needsSync := false // Assume doesn't need sync

	if err != nil || dagConfig == nil { // Check no existing dag config
		if *bootstrapNodeAddressFlag == "localhost" { // Check no bootstrap node
			return errNoBootstrap // Return error
		}

		dagConfig, err = p2p.BootstrapConfig(ctx, host, *bootstrapNodeAddressFlag, *networkFlag) // Bootstrap dag config

		if err != nil { // Check for errors
			return err // Return found error
		}

		needsSync = true // Set does need sync
	}

	dag, err := types.NewDag(dagConfig) // Init dag

	if err != nil { // Check for errors
		return err // Return found error
	}

	validator := validator.Validator(validator.NewBeaconDagValidator(dagConfig, dag)) // Initialize validator

	client := p2p.NewClient(*networkFlag, &validator) // Initialize client

	err = client.StartServingStreams(*networkFlag) // Start handlers

	if err != nil { // Check for errors
		return err // Return found error
	}

	if needsSync { // Check must sync
		err = client.SyncDag(ctx) // Sync network

		if err != nil { // Check for errors
			return err // Return found error
		}
	}

	client.StartIntermittentSync(intermittentSyncContext, 120*time.Second) // Sync every 120 seconds

	return nil // No error occurred, return nil
}
