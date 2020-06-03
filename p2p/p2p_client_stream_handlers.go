package p2p

import (
	"bufio"
	"encoding/hex"

	inet "github.com/libp2p/go-libp2p-net"
	protocol "github.com/libp2p/go-libp2p-protocol"
	"github.com/polaris-project/go-polaris/common"
	"github.com/polaris-project/go-polaris/types"
)

/* BEGIN EXPORTED METHODS */

/*
	BEGIN HANDLER REGISTRATION HELPERS
*/

// StartServingStreams attempts to start serving all necessary streams
func (client *Client) StartServingStreams(network string) error {
	logger.Infof("setting up node stream handlers") // Log set up stream handlers

	err := client.StartServingStream(GetStreamHeaderProtocolPath(network, PublishTransaction), client.HandleReceiveTransaction) // Register tx handler

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = client.StartServingStream(GetStreamHeaderProtocolPath(network, RequestConfig), client.HandleReceiveConfigRequest) // Register config request handler

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = client.StartServingStream(GetStreamHeaderProtocolPath(network, RequestBestTransaction), client.HandleReceiveBestTransactionRequest) // Register best tx request handler

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = client.StartServingStream(GetStreamHeaderProtocolPath(network, RequestTransaction), client.HandleReceiveTransactionRequest) // Register transaction request handler

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = client.StartServingStream(GetStreamHeaderProtocolPath(network, RequestGenesisHash), client.HandleReceiveGenesisHashRequest) // Register genesis hash request handler

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = client.StartServingStream(GetStreamHeaderProtocolPath(network, RequestChildHashes), client.HandleReceiveTransactionChildHashesRequest) // Register child hashes request handler

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// StartServingStream starts serving a stream on a given header protocol path.
func (client *Client) StartServingStream(streamHeaderProtocolPath string, handler func(inet.Stream)) error {
	if WorkingHost == nil { // Check no host
		return ErrNoWorkingHost // Return found error
	}

	WorkingHost.SetStreamHandler(protocol.ID(streamHeaderProtocolPath), handler) // Set handler

	return nil // No error occurred, return nil
}

/*
	END HANDLER REGISTRATION HELPERS
*/

/*
	BEGIN HANDLERS
*/

// HandleReceiveTransaction handles a new stream sending a transaction.
func (client *Client) HandleReceiveTransaction(stream inet.Stream) {
	logger.Infof("handling new publish transaction stream") // Log handle stream

	reader := bufio.NewReader(stream) // Initialize reader from stream

	transactionBytes, err := readAsync(reader) // Read async

	if err != nil { // Check for errors
		return // Return
	}

	transaction := types.TransactionFromBytes(transactionBytes) // Deserialize transaction

	logger.Infof("validating received transaction with hash: %s", hex.EncodeToString(transaction.Hash.Bytes())) // Log receive tx

	if err := (*client.Validator).ValidateTransaction(transaction); err == nil { // Check transaction valid
		logger.Infof("transaction was valid; adding to dag") // Log add to dag

		(*client.Validator).GetWorkingDag().AddTransaction(transaction) // Add transaction to working dag
	}
}

// HandleReceiveBestTransactionRequest handle a new stream requesting for the best transaction hash.
func (client *Client) HandleReceiveBestTransactionRequest(stream inet.Stream) {
	logger.Infof("handling new best transaction request stream") // Log handle stream

	writer := bufio.NewWriter(stream) // Initialize writer from stream

	defer writer.Flush() // Flush

	bestTransaction, _ := (*client.Validator).GetWorkingDag().GetBestTransaction() // Get best transaction

	logger.Infof("responding with best transaction hash %s", hex.EncodeToString(bestTransaction.Hash.Bytes())) // Log handle stream

	writer.Write(append(bestTransaction.Hash.Bytes(), byte('\f'))) // Write best transaction hash
}

// HandleReceiveTransactionRequest handles a new stream requesting transaction metadata with a given hash.
func (client *Client) HandleReceiveTransactionRequest(stream inet.Stream) {
	logger.Infof("handling new transaction with hash request stream") // Log handle stream

	readWriter := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream)) // Init reader/writer for stream

	defer readWriter.Flush() // Flush

	targetHashBytes, err := readAsync(readWriter.Reader) // Read async

	if err != nil { // Check for errors
		return // Return
	}

	logger.Infof("handling request for transaction with hash: %s", hex.EncodeToString(targetHashBytes)) // Log handle request

	transaction, _ := (*client.Validator).GetWorkingDag().GetTransactionByHash(common.NewHash(targetHashBytes)) // Get transaction with hash

	logger.Infof("responding with serialized transaction bytes: %s (len: %d), hash: %s", hex.EncodeToString(transaction.Bytes())[:36], len(transaction.Bytes()), hex.EncodeToString(transaction.Hash.Bytes())) // Log respond

	readWriter.Write(append(transaction.Bytes(), byte('\f'))) // Write transaction bytes
}

