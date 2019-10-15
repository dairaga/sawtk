package util

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

// SHA512 returns a SHA-512 hex string.
func SHA512(data []byte) string {
	hash := sha512.New()
	hash.Write(data)
	bytes := hash.Sum(nil)
	return hex.EncodeToString(bytes)
}

// SHA256 returns a SHA-256 hex string.
func SHA256(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	bytes := hash.Sum(nil)
	return hex.EncodeToString(bytes)
}
