package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// encryptGCM encrypts plaintext using AES-GCM (Authenticated Encryption)
// GCM provides both confidentiality and authenticity, protecting against tampering
func encryptGCM(plaintext string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	// GCM mode provides authenticated encryption
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Create nonce (number used once)
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Encrypt and authenticate
	// GCM's Seal appends the authentication tag to the ciphertext
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// decryptGCM decrypts ciphertext using AES-GCM
// Automatically verifies authenticity - will fail if data was tampered with
func decryptGCM(b64cipher string, key []byte) (string, error) {
	data, err := base64.StdEncoding.DecodeString(b64cipher)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	// Decrypt and verify authenticity
	// Open will return an error if the authentication tag doesn't match
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt or verify: %w", err)
	}

	return string(plaintext), nil
}

// CryptGCM implements authenticated encryption using AES-GCM
type CryptGCM struct {
	key []byte
}

// NewCryptGCM creates a new GCM-based encryption instance
func NewCryptGCM(key []byte) *CryptGCM {
	return &CryptGCM{key: key}
}

// Encrypt encrypts plaintext using AES-GCM
func (c *CryptGCM) Encrypt(plaintext string) (string, error) {
	return encryptGCM(plaintext, c.key)
}

// Decrypt decrypts ciphertext using AES-GCM
func (c *CryptGCM) Decrypt(ciphertext string) (string, error) {
	return decryptGCM(ciphertext, c.key)
}
