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

	"github.com/juju/loggo"

	"github.com/boltdb/bolt"
	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/config"
	"github.com/polaris-project/go-polaris/crypto"
)

var transactionBucket = []byte("transaction-bucket")

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

	// ErrNilGenesis represents an error describing a genesis value of nil.
	ErrNilGenesis = errors.New("dag does not have a valid genesis")

	// ErrNilSignature represents an error describing a transaction lacking a signature.
	ErrNilSignature = errors.New("transaction has no signature")

	// ErrDuplicateTransaction represents an error describing a duplicate transaction entry in the given working dag db.
	ErrDuplicateTransaction = errors.New("transaction already exists in dag")

	// ErrInvalidSignature represents an error describing an
	ErrInvalidSignature = errors.New("signature invalid")
)

// logger is the dag package logger.
var logger = getDagLogger()

// Dag is a simple struct used to abstract db reading and writing methods.
type Dag struct {
	DagConfig *config.DagConfig `json:"config"` // Dag config

	Genesis common.Hash `json:"genesis"` // Dag genesis

	LastTransaction common.Hash `json:"last_tx"` // Last transaction hash
}

/* BEGIN EXPORTED METHODS */

// NewDag creates a new dag with the given config, and writes the dag db to memory.
// The newly opened dag db is stored in the WorkingDagDB variable.
func NewDag(config *config.DagConfig) (*Dag, error) {
	logger.Infof("initializing dag instance") // Log init dag

	err := config.WriteToMemory() // Write dag config to persistent memory
	if err != nil {               // Check for errors
		return &Dag{}, err // Return found error
	}

	err = common.CreateDirIfDoesNotExist(common.DbDir) // Make database directory

	if err != nil { // Check for errors
		return &Dag{}, err // Return found error
	}

	logger.Infof("opening dag db") // Log open db

	dagDB, err := bolt.Open(filepath.FromSlash(fmt.Sprintf("%s/%s.db", common.DbDir, config.Identifier)), 0o644, &bolt.Options{Timeout: 5 * time.Second}) // Open DB with timeout
	if err != nil {                                                                                                                                       // Check for errors
		return &Dag{}, err // Return found error
	}

	WorkingDagDB = dagDB // Set dag DB

	err = createTransactionBucketIfNotExist() // Create transaction bucket if it doesn't already exist

	if err != nil { // Check for errors
		return &Dag{}, err // Return found error
	}

	logger.Infof("attempting to open dag db header") // Log open dag db header

	dagHeader, err := readDagDbHeaderFromMemory(config.Identifier) // Read dag db

	if err != nil || dagHeader == nil { // Check no existing dag
		logger.Infof("could not load local dag db header; initializing one instead") // Log initialize

		dagHeader = &Dag{
			DagConfig: config, // Set config
		} // Initialize dag db header

		logger.Infof("initialized dag db header, writing to memory") // Log write

		err = dagHeader.WriteToMemory() // Write dag db header to persistent memory

		if err != nil { // Check for errors
			return &Dag{}, err // Return found error
		}
	}

	logger.Infof("finished setting up dag") // Log setup dag

	return dagHeader, nil // Return initialized dag
}

// Close closes the working dag.
func (dag *Dag) Close() error {
	logger.Infof("closing dag db") // Log close

	if WorkingDagDB == nil { // Check no working dag db
		return ErrDagDbNotOpened // Return error
	}

	return WorkingDagDB.Close() // Close
}

