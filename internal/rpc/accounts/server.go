// Package accounts represents the accounts RPC server.
package accounts

import (
	"context"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"strings"

	account "github.com/polaris-project/go-polaris/accounts"
	"github.com/polaris-project/go-polaris/common"

	accountsProto "github.com/polaris-project/go-polaris/internal/proto/accounts"
)

// Server represents a Polaris RPC server.
type Server struct{}

/* BEGIN EXPORTED METHODS */

// NewAccount handles the NewAccount request method.
func (server *Server) NewAccount(ctx context.Context, request *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	account, err := account.NewAccount() // Initialize new account

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	err = account.WriteToMemory() // Write account to persistent memory

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: hex.EncodeToString(account.Address().Bytes())}, nil // Return account address
}

// GetAllAccounts handles the GetAllAccounts request method.
func (server *Server) GetAllAccounts(ctx context.Context, request *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	accounts := account.GetAllAccounts() // Get all accounts

	var accountStrings []string // Init string buffer

	for _, account := range accounts { // Iterate through accounts
		accountStrings = append(accountStrings, hex.EncodeToString(account.Bytes())) // Append hex encoded address
	}

	return &accountsProto.GeneralResponse{Message: strings.Join(accountStrings, ", ")}, nil // Return accounts
}

// AccountFromKey handles the AccountFromKey request method.
func (server *Server) AccountFromKey(ctx context.Context, request *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	decodedBytes, err := hex.DecodeString(request.PrivatePublicKey) // Decode hex-encoded private key string

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	block, _ := pem.Decode(decodedBytes) // Decode private key pem

	privateKey, err := x509.ParseECPrivateKey(block.Bytes) // Parse PEM block

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: hex.EncodeToString(account.AccountFromKey(privateKey).Address().Bytes())}, nil // Return account address
}

// Address handles the Address request method.
func (server *Server) Address(ctx context.Context, request *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	addressBytes, err := hex.DecodeString(request.PrivatePublicKey) // Decode address

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	account, err := account.ReadAccountFromMemory(common.NewAddress(addressBytes)) // Read account

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: hex.EncodeToString(account.Address().Bytes())}, nil // Return account address
}

// PublicKey handles the PublicKey request method.
func (server *Server) PublicKey(ctx context.Context, request *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	addressBytes, err := hex.DecodeString(request.PrivatePublicKey) // Decode address

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	account, err := account.ReadAccountFromMemory(common.NewAddress(addressBytes)) // Read account

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	publicKeyBytes := elliptic.Marshal(elliptic.P521(), account.PublicKey().X, account.PublicKey().Y) // Marshal public key

	return &accountsProto.GeneralResponse{Message: hex.EncodeToString(publicKeyBytes)}, nil // Return marshaled account public key
}

// PrivateKey handles the PrivateKey request method.
func (server *Server) PrivateKey(ctx context.Context, request *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	addressBytes, err := hex.DecodeString(request.PrivatePublicKey) // Decode address

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	account, err := account.ReadAccountFromMemory(common.NewAddress(addressBytes)) // Read account

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	marshaledPrivateKey, err := x509.MarshalECPrivateKey(account.PrivateKey()) // Marshal private key

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: marshaledPrivateKey}) // Encode to memory

	return &accountsProto.GeneralResponse{Message: hex.EncodeToString(pemEncoded)}, nil // Return hex encoded private key
}

// String handles the String request method.
func (server *Server) String(ctx context.Context, request *accountsProto.GeneralRequest) (*accountsProto.GeneralResponse, error) {
	addressBytes, err := hex.DecodeString(request.PrivatePublicKey) // Decode address

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	account, err := account.ReadAccountFromMemory(common.NewAddress(addressBytes)) // Read account

	if err != nil { // Check for errors
		return &accountsProto.GeneralResponse{}, err // Return found error
	}

	return &accountsProto.GeneralResponse{Message: account.String()}, nil // Return account string
}

/* END EXPORTED METHODS */
