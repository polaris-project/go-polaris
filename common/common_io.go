// Package common defines a set of commonly used helper methods and data types.
package common

import (
	"fmt"
	"path/filepath"
)

var (
	// DataDir is the global data directory definition.
	DataDir = getDataDir()

	// ConfigDir is the global configuration directory definition.
	ConfigDir = filepath.FromSlash(fmt.Sprintf("%s/config", DataDir))
)

/* BEGIN INTERNAL METHODS */

// getDataDir fetches the data direcotry
func getDataDir() string {
	abs, _ := filepath.Abs("./data") // Get absolute dir

	return filepath.FromSlash(abs) // Match OS slashes
}

/* END INTERNAL METHODS */
