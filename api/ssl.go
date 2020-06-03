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
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/polaris-project/go-polaris/common"
)

/* BEGIN INTERNAL METHODS */

// generateCert generates an ssl cert.
func generateCert(certName string, hosts []string) error {
	os.Remove(fmt.Sprintf("%sKey.pem", certName))  // Remove existent key
	os.Remove(fmt.Sprintf("%sCert.pem", certName)) // Remove existent cert

	privateKey, err := generateTLSKey(certName) // Generate key

	if err != nil { // Check for errors
		return err // Return found error
	}

	err = generateTLSCert(privateKey, certName, hosts) // Generate cert

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// generateTLSKey generates necessary TLS keys.
func generateTLSKey(keyName string) (*ecdsa.PrivateKey, error) {
	err := common.CreateDirIfDoesNotExist(common.CertificatesDir) // Create certs dir if does not exist

	if err != nil { // Check for errors
		return &ecdsa.PrivateKey{}, err // Return found error
	}

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader) // Generate private key

	if err != nil { // Check for errors
		return &ecdsa.PrivateKey{}, err // Return found error
	}

	marshaledPrivateKey, err := x509.MarshalECPrivateKey(privateKey) // Marshal private key

	if err != nil { // Check for errors
		return &ecdsa.PrivateKey{}, err // Return found error
	}

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: marshaledPrivateKey}) // Encode to memory

	err = ioutil.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/%sKey.pem", common.CertificatesDir, keyName)), pemEncoded, 0644) // Write pem

	if err != nil { // Check for errors
		return &ecdsa.PrivateKey{}, err // Return found error
	}

	return privateKey, nil // No error occurred, return nil
}

// generateTLSCert generates necessary TLS certs.
func generateTLSCert(privateKey *ecdsa.PrivateKey, certName string, hosts []string) error {
	err := common.CreateDirIfDoesNotExist(common.CertificatesDir) // Create certs dir if does not exist

	if err != nil { // Check for errors
		return err // Return found error
	}

	notBefore := time.Now() // Fetch current time

	notAfter := notBefore.Add(2 * 24 * time.Hour) // Fetch 'deadline'

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)     // Init limit
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit) // Init serial number

	if err != nil { // Check for errors
		return err // Return found error
	}

	template := x509.Certificate{ // Init template
		SerialNumber: serialNumber, // Generate w/serial number
		Subject: pkix.Name{ // Generate w/subject
			Organization: hosts, // Generate w/org
		},
		NotBefore: notBefore, // Generate w/not before
		NotAfter:  notAfter,  // Generate w/not after

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature, // Generate w/key usage
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},               // Generate w/ext key
		BasicConstraintsValid: true,                                                         // Generate w/basic constraints
		IsCA:                  true,                                                         // Force is CA
	}

	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey) // Generate certificate

	if err != nil { // Check for errors
		return err // Return found error
	}

	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert}) // Encode pem

	err = ioutil.WriteFile(filepath.FromSlash(fmt.Sprintf("%s/%sCert.pem", common.CertificatesDir, certName)), pemEncoded, 0644) // Write cert file

	if err != nil { // Check for errors
		return err // Return found error
	}

	return nil // No error occurred, return nil
}

// keyToFile writes a given key to a pem-encoded file.
func keyToFile(filename string, key *ecdsa.PrivateKey) error {
	err := common.CreateDirIfDoesNotExist(common.CertificatesDir) // Create cert dir if doesn't already exist

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
	err := common.CreateDirIfDoesNotExist(common.CertificatesDir) // Create cert dir if doesn't already exist

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
