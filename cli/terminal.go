// Package cli implements the terminal.
package cli

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/polaris-project/go-polaris/common"

	accountsProto "github.com/polaris-project/go-polaris/internal/proto/accounts"
	configProto "github.com/polaris-project/go-polaris/internal/proto/config"
	cryptoProto "github.com/polaris-project/go-polaris/internal/proto/crypto"
)

var (
	// ErrInvalidParams is an error definition describing invalid input parameters.
	ErrInvalidParams = errors.New("invalid parameters")
)

/* BEGIN EXPORTED METHODS */

// NewTerminal attempts to start a handler for term commands.
func NewTerminal(rpcPort uint, rpcAddress string) {
	reader := bufio.NewScanner(os.Stdin) // Init reader

	transport := &http.Transport{ // Init transport
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	for {
		fmt.Print("\n> ") // Print prompt

		reader.Scan() // Scan

		input := reader.Text() // Fetch string input

		input = strings.TrimSuffix(input, "\n") // Trim newline

		receiver, methodname, params, err := common.ParseStringMethodCall(input) // Attempt to parse as method call

		if err != nil { // Check for errors
			fmt.Println(err.Error()) // Log found error

			continue // Continue
		}

		handleCommand(receiver, methodname, params, rpcPort, rpcAddress, transport) // Handle command
	}
}

// handleCommand runs the handler for a given receiver.
func handleCommand(receiver string, methodname string, params []string, rpcPort uint, rpcAddress string, transport *http.Transport) {
	cryptoClient := cryptoProto.NewCryptoProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport})       // Init crypto client
	accountsClient := accountsProto.NewAccountsProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport}) // Init accounts client
	configClient := configProto.NewConfigProtobufClient("https://"+rpcAddress+":"+strconv.Itoa(int(rpcPort)), &http.Client{Transport: transport})       // Init config client

	switch receiver {
	case "crypto":
		err := handleCrypto(&cryptoClient, methodname, params) // Handle crypto

		if err != nil { // Check for errors
			fmt.Println("\n" + err.Error()) // Log found error
		}
	case "accounts":
		err := handleAccounts(&accountsClient, methodname, params) // Handle accounts

		if err != nil { // Check for errors
			fmt.Println("\n" + err.Error()) // Log found error
		}
	case "config":
		err := handleConfig(&configClient, methodname, params) // Handle config

		if err != nil { // Check for errors
			fmt.Println("\n" + err.Error()) // Log found error
		}
	default:
		fmt.Println("\n" + "unrecognized namespace " + `"` + receiver + `"` + ", available namespaces: crypto, accounts") // Log invalid namespace
	}
}

// handleCrypto handles the crypto receiver.
func handleCrypto(cryptoClient *cryptoProto.Crypto, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname { // Handle different methods
	case "Sha3", "Sha3d":
		if len(params) != 1 { // Check for invalid params
			return ErrInvalidParams // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&cryptoProto.GeneralRequest{B: []byte(params[0])})) // Append params
	case "Sha3n":
		if len(params) != 2 { // Check for invalid params
			return ErrInvalidParams // return error
		}

		intVal, _ := strconv.Atoi(params[1]) // Convert to int

		reflectParams = append(reflectParams, reflect.ValueOf(&cryptoProto.GeneralRequest{B: []byte(params[0]), N: float64(intVal)})) // Append params
	case "AddressFromPrivateKey", "AddressFromPublicKey":
		if len(params) != 1 { // Check for invalid params
			return ErrInvalidParams // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&cryptoProto.GeneralRequest{PrivatePublicKey: params[0]})) // Append params
	default:
		return errors.New("illegal method: " + methodname + ", available methods: Sha3(), Sha3d(), Sha3n(), AddressFromPrivateKey(), AddressFromPublicKey()") // Return error
	}

	result := reflect.ValueOf(*cryptoClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*cryptoProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		return result[1].Interface().(error) // Return error
	}

	fmt.Println(response.Message) // Log response

	return nil // No error occurred, return nil
}

// handleAccounts handles the accounts receiver.
func handleAccounts(accountsClient *accountsProto.Accounts, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname { // Handle different methods
	case "NewAccount":
		if len(params) != 0 { // Check for invalid params
			return ErrInvalidParams // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&cryptoProto.GeneralRequest{})) // Append params
	case "AccountFromKey", "Address", "PublicKey", "PrivateKey", "String":
		if len(params) != 1 { // Check for invalid params
			return ErrInvalidParams // return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&cryptoProto.GeneralRequest{PrivatePublicKey: params[0]})) // Append params
	default:
		return errors.New("illegal method: " + methodname + ", available methods: Sha3(), Sha3d(), Sha3n(), AddressFromPrivateKey(), AddressFromPublicKey()") // Return error
	}

	result := reflect.ValueOf(*accountsClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*accountsProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		return result[1].Interface().(error) // Return error
	}

	fmt.Println(response.Message) // Log response

	return nil // No error occurred, return nil
}

// handleConfig handles the config receiver.
func handleConfig(configClient *configProto.Config, methodname string, params []string) error {
	reflectParams := []reflect.Value{} // Init buffer

	reflectParams = append(reflectParams, reflect.ValueOf(context.Background())) // Append request context

	switch methodname { // Handle different methods
	case "GetAllConfigs":
		if len(params) != 0 { // Check for invalid params
			return ErrInvalidParams // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&configProto.GeneralRequest{})) // Append params
	case "NewDagConfig":
		if len(params) != 1 { // Check for invalid params
			return ErrInvalidParams // Return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&configProto.GeneralRequest{FilePath: params[0]})) // Append params
	case "GetConfig":
		if len(params) != 1 { // Check for invalid params
			return ErrInvalidParams // return error
		}

		reflectParams = append(reflectParams, reflect.ValueOf(&configProto.GeneralRequest{Network: params[0]})) // Append params
	default:
		return errors.New("illegal method: " + methodname + ", available methods: Sha3(), Sha3d(), Sha3n(), AddressFromPrivateKey(), AddressFromPublicKey()") // Return error
	}

	result := reflect.ValueOf(*configClient).MethodByName(methodname).Call(reflectParams) // Call method

	response := result[0].Interface().(*configProto.GeneralResponse) // Get response

	if result[1].Interface() != nil { // Check for errors
		return result[1].Interface().(error) // Return error
	}

	fmt.Println(response.Message) // Log response

	return nil // No error occurred, return nil
}

/* END EXPORTED METHODS */
