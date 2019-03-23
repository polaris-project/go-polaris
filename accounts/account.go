// Package accounts defines a set of ECDSA private-public keypair management utilities and helper methods.
package accounts

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"math/big"

	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/crypto"
)

// Account represents an ECDSA private-public keypair.
// Only an account's x and y curve values are stored persistently.
// ecdsa.PrivateKey and ecdsa.PublicKey references can be obtained at runtime via .PrivateKey() and .PublicKey().
type Account struct {
	X *big.Int `json:"x"` // X value
	Y *big.Int `json:"y"` // Y value
	D *big.Int `json:"d"` // D value
}

/* BEGIN EXPORTED METHODS */

// NewAccount generates a new ECDSA private-public keypair, returns the initialized account.
// Does not write the new account to persistent memory on creation.
func NewAccount() (*Account, error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		return &Account{}, err // Return found error
	}

	return &Account{
		privateKey.X, // Set X
		privateKey.Y, // Set Y
		privateKey.D, // Set D
	}, nil // Return initialized account
}

// AccountFromKey initializes a new account instance from a given ECDSA private key.
func AccountFromKey(privateKey *ecdsa.PrivateKey) *Account {
	return &Account{
		X: privateKey.X, // Set X
		Y: privateKey.Y, // Set Y
	} // Return initialized account
}

// Address attempts to derive an address from the given account.
func (account *Account) Address() *common.Address {
	publicKey := &ecdsa.PublicKey{
		X: account.X, // Set X
		Y: account.Y, // Set Y
	}

	return crypto.AddressFromPublicKey(publicKey) // Return address value
}

/* END EXPORTED METHODS */
