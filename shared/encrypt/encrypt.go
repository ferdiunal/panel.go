package encrypt

import (
	"encoding/hex"
	"log"
)

type Crypt interface {
	Encrypt(plaintext string) (string, error)
	Decrypt(ciphertext string) (string, error)
}

var _crypt Crypt

type crypt struct {
	key []byte
}

func NewCrypt(key string) Crypt {
	if _crypt != nil {
		return _crypt
	}

	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		log.Fatalf("Failed to decode encryption key: %v", err)
	}

	_crypt = &crypt{key: keyBytes}

	return _crypt
}

func (c *crypt) Encrypt(plaintext string) (string, error) {
	return encrypt(plaintext, c.key)
}

func (c *crypt) Decrypt(ciphertext string) (string, error) {
	return decrypt(ciphertext, c.key)
}
