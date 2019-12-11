package util

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"

	"golang.org/x/crypto/ripemd160"
)

// SHA512Raw ...
func SHA512Raw(data []byte) []byte {
	hash := sha512.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// SHA512 returns a SHA-512 hex string.
func SHA512(data []byte) string {
	return hex.EncodeToString(SHA512Raw(data))
}

// SHA256Raw ...
func SHA256Raw(data []byte) []byte {
	hash := sha256.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// SHA256 returns a SHA-256 hex string.
func SHA256(data []byte) string {
	return hex.EncodeToString(SHA256Raw(data))
}

// RIPEMD160Raw ...
func RIPEMD160Raw(data []byte) []byte {
	hash := ripemd160.New()
	hash.Write(data)
	return hash.Sum(nil)
}

// RIPEMD160 ...
func RIPEMD160(data []byte) string {
	return hex.EncodeToString(RIPEMD160Raw(data))
}
