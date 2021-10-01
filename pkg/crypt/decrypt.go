package crypt

import (
	"crypto/aes"
	"crypto/cipher"
)

func DecryptChip(credentials []byte, secret []byte) ([]byte, error) {
	block, err := aes.NewCipher(secret)
	if err != nil {
		return nil, err
	}

	aed, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aed.NonceSize()
	if len(credentials) < nonceSize {
		return nil, nil
	}

	nonce, cipherText := credentials[:nonceSize], credentials[nonceSize:]
	credentials, err = aed.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return nil, err
	}
	return credentials, nil
}

