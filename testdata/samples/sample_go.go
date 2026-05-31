// Package samples contains deliberate cryptographic usage for testing.
package samples

import (
	"crypto/aes"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
)

// GenerateRSAKey creates a 2048-bit RSA private key.
// This uses RSA key generation which is quantum-vulnerable.
func GenerateRSAKey() (*rsa.PrivateKey, error) {
	// RSA is vulnerable to Shor's algorithm on quantum computers.
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("generating RSA key: %w", err)
	}
	return key, nil
}

// HashWithSHA1 computes a SHA-1 hash of the given data.
// SHA-1 is cryptographically broken and should not be used.
func HashWithSHA1(data []byte) []byte {
	h := sha1.New()
	h.Write(data)
	return h.Sum(nil)
}

// GenerateECDSAKey creates an ECDSA key using P-256.
// ECDSA is vulnerable to quantum attacks via Shor's algorithm.
func GenerateECDSAKey() (*ecdsa.PrivateKey, error) {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("generating ECDSA key: %w", err)
	}
	return key, nil
}

// EncryptAES creates an AES cipher block from the given key.
// AES-128 has reduced security under Grover's algorithm (64-bit effective).
// AES-256 remains quantum-safe (128-bit effective security).
func EncryptAES(key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("creating AES cipher: %w", err)
	}

	plaintext := []byte("quantum computing test data block")
	ciphertext := make([]byte, block.BlockSize())
	block.Encrypt(ciphertext, plaintext[:block.BlockSize()])
	return ciphertext, nil
}

// HashWithSHA256 computes a SHA-256 hash. SHA-256 is considered quantum-safe
// with 128-bit post-quantum security.
func HashWithSHA256(r io.Reader) ([]byte, error) {
	h := sha256.New()
	if _, err := io.Copy(h, r); err != nil {
		return nil, fmt.Errorf("hashing with SHA-256: %w", err)
	}
	return h.Sum(nil), nil
}
