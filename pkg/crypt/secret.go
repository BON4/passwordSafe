package crypt

import "crypto/sha256"

func ValidateSecretKey(key []byte) []byte {
	shakey := sha256.Sum256(key)
	return shakey[:16]
}