// MakeGenesis makes the dag's genesis transaction set.
// If the dag already has a genesis transaction, an ErrDuplicateTransaction error is returned.
func (dag *Dag) MakeGenesis() ([]*Transaction, error) {
	logger.Infof("making genesis transaction set") // Log make genesis

	if !dag.Genesis.IsNil() { // Check genesis already exists
		return nil, ErrDuplicateTransaction // Return found error
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader) // Generate gensis address
	if err != nil {                                                    // Check for errors
		return nil, err // Return found error
	}

	logger.Infof("created genesis private key") // Log init private key

	totalGenesisValue := 0.0 // Init total value buffer

	for _, value := range dag.DagConfig.Alloc { // Iterate through alloc
		totalGenesisValue += value // Increment value
	}

	genesisTransactions := []*Transaction{} // Initialize genesis transactions

	logger.Infof("creating genesis transaction") // Log init genesis

	genesisTransaction := NewTransaction(0, big.NewFloat(totalGenesisValue), nil, crypto.AddressFromPrivateKey(privateKey), nil, 0, big.NewInt(0), []byte("genesis")) // Initialize genesis transaction

	err = dag.forceAddTransaction(genesisTransaction) // Add genesis transaction

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	logger.Infof("added genesis transaction to dag") // Log add genesis

	(*dag).Genesis = genesisTransaction.Hash // Set genesis

	err = (*dag).WriteToMemory() // Write dag header to persistent memory

	if err != nil { // Check for errors
		return nil, err // Return found error
	}

	genesisTransactions = append(genesisTransactions, genesisTransaction) // Append genesis

	lastParent := genesisTransaction // Set last parent

	x := uint64(0) // Init nonce

	for key, value := range dag.DagConfig.Alloc { // Iterate through alloc
		logger.Infof("creating genesis child transaction for allocation address: %s", key) // Log create genesis child

		decodedKey, err := hex.DecodeString(key) // Decode key
		if err != nil {                          // Check for errors
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

		logger.Infof("added genesis child transaction with hash: %s and alloc address: %s", hex.EncodeToString(transaction.Hash.Bytes()), key) // Log add genesis child

		genesisTransactions = append(genesisTransactions, transaction) // Append transaction

		lastParent = transaction // Set last parent

		x++ // Increment nonce
	}

	return genesisTransactions, nil // No error occurred, return nil
}

// OpenDag attempts to open all dag-related resources.
func OpenDag(identifier string) (*Dag, error) {
	logger.Infof("opening dag db header with identifier: %s", identifier) // Log open dag

	dagDbHeader, err := readDagDbHeaderFromMemory(identifier) // Read dag db header
	if err != nil {                                           // Check for errors
		return &Dag{}, err // Return found error
	}

	logger.Infof("finished opening dag db header with identifier: %s", identifier) // Log opened dag

	WorkingDagDB, err = bolt.Open(filepath.FromSlash(fmt.Sprintf("%s/%s.db", common.DbDir, identifier)), 0o644, &bolt.Options{Timeout: 5 * time.Second}) // Open DB with timeout

	if err != nil { // Check for errors
		return &Dag{}, err // Return found error
	}

	logger.Infof("opened dag db with identifier: %s", identifier) // Log opened dag db

	return dagDbHeader, nil // Return dag db header
}

// AddTransaction appends a given transaction to the working dag.
// Returns an ErrDagDbNotOpened error if the working dag db is nil (has been not opened).
// Return an ErrNilTransaction error if the given transaction pointer is nil.
// Returns an ErrDuplicateTransaction error if the transaction already exists in the working dag db.
// Returns an ErrNilSignature error if the transaction has not been signed.
// Return an ErrInvalidSignature error if the transaction's signature is invalid.
func (dag *Dag) AddTransaction(transaction *Transaction) error {
	logger.Infof("adding transaction with hash: %s", hex.EncodeToString(transaction.Hash.Bytes())) // Log add transaction

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
	if err != nil {                            // Check for errors
		return err // Return found error
	}

	logger.Infof("checking transaction with hash: %s already exists in dag", hex.EncodeToString(transaction.Hash.Bytes())) // Log check tx already exists

	_, err = dag.GetTransactionByHash(transaction.Hash) // Get transaction by hash

	if err == nil { // Check tx already exists
		return ErrDuplicateTransaction // Return found error
	}

	logger.Infof("verifying transaction signature with hash: %s", hex.EncodeToString(transaction.Hash.Bytes())) // Log verify tx signature

	if !transaction.Signature.Verify(transaction.Sender) { // Check transaction signature invalid
		return ErrInvalidSignature // Return found error
	}

	logger.Infof("transaction signature with hash: %s verified", hex.EncodeToString(transaction.Hash.Bytes())) // Log verified signature

	return WorkingDagDB.Update(func(tx *bolt.Tx) error {
		workingTransactionBucket := tx.Bucket(transactionBucket) // Get transaction bucket

		logger.Infof("adding transaction with hash: %s to dag db", hex.EncodeToString(transaction.Hash.Bytes())) // Log add tx to dag db

		return workingTransactionBucket.Put(transaction.Hash.Bytes(), transaction.Bytes()) // Put transaction
	}) // Write transaction
}

/*
	BEGIN DB READING HELPER METHODS
*/

// GetTransactionByHash attempts to query the working dag db by the given transaction hash.
// If no transaction exists at this hash, an
func (dag *Dag) GetTransactionByHash(transactionHash common.Hash) (*Transaction, error) {
	logger.Infof("attempting to query transaction by hash: %s", hex.EncodeToString(transactionHash.Bytes())) // Log get tx

	var txBytes []byte // Init buffer

	if WorkingDagDB == nil { // Check no working db
		return &Transaction{}, ErrDagDbNotOpened // Return found error
	}

	err := createTransactionBucketIfNotExist() // Create tx bucket if doesn't already exist to prevent nil pointer dereferences
	if err != nil {                            // Check for errors
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
	logger.Infof("attempting to query transaction children for tx with hash: %s", hex.EncodeToString(transactionHash.Bytes())) // Log query tx children

	if WorkingDagDB == nil { // Check no dag db
		return []*Transaction{}, ErrDagDbNotOpened // Return found error
	}

	transactions := []*Transaction{} // Initialize tx buffer

	err := createTransactionBucketIfNotExist() // Create transaction bucket if not exist
	if err != nil {                            // Check for errors
		return []*Transaction{}, err // Return found error
	}

	return transactions, WorkingDagDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(transactionBucket) // Get transaction bucket

		c := bucket.Cursor() // Get cursor

		for currentTransactionHash, transactionBytes := c.First(); currentTransactionHash != nil; currentTransactionHash, transactionBytes = c.Next() { // Iterate through tx set
			transaction := TransactionFromBytes(transactionBytes) // Deserialize transaction

			for _, parentHash := range transaction.ParentTransactions { // Iterate through tx parents
				if bytes.Equal(parentHash.Bytes(), transactionHash.Bytes()) { // Check matching parent
					logger.Infof("found child with hash: %s and parent hash: %s", hex.EncodeToString(transaction.Hash.Bytes()), hex.EncodeToString(transactionHash.Bytes())) // Log found child

					transactions = append(transactions, transaction) // Append to transactions
				}
			}
		}

		return nil // No error occurred, return nil
	}) // Return filtered transactions
}

