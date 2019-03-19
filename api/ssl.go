// Package api contains all rpc and rest-related api helper methods, structs, etc...
package api

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/polaris-project/go-polaris/common"
)

/* BEGIN INTERNAL METHODS */

// generateCert generates an ssl cert.
func generateCert(certName string, hosts []string) error {
	notBefore := time.Now() // Get not before time

	notAfter := notBefore.Add(2 * (356 * (24 * time.Hour))) // Get not after

	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // Generate key

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = keyToFile(fmt.Sprintf("%sRootKey.pem", key), key) // Write key to file

	if err != nil { // Check for errors
		return err // Return found error
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128) // Get limit

	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit) // Generate serial number

	if err != nil { // Check for errors
		return err // Return found error
	}

	rootTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Polaris Node"},
			CommonName:   "Root CA",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
	} // Create template

	derBytes, err := x509.CreateCertificate(rand.Reader, &rootTemplate, &rootTemplate, &key.PublicKey, key) // Generate root bytes

	if err != nil { // Check for errors
		return err // Return found error
	}

	leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // Generate leaf key

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = keyToFile(fmt.Sprintf("%sLeafKey.pem", leafKey), leafKey) // Write leaf key to file

	if err != nil { // Check for errors
		return err // Return found error
	}

	serialNumber, err = rand.Int(rand.Reader, serialNumberLimit) // Seed rand

	if err != nil {
		return err // Return found error
	}

	leafTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Polaris Node"},
			CommonName:   "leaf_cert_1",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	} // Init leaf template

	for _, host := range hosts { // Iterate through hosts
		if ip := net.ParseIP(host); ip != nil { // Parse IP address
			leafTemplate.IPAddresses = append(leafTemplate.IPAddresses, ip) // Append parsed IP
		} else { // Could not parse
			leafTemplate.DNSNames = append(leafTemplate.DNSNames, host) // Append hostname
		}
	}

	derBytes, err = x509.CreateCertificate(rand.Reader, &leafTemplate, &rootTemplate, &leafKey.PublicKey, key) // Create certificate

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = certToFile(fmt.Sprintf("%sCertTemplate.pem", certName), derBytes) // Write certificate to file

	if err != nil { // Check for errors
		return err // Return found error
	}

	clientKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // Generate client key

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = keyToFile(fmt.Sprintf("%sClient.key", certName), clientKey) // Write client key to persistent memory

	if err != nil { // Check for errors
		return err // Return found error
	}

	clientTemplate := x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(4),
		Subject: pkix.Name{
			Organization: []string{"Polaris Node"},
			CommonName:   "client_cert",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
	} // Initialize client template

	derBytes, err = x509.CreateCertificate(rand.Reader, &clientTemplate, &rootTemplate, &clientKey.PublicKey, key) // Create certificate

	if err != nil {
		return err // Return found error
	}

	err = certToFile(fmt.Sprintf("%sClient.pem", certName), derBytes) // Write cert

	if err != nil {
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// keyToFile writes a given key to a pem-encoded file.
func keyToFile(filename string, key *ecdsa.PrivateKey) error {
	err := common.CreateDirIfDoesNotExit(common.CertificatesDir) // Create cert dir if doesn't already exist

	if err != nil { // Check for errors
		return err // Return found error
	}

	file, err := os.Create(filepath.FromSlash(fmt.Sprintf("%s/%s", common.CertificatesDir, filename))) // Create key file

	if err != nil { // Check for errors
		return err // Return found error
	}

	b, err := x509.MarshalECPrivateKey(key) // Marshal private key

	if err != nil { // Check for errors
		file.Close() // Close file

		return err // Return found error
	}

	if err := pem.Encode(file, &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}); err != nil { // Encode
		file.Close() // Close file

		return err // Return found error
	}

	return file.Close() // Close file
}

// certToFile writes a given certificate to a file.
func certToFile(filename string, derBytes []byte) error {
	err := common.CreateDirIfDoesNotExit(common.CertificatesDir) // Create cert dir if doesn't already exist

	if err != nil { // Check for errors
		return err // Return found error
	}

	certOut, err := os.Create(filepath.FromSlash(fmt.Sprintf("%s/%s", common.CertificatesDir, filename))) // Create file

	if err != nil {
		return err // Return found error
	}

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		certOut.Close() // Close cert file

		return err // Return found error
	}

	return certOut.Close() // Close cert file
}

/* END INTERNAL METHODS */
