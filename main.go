// Package main is the main polaris mode entry point.
package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	routed "github.com/libp2p/go-libp2p/p2p/host/routed"

	"github.com/juju/loggo"
	"github.com/juju/loggo/loggocolor"

	"github.com/polaris-project/go-polaris/config"
	"github.com/polaris-project/go-polaris/types"

	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/p2p"
	"github.com/polaris-project/go-polaris/validator"

	"github.com/polaris-project/go-polaris/api"
	"github.com/polaris-project/go-polaris/cli"
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
	disableLogFileFlag       = flag.Bool("no-logs", false, "disable writing logs to a logs.txt file")                                                           // Init disable logs file flag
	debugFlag                = flag.Bool("debug", false, "force node to log in debug mode")                                                                     // Init debug flag
	disableAutoGenesisFlag   = flag.Bool("no-genesis", false, "disables the automatic creation of a genesis transaction set if no dag can be bootstrapped")     // Init disable auto genesis flag
	apiPortFlag              = flag.Int("api-port", 8000, "incrementally start APIs on given port (i.e. RPC = 8080, HTTP = 8081, etc...)")                      // Init API port flag
	disableAPIFlag           = flag.Bool("disable-api", false, "disable API")                                                                                   // Init disable API flag
	terminalFlag             = flag.Bool("terminal", false, "launch with terminal")                                                                             // Init terminal flag
	rpcPortFlag              = flag.Int("rpc-port", 8000, "port to connect to via RPC")                                                                         // Init RPC port flag
	rpcAddrFlag              = flag.String("rpc-address", "localhost", "RPC addr to connect to")                                                                // Init RPC addr flag

	logger = loggo.GetLogger("") // Get logger

	logFile *os.File   // Log file
	dag     *types.Dag // Dag reference

	intermittentSyncContext, cancelIntermittent = context.WithCancel(context.Background()) // Get background sync context
)

// Main starts all necessary Polaris services, and parses command line arguments.
func main() {
	flag.Parse() // Parse flags

	err := setUserParams() // Set common params

	if err != nil { // Check for errors
		logger.Criticalf("main panicked: %s", err.Error()) // Log pending panic

		os.Exit(1) // Panic
	}

	defer dag.Close()     // Close dag
	defer logFile.Close() // Close log file

	if !checkNodeAlreadyUp() { // Check no node already up
		defer cancelIntermittent() // Cancel

		err = startNode() // Start node

		if err != nil { // Check for errors
			logger.Criticalf("main panicked: %s", err.Error()) // Log pending panic

			os.Exit(1) // Panic
		}
	}

	if *terminalFlag { // Check should launch terminal
		cli.NewTerminal(uint(*rpcPortFlag), *rpcAddrFlag) // Initialize terminal
	}
}

