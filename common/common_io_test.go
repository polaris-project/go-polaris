package common

import "testing"

/* BEGIN EXPORTED METHODS TESTS */

// TestCreateDirIfDoesNotExist tests the functionality of the CreateDirIfDoesNotExit() helper method.
func TestCreateDirIfDoesNotExist(t *testing.T) {
	err := CreateDirIfDoesNotExit(DataDir) // Create data dir (just as an example)

	if err != nil { // Check for errors
		t.Fatal(err) // Panic
	}
}

/* END EXPORTED METHODS TESTS */

/* BEGIN INTERNAL METHODS TESTS */

// TestGetDataDir tests the functionality of the getDataDir() helper method.
func TestGetDataDir(t *testing.T) {
	t.Log(getDataDir()) // Log data dir
}

/* END INTERNAL METHODS TESTS */
