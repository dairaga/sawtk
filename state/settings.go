package state

import (
	"github.com/dairaga/sawtk/errors"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/setting_pb2"
)

// Setting returns the value in sawtooth settings family.
func Setting(ctx *processor.Context, addr string) (string, *processor.InvalidTransactionError) {
	setting := new(setting_pb2.Setting)
	if ok, err := Get(ctx, addr, setting); err != nil {
		return "", err
	} else if !ok {
		return "", errors.NotFound.TxErrorf("settings not found: %s", addr)
	}

	return setting.Entries[0].Value, nil
}
