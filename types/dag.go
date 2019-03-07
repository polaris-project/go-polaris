// Package types provides core primitives for the operation
// of the Polaris protocol.
package types

import (
	"bytes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"path/filepath"
	"time"

	"github.com/boltdb/bolt"
	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/config"
	"github.com/polaris-project/go-polaris/crypto"
)

var (
	transactionBucket = []byte("transaction-bucket")
)

var (
	// WorkingDagDB represents the current opened dag database.
	WorkingDagDB *bolt.DB

	// ErrDagAlreadyExists represents an error describing
	// the attempted overwriting of an existing DAG.
	ErrDagAlreadyExists = errors.New("dag already exists")

	// ErrDagDbNotOpened represents an error describing the attempted appending of a data set
	// to the working dag, despite the fact that the dag db has not yet been opened.
	ErrDagDbNotOpened = errors.New("dag db has not been opened")

	// ErrNilTransaction represents an error describing a transaction pointer of nil value.
	ErrNilTransaction = errors.New("transaction pointer is nil")

	// ErrNilTransactionAtHash represents an error describing a transaction pointer of nil value discovered through the
	// querying of the working db for an invalid hash.
	ErrNilTransactionAtHash = errors.New("no transaction exists in the dag with given hash")

	// ErrNilSignature represents an error describing a transaction lacking a signature.
	ErrNilSignature = errors.New("transaction has no signature")

	// ErrDuplicateTransaction represents an error describing a duplicate transaction entry in the given working dag db.
	ErrDuplicateTransaction = errors.New("transaction already exists in dag")

	// ErrInvalidSignature represents an error describing an
	ErrInvalidSignature = errors.New("signature invalid")
)

// Dag is a simple struct used to abstract db reading and writing methods.
type Dag struct {
	DagConfig *config.DagConfig `json:"config"` // Dag config

	Genesis common.Hash `json:"genesis"` // Dag genesis
}

/* BEGIN EXPORTED METHODS */

// NewDag creates a new dag with the given config, and writes the dag db to memory.
// The newly opened dag db is stored in the WorkingDagDB variable.
func NewDag(config *config.DagConfig) (*Dag, error) {
	err := config.WriteToMemory() // Write dag config to persistent memory

	if err != nil { // Check for errors
		return &Dag{}, err // Return found error
	}

	err = common.CreateDirIfDoesNotExit(common.DbDir) // Make database directory

	if err != nil { // Check for errors
		return &Dag{}, err // Return found error
	}

	dagDB, err := bolt.Open(filepath.FromSlash(fmt.Sprintf("%s/%s.db", common.DbDir, config.Identifier)), 0644, &bolt.Options{Timeout: 5 * time.Second}) // Open DB with timeout

	if err != nil { // Check for errors
		return &Dag{}, err // Return found error
	}

	WorkingDagDB = dagDB // Set dag DB

	err = createTransactionBucketIfNotExist() // Create transaction bucket if it doesn't already exist

	if err != nil { // Check for errors
		return &Dag{}, err // Return found error
	}

	dag := &Dag{
		DagConfig: config, // Set config
	} // Initialize dag db header

	err = dag.writeToMemory() // Write dag db header to persistent memory

	if err != nil { // Check for errors
		return &Dag{}, err // Return found error
	}

	return dag, nil // Return initialized dag
}

