// Package crypto represents the crypto RPC server.
package crypto

import (
	"context"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"

	"github.com/polaris-project/go-polaris/crypto"
	cryptoProto "github.com/polaris-project/go-polaris/internal/proto/crypto"
)

// Server represents a Polaris RPC server.
type Server struct{}

/* BEGIN EXPORTED METHODS */

// AddressFromPrivateKey handles the AddressFromPrivateKey request method.
func (server *Server) AddressFromPrivateKey(ctx context.Context, request *cryptoProto.GeneralRequest) (*cryptoProto.GeneralResponse, error) {
	decodedBytes, err := hex.DecodeString(request.PrivatePublicKey) // Decode hex-encoded private key string

	if err != nil { // Check for errors
		return &cryptoProto.GeneralResponse{}, err // Return found error
	}

	block, _ := pem.Decode(decodedBytes) // Decode private key pem

	privateKey, err := x509.ParseECPrivateKey(block.Bytes) // Parse PEM block

	if err != nil { // Check for errors
		return &cryptoProto.GeneralResponse{}, err // Return found error
	}

	return &cryptoProto.GeneralResponse{Message: hex.EncodeToString(crypto.AddressFromPrivateKey(privateKey).Bytes())}, nil // Return address value
}

/* END EXPORTED METHODS */