// GetTransactionsByAddress attempts to filter the dag by a given sending or receiving address.
func (dag *Dag) GetTransactionsByAddress(address *common.Address) ([]*Transaction, error) {
	logger.Infof("attempting to query transactions by sender or recipient: %s", hex.EncodeToString(address.Bytes())) // Log query tx

	if WorkingDagDB == nil { // Check no dag db
		return []*Transaction{}, ErrDagDbNotOpened // Return found error
	}

	transactions := []*Transaction{} // Init tx buffer

	err := createTransactionBucketIfNotExist() // Create transaction bucket if not exist
	if err != nil {                            // Check for errors
		return []*Transaction{}, err // Return found error
	}

	return transactions, WorkingDagDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(transactionBucket) // Get transaction bucket

		c := bucket.Cursor() // Get cursor

		for transactionHash, transactionBytes := c.First(); transactionHash != nil; transactionHash, transactionBytes = c.Next() { // Iterate through tx set
			transaction := TransactionFromBytes(transactionBytes) // Deserialize transaction

			if bytes.Equal(transaction.Sender.Bytes(), address.Bytes()) || bytes.Equal(transaction.Recipient.Bytes(), address.Bytes()) { // Check relevant
				logger.Infof("found transaction with recipient/sender: %s", hex.EncodeToString(transaction.Hash.Bytes())) // Log found tx

				transactions = append(transactions, transaction) // Append transaction
			}
		}

		return nil // No error occurred, return nil
	}) // Return filtered transactions
}

// GetTransactionsBySender attempts to filter the dag by a given sending address.
func (dag *Dag) GetTransactionsBySender(sender *common.Address) ([]*Transaction, error) {
	logger.Infof("attempting to query transactions by sender: %s", hex.EncodeToString(sender.Bytes())) // Log query tx

	if WorkingDagDB == nil { // Check no dag db
		return []*Transaction{}, ErrDagDbNotOpened // Return found error
	}

	transactions := []*Transaction{} // Init tx buffer

	err := createTransactionBucketIfNotExist() // Create transaction bucket if not exist
	if err != nil {                            // Check for errors
		return []*Transaction{}, err // Return found error
	}

	return transactions, WorkingDagDB.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(transactionBucket) // Get transaction bucket

		c := bucket.Cursor() // Get cursor

		for transactionHash, transactionBytes := c.First(); transactionHash != nil; transactionHash, transactionBytes = c.Next() { // Iterate through tx set
			transaction := TransactionFromBytes(transactionBytes) // Deserialize transaction

			if bytes.Equal(transaction.Sender.Bytes(), sender.Bytes()) { // Check is sender
				logger.Infof("found transaction with sender: %s", hex.EncodeToString(transaction.Hash.Bytes())) // Log found tx

				transactions = append(transactions, transaction) // Append transaction
			}
		}

		return nil // No error occurred, return nil
	}) // Return filtered transactions
}