// MakeGenesis makes the dag's genesis transaction set.
// If the dag already has a genesis transaction, an ErrDuplicateTransaction error is returned.
func (dag *Dag) MakeGenesis() ([]*Transaction, error) {
	if !dag.Genesis.IsNil() { // Check genesis already exists
		return nil, ErrDuplicateTransaction // Return found error
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate gensis address

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	totalGenesisValue := 0.0 // Init total value buffer

	for _, value := range dag.DagConfig.Alloc { // Iterate through alloc
		totalGenesisValue += value // Increment value
	}

	genesisTransactions := []*Transaction{} // Initialize genesis transactions

	genesisTransaction := NewTransaction(0, big.NewFloat(totalGenesisValue), nil, crypto.AddressFromPrivateKey(privateKey), nil, 0, big.NewInt(0), []byte("genesis")) // Initialize genesis transaction

	err = dag.forceAddTransaction(genesisTransaction) // Add genesis transaction

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	genesisTransactions = append(genesisTransactions, genesisTransaction) // Append genesis

	lastParent := genesisTransaction // Set last parent

	x := uint64(1) // Init nonce

	for key, value := range dag.DagConfig.Alloc { // Iterate through alloc
		decodedKey, err := hex.DecodeString(key) // Decode key

		if err != nil { // Check for errors
			return nil, err // Return found error
		}

		decodedAddress := common.NewAddress(decodedKey) // Decode address

		transaction := NewTransaction(x, big.NewFloat(value), crypto.AddressFromPrivateKey(privateKey), decodedAddress, []common.Hash{lastParent.Hash}, 0, big.NewInt(0), []byte("genesis_child")) // Initialize new genesis child transaction

		err = SignTransaction(transaction, privateKey) // Sign transaction

		if err != nil { // Check for errors
			return nil, err // Return found error
		}

		err = dag.AddTransaction(transaction) // Add transaction

		if err != nil { // Check for errors
			return nil, err // Return found error
		}

		genesisTransactions = append(genesisTransactions, transaction) // Append transaction

		lastParent = transaction // Set last parent

		x++ // Increment nonce
	}

	return nil, nil // No error occurred, return nil
}

// OpenDag attempts to open all dag-related resources.
func OpenDag(identifier string) (*Dag, error) {
	dagDbHeader, err := readDagDbHeaderFromMemory(identifier) // Read dag db header

	if err != nil { // Check for errors
		return &Dag{}, err // Return found error
	}

	WorkingDagDB, err = bolt.Open(filepath.FromSlash(fmt.Sprintf("%s/%s.db", common.DbDir, identifier)), 0644, &bolt.Options{Timeout: 5 * time.Second}) // Open DB with timeout

	if err != nil { // Check for errors
		return &Dag{}, err // Return found error
	}

	return dagDbHeader, nil // Return dag db header
}

// AddTransaction appends a given transaction to the working dag.
// Returns an ErrDagDbNotOpened error if the working dag db is nil (has been not opened).
// Return an ErrNilTransaction error if the given transaction pointer is nil.
// Returns an ErrDuplicateTransaction error if the transaction already exists in the working dag db.
// Returns an ErrNilSignature error if the transaction has not been signed.
// Return an ErrInvalidSignature error if the transaction's signature is invalid.
func (dag *Dag) AddTransaction(transaction *Transaction) error {
	if WorkingDagDB == nil { // Check dag db not opened
		return ErrDagDbNotOpened // Return found error
	}

	if transaction == nil { // Check nil pointer
		return ErrNilTransaction // Return found error
	}

	if transaction.Signature == nil { // Check no signature
		return ErrNilSignature // Return found error
	}

	err := createTransactionBucketIfNotExist() // Create transaction bucket if it doesn't already exist

	if err != nil { // Check for errors
		return err // Return found error
	}

	_, err = dag.GetTransactionByHash(transaction.Hash) // Get transaction by hash

	if err == nil { // Check tx already exists
		return ErrDuplicateTransaction // Return found error
	}

	if !transaction.Signature.Verify(transaction.Sender) { // Check transaction signature invalid
		return ErrInvalidSignature // Return found error
	}

	return WorkingDagDB.Update(func(tx *bolt.Tx) error {
		workingTransactionBucket := tx.Bucket(transactionBucket) // Get transaction bucket

		return workingTransactionBucket.Put(transaction.Hash.Bytes(), transaction.Bytes()) // Put transaction
	}) // Write transaction
}

/*
	BEGIN DB READING HELPER METHODS
*/

// GetTransactionByHash attempts to query the working dag db by the given transaction hash.
// If no transaction exists at this hash, an
func (dag *Dag) GetTransactionByHash(transactionHash common.Hash) (*Transaction, error) {
	var txBytes []byte // Init buffer

	if WorkingDagDB == nil { // Check no working db
		return &Transaction{}, ErrDagDbNotOpened // Return found error
	}

	err := createTransactionBucketIfNotExist() // Create tx bucket if doesn't already exist to prevent nil pointer dereferences

	if err != nil { // Check for errors
		return &Transaction{}, err // Return found error
	}

	err = WorkingDagDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(transactionBucket) // Get tx bucket

		txBytes = bucket.Get(transactionHash.Bytes()) // Get tx at hash

		if txBytes == nil { // Check no transaction at hash
			return ErrNilTransactionAtHash // Return error
		}

		return nil // No error occurred, return nil
	})

	if err != nil { // Check for errors
		return &Transaction{}, err // Return found error
	}

	return TransactionFromBytes(txBytes), nil // Return deserialized tx
}

// GetTransactionChildren iterates through the dag's transactions, and finds transactions with the given hash as a parent.
func (dag *Dag) GetTransactionChildren(transactionHash common.Hash) ([]*Transaction, error) {
	if WorkingDagDB == nil { // Check no dag db
		return []*Transaction{}, ErrDagDbNotOpened // Return found error
	}

	transactions := []*Transaction{} // Initialize tx buffer

	err := createTransactionBucketIfNotExist() // Create transaction bucket if not exist

	if err != nil { // Check for errors
		return []*Transaction{}, err // Return found error
	}

	return transactions, WorkingDagDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(transactionBucket) // Get transaction bucket

		c := bucket.Cursor() // Get cursor

		for currentTransactionHash, transactionBytes := c.First(); currentTransactionHash != nil; currentTransactionHash, transactionBytes = c.Next() { // Iterate through tx set
			transaction := TransactionFromBytes(transactionBytes) // Deserialize transaction

			for _, parentHash := range transaction.ParentTransactions { // Iterate through tx parents
				if bytes.Equal(parentHash.Bytes(), transactionHash.Bytes()) { // Check matching parent
					transactions = append(transactions, transaction) // Append to transactions
				}
			}
		}

		return nil // No error occurred, return nil
	}) // Return filtered transactions
}

