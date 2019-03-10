package p2p

import (
	"bufio"

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
	reader := bufio.NewReader(stream) // Initialize reader from stream

	var transactionBytes []byte // Initialize transaction bytes buffer

	for readBytes, err := reader.ReadByte(); err != nil; { // Read until EOF
		transactionBytes = append(transactionBytes, readBytes) // Append read bytes
	}

	transaction := types.TransactionFromBytes(transactionBytes) // Deserialize transaction

	if err := (*client.Validator).ValidateTransaction(transaction); err == nil { // Check transaction valid
		(*client.Validator).GetWorkingDag().AddTransaction(transaction) // Add transaction to working dag
	}
}

// HandleReceiveBestTransactionRequest handle a new stream requesting for the best transaction hash.
func (client *Client) HandleReceiveBestTransactionRequest(stream inet.Stream) {
	writer := bufio.NewWriter(stream) // Initialize writer from stream

	bestTransaction, _ := (*client.Validator).GetWorkingDag().GetBestTransaction() // Get best transaction

	writer.Write(bestTransaction.Bytes()) // Write best transaction
}

// HandleReceiveTransactionRequest handles a new stream requesting transaction metadata with a given hash.
func (client *Client) HandleReceiveTransactionRequest(stream inet.Stream) {
	readWriter := bufio.NewReadWriter(bufio.NewReader(stream), bufio.NewWriter(stream)) // Init reader/writer for stream

	var targetHashBytes []byte // Initialize tx hash bytes buffer

	for readBytes, err := readWriter.ReadByte(); err != nil; { // Read until EOF
		targetHashBytes = append(targetHashBytes, readBytes) // Append read bytes
	}

	transaction, _ := (*client.Validator).GetWorkingDag().GetTransactionByHash(common.NewHash(targetHashBytes)) // Get transaction with hash

	readWriter.Write(transaction.Bytes()) // Write transaction bytes
}

// HandleReceiveConfigRequest handles a new stream requesting a
func (client *Client) HandleReceiveConfigRequest(stream inet.Stream) {
	writer := bufio.NewWriter(stream) // Initialize writer

	writer.Write((*client.Validator).GetWorkingConfig().Bytes()) // Write config bytes
}

/*
	END HANDLERS
*/

/* END EXPORTED METHODS */
