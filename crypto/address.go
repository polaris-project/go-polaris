// Package crypto provides cryptography helper methods.
package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"

	"github.com/polaris-project/go-polaris/common"
)

/* BEGIN EXPORTED METHODS */

// AddressFromPrivateKey serializes and converts an ecdsa private key into an address (uses private key pub).
func AddressFromPrivateKey(privateKey *ecdsa.PrivateKey) *common.Address {
	return AddressFromPublicKey(&privateKey.PublicKey) // Return address
}

// AddressFromPublicKey serializes and converts an ecdsa public key into an address.
func AddressFromPublicKey(publicKey *ecdsa.PublicKey) *common.Address {
	publicKeyBytes := elliptic.Marshal(elliptic.P521(), publicKey.X, publicKey.Y) // Marshal public key

	return common.NewAddress(Sha3(publicKeyBytes).Bytes()) // Return address value
}

/* END EXPORTED METHODS */
