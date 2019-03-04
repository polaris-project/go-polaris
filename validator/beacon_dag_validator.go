// Package validator represents a helper methods useful for validators in the Polaris network.
// Methods in the validator package are specified in terms of a validator interface, that of which is
// also implemented in the validator package.
package validator

import "github.com/polaris-project/go-polaris/config"

// BeaconDagValidator represents a main dag validator.
type BeaconDagValidator struct {
	Config *config.DagConfig `json:"config"` // Config represents the beacon dag config
}

/* BEGIN EXPORTED METHODS */

/* END EXPORTED METHODS */
