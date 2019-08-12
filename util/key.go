package util

// IsPublicKey returns key is a valid public key generated from ECDSA Secp256k1.
func IsPublicKey(key string) bool {
	//keybytes, err := hex.DecodeString(key)
	//return err == nil && len(keybytes) == 33

	if len(key) != 66 {
		return false
	}

	return IsHexString(key)
}

// IsHexString return s is a hex string or not.
func IsHexString(s string) bool {
	if len(s)%2 != 0 {
		return false
	}

	for _, x := range s {
		if !(('a' <= x && x <= 'f') || ('0' <= x && x <= '9')) {
			return false
		}
	}

	return true
}

// IsSignature returns s is a valid signature or not.
func IsSignature(s []byte) bool {
	return len(s) == 64
}
