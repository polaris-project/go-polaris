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
)

/* BEGIN EXPORTED METHODS */

// CreateDirIfDoesNotExit creates a given directory if it does not already exist.
func CreateDirIfDoesNotExit(dir string) error {
	safeDir := filepath.FromSlash(dir) // Just to be safe

	if _, err := os.Stat(safeDir); os.IsNotExist(err) { // Check dir exists
		err = os.MkdirAll(safeDir, 0755) // Create directory

		if err != nil { // Check for errors
			return err // Return error
		}
	}

	return nil // No error occurred
}

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

// getDataDir fetches the data direcotry
func getDataDir() string {
	abs, _ := filepath.Abs("./data") // Get absolute dir

	return filepath.FromSlash(abs) // Match OS slashes
}

/* END INTERNAL METHODS */
