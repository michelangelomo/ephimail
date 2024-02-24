// internal/encryption/encryption.go
package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
)

// EncryptEmail encrypts an email body using the recipient's public key
func EncryptEmail(emailBody, publicKeyStr string) (string, error) {
	// Decode the public key from base64
	publicKeyBytes, err := base64.StdEncoding.DecodeString(publicKeyStr)
	if err != nil {
		return "", fmt.Errorf("failed to decode public key: %w", err)
	}

	// Parse the public key
	block, _ := pem.Decode(publicKeyBytes)
	if block == nil {
		// Try parsing as DER if PEM fails
		block = &pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: publicKeyBytes,
		}
	}

	publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", fmt.Errorf("failed to parse public key: %w", err)
	}

	rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("public key is not an RSA key")
	}

	// Encrypt the email body
	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		rsaPublicKey,
		[]byte(emailBody),
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt email: %w", err)
	}

	// Encode the encrypted data to base64
	encryptedStr := base64.StdEncoding.EncodeToString(encryptedBytes)
	return encryptedStr, nil
}

// IsEncrypted checks if the text is likely to be encrypted (base64 encoded)
func IsEncrypted(text string) bool {
	_, err := base64.StdEncoding.DecodeString(text)
	return err == nil
}

// EncryptWithSymmetricKey implements hybrid encryption for larger emails
func EncryptWithSymmetricKey(emailBody, publicKeyStr string) (string, error) {
	// For larger emails, we would implement hybrid encryption here:
	// 1. Generate a random symmetric key
	// 2. Encrypt the email with the symmetric key
	// 3. Encrypt the symmetric key with the recipient's public key
	// 4. Combine the encrypted symmetric key and encrypted email

	// This is a placeholder for the implementation
	return EncryptEmail(emailBody, publicKeyStr)
}