// GetBestTransaction gets the last transaction in the dag. If more than one last child exists, the child with
// the latest timestamp is returned.
func (dag *Dag) GetBestTransaction() (*Transaction, error) {
	logger.Infof("attempting to find best transaction in working dag") // Log query tx

	if dag.Genesis.IsNil() { // Check no genesis
		return &Transaction{}, nil // No best tx
	}

	if dag.LastTransaction.IsNil() { // Check no best transaction
		logger.Infof("dag db header does not have last tx; setting to genesis") // Log set last tx

		dag.LastTransaction = dag.Genesis // Set last transaction
	}

	err := dag.WriteToMemory() // Write to persistent memory
	if err != nil {            // Check for errors
		return &Transaction{}, err // Return found error
	}

	lastTransaction, err := dag.GetTransactionByHash(dag.LastTransaction) // Initialize last transaction buffer
	if err != nil {                                                       // Check for errors
		return &Transaction{}, err // Return found error
	}

	logger.Infof("starting get best job from root tx: %s", hex.EncodeToString(lastTransaction.Hash.Bytes())) // Log starting tx

	for { // Do until found parent without children
		children, _ := dag.GetTransactionChildren(lastTransaction.Hash) // Get children

		if len(children) == 0 { // Check no children
			break // Break
		}

		bestTransaction := children[0] // Get first tx pointer

		for _, child := range children { // Iterate through children
			if child.Timestamp.After(bestTransaction.Timestamp) { // Check latest child
				bestTransaction = child // Set best transaction
			}
		}

		lastTransaction = bestTransaction // Set best transaction
	}

	dag.LastTransaction = lastTransaction.Hash // Set last transaction

	logger.Infof("found last transaction: %s", hex.EncodeToString(lastTransaction.Hash.Bytes())) // Log found best tx

	return lastTransaction, dag.WriteToMemory() // Return last transaction
}

/*
	END DB READING HELPER METHODS
*/

/*
	BEGIN HELPER METHODS
*/

// CalculateAddressBalance calculates the total balance of an address from genesis to latest tx.
func (dag *Dag) CalculateAddressBalance(address *common.Address) (*big.Float, error) {
	logger.Infof("calculating balance for address: %s", hex.EncodeToString(address.Bytes())) // Log calculate balance

	transactionsRegardingAddress, err := dag.GetTransactionsByAddress(address) // Filter by pertaining to
	if err != nil {                                                            // Check for errors
		return &big.Float{}, err // Return found error
	}

	logger.Infof("found %d transactions related to address %s", len(transactionsRegardingAddress), hex.EncodeToString(address.Bytes())) // Log found related

	balance := big.NewFloat(0) // Init balance buffer

	for _, transaction := range transactionsRegardingAddress { // Iterate through transactions
		if bytes.Equal(transaction.Sender.Bytes(), address.Bytes()) { // Check was sender
			balance.Sub(balance, transaction.CalculateTotalValue()) // Subtract transaction value
		}

		if bytes.Equal(transaction.Recipient.Bytes(), address.Bytes()) { // Check was recipient
			balance.Add(balance, transaction.Amount) // Add transaction amount
		}
	}

	logger.Infof("calculated balance of address %s: %d", hex.EncodeToString(address.Bytes()), balance) // Log calculated balance

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
	if err != nil {                            // Check for errors
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

// getDagLogger gets the dag package logger, and sets the levels of said logger.
func getDagLogger() loggo.Logger {
	logger := loggo.GetLogger("dag") // Get logger

	loggo.ConfigureLoggers("dag=INFO") // Configure loggers

	return logger // Return logger
}

/*
	END DB BUCKET HELPER METHODS
*/

/* END INTERNAL METHODS */