// HandleReceiveConfigRequest handles a new stream requesting the working dag config.
func (client *Client) HandleReceiveConfigRequest(stream inet.Stream) {
	logger.Infof("handling new config request stream") // Log handle stream

	writer := bufio.NewWriter(stream) // Initialize writer

	defer writer.Flush() // Flush

	logger.Infof("responding with serialized config bytes: %s", hex.EncodeToString((*client.Validator).GetWorkingConfig().Bytes())[:36]) // Log response

	writer.Write(append((*client.Validator).GetWorkingConfig().Bytes(), byte('\f'))) // Write config bytes
}

// HandleReceiveGenesisHashRequest handles a new stream requesting for the genesis hash of the working dag.
func (client *Client) HandleReceiveGenesisHashRequest(stream inet.Stream) {
	logger.Infof("handling new genesis hash request stream") // Log handle stream

	writer := bufio.NewWriter(stream) // Initialize writer

	defer writer.Flush() // Flush

	logger.Infof("responding with genesis hash: %s", hex.EncodeToString((*client.Validator).GetWorkingDag().Genesis.Bytes())) // Log response

	writer.Write(append((*client.Validator).GetWorkingDag().Genesis.Bytes(), byte('\f'))) // Write genesis hash
}

// HandleReceiveTransactionChildHashesRequest handles a new stream requesting for the child hashes of a given transaction.
func (client *Client) HandleReceiveTransactionChildHashesRequest(stream inet.Stream) {
	logger.Infof("handling new child hash request stream") // Log handle stream

	readWriter := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream)) // Initialize reader/writer from stream

	defer readWriter.Flush() // Flush

	parentHashBytes, err := readAsync(readWriter.Reader) // Read async

	if err != nil { // Check for errors
		return // Return
	}

	children, err := (*client.Validator).GetWorkingDag().GetTransactionChildren(common.NewHash(parentHashBytes)) // Get children

	if err == nil { // Check no error
		var summarizedChildHashes []byte // Init summarized child hashes buffer

		for x, child := range children { // Iterate through children
			if x == len(children)-1 { // Check is last
				summarizedChildHashes = append(summarizedChildHashes, child.Hash[:]...) // Append hash
			}

			summarizedChildHashes = append(summarizedChildHashes, append(child.Hash[:], []byte("end_hash")...)...) // Append hash
		}

		if hexEncodedChildHashes := hex.EncodeToString(summarizedChildHashes); hexEncodedChildHashes != "" { // Check can log
			logger.Infof("responding with child hashes: %s", hex.EncodeToString(summarizedChildHashes)[:36]) // Log response
		}

		readWriter.Write(append(summarizedChildHashes, byte('\f'))) // Write child hashes

		readWriter.Flush() // Flush
	}
}

/*
	END HANDLERS
*/

/* END EXPORTED METHODS */