// GetTransactionsByAddress attempts to filter the dag by a given sending or receiving address.
func (dag *Dag) GetTransactionsByAddress(address *common.Address) ([]*Transaction, error) {
	if WorkingDagDB == nil { // Check no dag db
		return []*Transaction{}, ErrDagDbNotOpened // Return found error
	}

	transactions := []*Transaction{} // Init tx buffer

	err := createTransactionBucketIfNotExist() // Create transaction bucket if not exist

	if err != nil { // Check for errors
		return []*Transaction{}, err // Return found error
	}

	return transactions, WorkingDagDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(transactionBucket) // Get transaction bucket

		c := bucket.Cursor() // Get cursor

		for transactionHash, transactionBytes := c.First(); transactionHash != nil; transactionHash, transactionBytes = c.Next() { // Iterate through tx set
			transaction := TransactionFromBytes(transactionBytes) // Deserialize transaction

			if bytes.Equal(transaction.Sender.Bytes(), address.Bytes()) || bytes.Equal(transaction.Recipient.Bytes(), address.Bytes()) { // Check relevant
				transactions = append(transactions, transaction) // Append transaction
			}
		}

		return nil // No error occurred, return nil
	}) // Return filtered transactions
}

// GetTransactionsBySender attempts to filter the dag by a given sending address.
func (dag *Dag) GetTransactionsBySender(sender *common.Address) ([]*Transaction, error) {
	if WorkingDagDB == nil { // Check no dag db
		return []*Transaction{}, ErrDagDbNotOpened // Return found error
	}

	transactions := []*Transaction{} // Init tx buffer

	err := createTransactionBucketIfNotExist() // Create transaction bucket if not exist

	if err != nil { // Check for errors
		return []*Transaction{}, err // Return found error
	}

	return transactions, WorkingDagDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(transactionBucket) // Get transaction bucket

		c := bucket.Cursor() // Get cursor

		for transactionHash, transactionBytes := c.First(); transactionHash != nil; transactionHash, transactionBytes = c.Next() { // Iterate through tx set
			transaction := TransactionFromBytes(transactionBytes) // Deserialize transaction

			if bytes.Equal(transaction.Sender.Bytes(), sender.Bytes()) { // Check is sender
				transactions = append(transactions, transaction) // Append transaction
			}
		}

		return nil // No error occurred, return nil
	}) // Return filtered transactions
}

/*
	END DB READING HELPER METHODS
*/

/*
	BEGIN HELPER METHODS
*/

// CalculateAddressBalance calculates the total balance of an address from genesis to latest tx.
func (dag *Dag) CalculateAddressBalance(address *common.Address) (*big.Float, error) {
	transactionsRegardingAddress, err := dag.GetTransactionsByAddress(address) // Filter by pertaining to

	if err != nil { // Check for errors
		return &big.Float{}, err // Return found error
	}

	balance := big.NewFloat(0) // Init balance buffer

	for _, transaction := range transactionsRegardingAddress { // Iterate through transactions
		if bytes.Equal(transaction.Sender.Bytes(), address.Bytes()) { // Check was sender
			balance.Sub(balance, transaction.CalculateTotalValue()) // Subtract transaction value
		}

		if bytes.Equal(transaction.Recipient.Bytes(), address.Bytes()) { // Check was recipient
			balance.Add(balance, transaction.Amount) // Add transaction amount
		}
	}

	return balance, nil // Return balance
}

/*
	END HELPER METHODS
*/

/* END EXPORTED METHODS */

/* BEGIN INTERNAL METHODS */

/*
	BEGIN DB BUCKET HELPER METHODS
*/

// forceAddTransaction forces the adding of a given transaction to the dag (only useful for adding a genesis tx).
func (dag *Dag) forceAddTransaction(transaction *Transaction) error {
	err := createTransactionBucketIfNotExist() // Create transaction bucket if it doesn't already exist

	if err != nil { // Check for errors
		return err // Return found error
	}

	_, err = dag.GetTransactionByHash(transaction.Hash) // Get transaction by hash

	if err == nil { // Check tx already exists
		return ErrDuplicateTransaction // Return found error
	}

	return WorkingDagDB.Update(func(tx *bolt.Tx) error {
		workingTransactionBucket := tx.Bucket(transactionBucket) // Get transaction bucket

		return workingTransactionBucket.Put(transaction.Hash.Bytes(), transaction.Bytes()) // Put transaction
	}) // Write transaction // No error occurred, return nil
}

// createTransactionBucketIfNotExist attempts to create the "transaction" bucket in the working dag db.
func createTransactionBucketIfNotExist() error {
	if WorkingDagDB == nil { // Check no working db
		return ErrDagDbNotOpened // Return found error
	}

	err := WorkingDagDB.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(transactionBucket) // Create tx bucket if it doesn't already exist

		return err // Return error
	}) // Create tx bucket if it doesn't already exist

	return err // Return error
}

/*
	END DB BUCKET HELPER METHODS
*/

/* END INTERNAL METHODS */