// setUserParams sets the default values in the common, p2p package.
func setUserParams() error {
	common.DataDir = filepath.FromSlash(*dataDirFlag) // Set data dir

	p2p.NodePort = *nodePortFlag // Set node port

	if !*disableColoredOutputFlag { // Check can log colored output
		if !*disableLogFileFlag { // Check can have log files
			err := common.CreateDirIfDoesNotExit(filepath.FromSlash(common.LogsDir)) // Create log dir

			if err != nil { // Check for errors
				return err // Return found error
			}

			logFile, err = os.OpenFile(filepath.FromSlash(fmt.Sprintf("%s/logs_%s.txt", common.LogsDir, time.Now().Format("2006-01-02_15-04-05"))), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) // Create log file

			if err != nil { // Check for errors
				return err // Return found error
			}

			loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr)) // Enabled colored output

			loggo.RegisterWriter("logs", loggo.NewSimpleWriter(logFile, loggo.DefaultFormatter)) // Register file writer
		} else {
			loggo.ReplaceDefaultWriter(loggocolor.NewWriter(os.Stderr)) // Enabled colored output
		}
	}

	if *silencedFlag { // Check should silence
		loggo.ResetLogging() // Silence
	}

	return nil // No error occurred, return nil.
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

	dagConfig, needsSync, err := getDagConfig(ctx, host) // Get dag config

	if err != nil { // Check for errors
		return err // Return found error
	}

	dag, err = types.NewDag(dagConfig) // Init dag

	if err != nil { // Check for errors
		return err // Return found error
	}

	c := make(chan os.Signal) // Get control c

	signal.Notify(c, os.Interrupt, syscall.SIGTERM) // Notify

	go func() {
		<-c // Wait for ^c

		cancelIntermittent() // Cancel intermittent sync
		cancel()             // Cancel

		err = dag.Close() // Close dag

		if err != nil { // Check for errors
			logger.Criticalf("dag close errored: %s", err.Error()) // Return found error
		}

		logFile.Close() // Close log file

		os.Exit(0) // Exit
	}()

	validator := validator.Validator(validator.NewBeaconDagValidator(dagConfig, dag)) // Initialize validator

	client := p2p.NewClient(*networkFlag, &validator) // Initialize client

	needsSync, err = startInitialSync(ctx, needsSync, client) // Start initial sync

	if err != nil { // Check for errors
		return err // Return found error
	}

	p2p.WorkingClient = client // Set working client

	err = p2p.WorkingClient.StartServingStreams(*networkFlag) // Start handlers

	if err != nil { // Check for errors
		return err // Return found error
	}

	if !*disableAPIFlag { // Check API enabled
		err = startServingRPC(ctx, dagConfig.Identifier) // Start serving RPC

		if err != nil { // Check for errors
			return err // Return found error
		}
	}

	if needsSync { // Check must sync
		err = p2p.WorkingClient.SyncDag(ctx) // Sync network

		if err != nil { // Check for errors
			return err // Return found error
		}
	}

	if !*terminalFlag { // Check no terminal
		p2p.WorkingClient.StartIntermittentSync(intermittentSyncContext, 120*time.Second) // Sync every 120 seconds
	} else {
		go p2p.WorkingClient.StartIntermittentSync(intermittentSyncContext, 120*time.Second) // Sync every 120 seconds
	}

	return nil // No error occurred, return nil
}

// startInitialSync starts an initial sync with a given client.
func startInitialSync(ctx context.Context, needsSync bool, client *p2p.Client) (bool, error) {
	localBestTransaction, _ := (*client.Validator).GetWorkingDag().GetBestTransaction() // Get local best transaction

	remoteBestTransactionHash, _ := client.RequestBestTransactionHash(ctx, 64) // Request best tx hash

	if !bytes.Equal(localBestTransaction.Hash.Bytes(), remoteBestTransactionHash.Bytes()) && !localBestTransaction.Hash.IsNil() && !remoteBestTransactionHash.IsNil() { // Check up to date
		needsSync = true // Set does need sync
	}

	if localBestTransaction.Hash.IsNil() && remoteBestTransactionHash.IsNil() { // Check nil genesis
		_, err := (*client.Validator).GetWorkingDag().MakeGenesis() // Make genesis

		if err != nil { // Check for errors
			return false, err // Return found error
		}
	}

	return needsSync, nil // No error occurred, return nil
}

// startServingRPC starts the RPC node interface.
func startServingRPC(ctx context.Context, network string) error {
	rpcAPI, err := api.NewRPCAPI(network, fmt.Sprintf(":%d", *apiPortFlag)) // Initialize API

	if err != nil { // Check for errors
		return err // Return found error
	}

	go rpcAPI.StartServing(ctx) // Start serving

	return nil // No error occurred, return nil
}

// getDagConfig attempts to read an existing dag config, or bootstrap one.
func getDagConfig(ctx context.Context, host *routed.RoutedHost) (*config.DagConfig, bool, error) {
	needsSync := false // Init buffer

	dagConfig, err := config.ReadDagConfigFromMemory(*networkFlag) // Read config

	if err != nil || dagConfig == nil { // Check no existing dag config
		if *bootstrapNodeAddressFlag == "localhost" { // Check no bootstrap node
			return &config.DagConfig{}, needsSync, errNoBootstrap // Return error
		}

		dagConfig, err = p2p.BootstrapConfig(ctx, host, *bootstrapNodeAddressFlag, *networkFlag) // Bootstrap dag config

		if err != nil { // Check for errors
			return &config.DagConfig{}, needsSync, err // Return found error
		}

		needsSync = true // Set does need sync
	}

	return dagConfig, needsSync, nil // Return found config
}

// checkNodeAlreadyUp checks if a node RPC API is already running.
func checkNodeAlreadyUp() bool {
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(*apiPortFlag)) // Attempt to listen

	if err != nil {
		return true // Already running
	}

	ln.Close() // Close

	return false // Not running
}
