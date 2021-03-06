// Package common defines a set of commonly used helper methods and data types.
package common

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	// DataDir is the global data directory definition.
	DataDir = getDataDir()

	// ConfigDir is the global configuration directory definition.
	ConfigDir = filepath.FromSlash(fmt.Sprintf("%s/config", DataDir))

	// DbDir is the global database directory definition.
	DbDir = filepath.FromSlash(fmt.Sprintf("%s/db", DataDir))

	// PeerIdentityDir is the global p2p identity directory definition.
	PeerIdentityDir = filepath.FromSlash(fmt.Sprintf("%s/p2p", DataDir))

	// LogsDir is the global logs directory definition.
	LogsDir = filepath.FromSlash(fmt.Sprintf("%s/logs", DataDir))

	// CertificatesDir is the global certificate directory definition.
	CertificatesDir = filepath.FromSlash(fmt.Sprintf("%s/certs", DataDir))

	// MempoolDir is the global mempool directory definition.
	// The mempool is used to store cached transactions that have not been added to the public dag.
	MempoolDir = filepath.FromSlash(fmt.Sprintf("%s/mempool", DataDir))

	// KeystoreDir is the global keystore directory definition.
	KeystoreDir = filepath.FromSlash(fmt.Sprintf("%s/keystore", DataDir))
)

/* BEGIN EXPORTED METHODS */

// CreateDirIfDoesNotExist creates a given directory if it does not already exist.
func CreateDirIfDoesNotExist(dir string) error {
	safeDir := filepath.FromSlash(dir) // Just to be safe

	if _, err := os.Stat(safeDir); os.IsNotExist(err) { // Check dir exists
		err = os.MkdirAll(safeDir, 0o755) // Create directory

		if err != nil { // Check for errors
			return err // Return error
		}
	}

	return nil // No error occurred
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// getDataDir fetches the data directory
func getDataDir() string {
	abs, _ := filepath.Abs("./data") // Get absolute dir

	return filepath.FromSlash(abs) // Match OS slashes
}

/* END INTERNAL METHODS */
