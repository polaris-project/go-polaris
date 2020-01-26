// Package accounts implements types and methods for aiding in the generation and serialization of polaris accounts.
package accounts

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/polaris-project/go-polaris/crypto"
)

// Account represents a Polaris network primtiive
type Account struct {
	// privateKey is the private key of the account
	privateKey ed25519.PrivateKey
}

// NewAccount generates a new Account, along side a new private key.
func NewAccount() (Account, error) {
	// Generate a private key for the account
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)

	// Check for any errors that might arise whilst generating the key
	if err != nil {
		// Return the errors
		return Account{}, err
	}

	// Generate and return an account from the crypto lib rand reader
	return Account{
		privateKey: privateKey,
	}, nil
}

// Address derives a Polaris address from the account's ed25519 private key.
func (acc *Account) Address() crypto.Address {
	// Return a hash of the account's public key
	return crypto.AddressFromPrivateKey(&acc.privateKey)
}
