// Package crypto represents the crypto RPC server.
package crypto

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
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
	if err != nil {                                                 // Check for errors
		return &cryptoProto.GeneralResponse{}, err // Return found error
	}

	block, _ := pem.Decode(decodedBytes) // Decode private key pem

	privateKey, err := x509.ParseECPrivateKey(block.Bytes) // Parse PEM block
	if err != nil {                                        // Check for errors
		return &cryptoProto.GeneralResponse{}, err // Return found error
	}

	return &cryptoProto.GeneralResponse{Message: hex.EncodeToString(crypto.AddressFromPrivateKey(privateKey).Bytes())}, nil // Return address value
}

// AddressFromPublicKey handles the AddressFromPublicKey request method.
func (server *Server) AddressFromPublicKey(ctx context.Context, request *cryptoProto.GeneralRequest) (*cryptoProto.GeneralResponse, error) {
	decodedBytes, err := hex.DecodeString(request.PrivatePublicKey) // Decode hex-encoded private key string
	if err != nil {                                                 // Check for errors
		return &cryptoProto.GeneralResponse{}, err // Return found error
	}

	x, y := elliptic.Unmarshal(elliptic.P521(), decodedBytes) // Unmarshal public key

	publicKey := &ecdsa.PublicKey{
		Curve: elliptic.P521(), // Set curve
		X:     x,               // Set x
		Y:     y,               // Set Y
	} // Init public key instance

	return &cryptoProto.GeneralResponse{Message: hex.EncodeToString(crypto.AddressFromPublicKey(publicKey).Bytes())}, nil // Return address value
}

// Sha3 handles the Sha3 request method.
func (server *Server) Sha3(ctx context.Context, request *cryptoProto.GeneralRequest) (*cryptoProto.GeneralResponse, error) {
	return &cryptoProto.GeneralResponse{Message: hex.EncodeToString(crypto.Sha3(request.B).Bytes())}, nil // Return hash value
}

// Sha3n handles the Sha3n request method.
func (server *Server) Sha3n(ctx context.Context, request *cryptoProto.GeneralRequest) (*cryptoProto.GeneralResponse, error) {
	return &cryptoProto.GeneralResponse{Message: hex.EncodeToString(crypto.Sha3n(request.B, uint(request.N)).Bytes())}, nil // Return hash value
}

// Sha3d handles the Sha3d request method.
func (server *Server) Sha3d(ctx context.Context, request *cryptoProto.GeneralRequest) (*cryptoProto.GeneralResponse, error) {
	return &cryptoProto.GeneralResponse{Message: hex.EncodeToString(crypto.Sha3d(request.B).Bytes())}, nil // Return hash value
}

/* END EXPORTED METHODS */
