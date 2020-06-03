// Package accounts defines a set of ECDSA private-public keypair management utilities and helper methods.
package accounts

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/crypto"
)

/* BEGIN EXPORTED METHODS */

// GetAllAccounts gets a list of all the owned accounts.
func GetAllAccounts() []common.Address {
	buffer := []common.Address{} // Init buffer

	files, err := ioutil.ReadDir(common.KeystoreDir) // Walk keystore dir

	if err != nil { // Check for errors
		return []common.Address{} // Return nil
	}

	for _, file := range files { // Iterate through files
		addressBytes, _ := hex.DecodeString(strings.Split(strings.Split(file.Name(), "account_")[1], ".json")[0]) // Decode hex addr

		buffer = append(buffer, *common.NewAddress(addressBytes)) // Append address
	}

	return buffer // No error occurred, return success
}

// String marshals a given account's contents to a JSON-encoded string.
func (account *Account) String() string {
	marshaledVal, _ := json.MarshalIndent(*account, "", "  ") // Marshal JSON

	return string(marshaledVal) // Return JSON
}

// Bytes encodes a given account's contents to a JSON-encoded byte array.
func (account *Account) Bytes() []byte {
	marshaledVal, _ := json.MarshalIndent(*account, "", "  ") // Marshal JSON

	return marshaledVal // Return JSON
}

// WriteToMemory writes the given account's contents to persistent memory
func (account *Account) WriteToMemory() error {
	err := common.CreateDirIfDoesNotExist(common.KeystoreDir) // Create keystore dir if necessary

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = ioutil.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/account_%s.json", common.KeystoreDir, hex.EncodeToString(crypto.AddressFromPublicKey(account.PublicKey()).Bytes()))), account.Bytes(), 0644) // Write account to persistent memory

	if err != nil { // Check for errors
		return err // Return error
	}

	return nil // No error occurred, return nil
}

// ReadAccountFromMemory reads an account with a given address from persistent memory.
func ReadAccountFromMemory(address *common.Address) (*Account, error) {
	data, err := ioutil.ReadFile(filepath.FromSlash(fmt.Sprintf("%s/account_%s.json", common.KeystoreDir, hex.EncodeToString(address.Bytes())))) // Read account

	if err != nil { // Check for errors
		return &Account{}, err // Return found error
	}

	buffer := &Account{} // Initialize buffer

	err = json.Unmarshal(data, buffer) // Deserialize JSON into buffer.

	if err != nil { // Check for errors
		return &Account{}, err // Return found error
	}

	return buffer, nil // No error occurred, return read account
}

/* END EXPORTED METHODS */
