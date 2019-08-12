package tp

import (
	"github.com/dairaga/sawtk/wallet"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
)

// MakeWallet return a wallet address from public key
func MakeWallet(pubkey string) (string, *processor.InvalidTransactionError) {
	w, err := wallet.MakeFromHex(pubkey)
	if err != nil {
		return "", Wallet.TxErrorf("make wallet %s: %v", pubkey, err)
	}
	return w, nil
}
