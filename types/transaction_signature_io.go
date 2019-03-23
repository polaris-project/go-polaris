// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"encoding/json"
)

/* BEGIN EXPORTED METHODS */

// String marshals the contents of a given ECDA signature via JSON.
func (signature *Signature) String() string {
	marshaledVal, _ := json.MarshalIndent(*signature, "", "  ") // Marshal indent

	return string(marshaledVal) // Return JSON string
}

/* END EXPORTED METHODS */
