// Package common defines a set of commonly used helper methods and data types.
package common

import "path/filepath"

var (
	// DataDir is the global data directory definition.
	DataDir = getDataDir()
)

/* BEGIN INTERNAL METHODS */

// getDataDir fetches the data direcotry
func getDataDir() string {
	abs, _ := filepath.Abs("./data") // Get absolute dir

	return filepath.FromSlash(abs) // Match OS slashes
}

/* END INTERNAL METHODS */
