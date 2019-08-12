package wallet

import (
	"encoding/hex"

	"github.com/btcsuite/btcutil/base58"
)

// Make returns a wallet address.
func Make(raw []byte) string {
	return base58.CheckEncode(raw, raw[0])
}

// Validate returns input wallet is ok or not.
func Validate(wallet string) bool {
	_, _, err := base58.CheckDecode(wallet)
	return err == nil
}

// MakeFromHex returns a wallet address from hex string.
func MakeFromHex(hexstr string) (string, error) {
	raw, err := hex.DecodeString(hexstr)
	if err != nil {
		return "", err
	}

	return Make(raw), nil
}

// Must retuns a wallet address or panic.
func Must(hexstr string) string {
	ret, err := MakeFromHex(hexstr)
	if err != nil {
		panic(err)
	}

	return ret
}
