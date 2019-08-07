// Package state is utility to read/write sawtooth state.
package state

import (
	"github.com/dairaga/sawtk/errors"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
)

// GetAll returns states with multiple addresses.
func GetAll(ctx *processor.Context, data map[string]proto.Message) *processor.InvalidTransactionError {
	size := len(data)
	if size <= 0 {
		return errors.BadParameters.TxErrorf("data length must larger than zero: %d", size)
	}

	tmp := make([]string, 0, size)
	for k := range data {
		tmp = append(tmp, k)
	}

	result, err := ctx.GetState(tmp)
	if err != nil {
		return errors.GetState.TxErrore(err)
	}

	for k, v := range data {
		b, ok := result[k]
		if !ok {
			delete(data, k)
			continue
		}
		if v != nil {
			if err := proto.Unmarshal(b, v); err != nil {
				return errors.Unmarshal.TxErrore(err)
			}
		}
	}

	return nil
}

// Get returns a state with address.
func Get(ctx *processor.Context, address string, data proto.Message) (bool, *processor.InvalidTransactionError) {

	tmp := map[string]proto.Message{address: data}

	err := GetAll(ctx, tmp)
	if err != nil {
		return false, err
	}

	_, ok := tmp[address]
	return ok, nil
}

// SetAll sets all data into chain.
func SetAll(ctx *processor.Context, data map[string]proto.Message) *processor.InvalidTransactionError {
	if len(data) <= 0 {
		return errors.BadParameters.TxErrorf("data length must be larger than zero: %d", len(data))
	}

	tmp := make(map[string][]byte)

	for k, v := range data {
		dataBytes, err := proto.Marshal(v)
		if err != nil {
			return errors.Marshal.TxErrore(err)
		}
		tmp[k] = dataBytes
	}

	resp, err := ctx.SetState(tmp)
	if err != nil {
		return errors.SetState.TxErrore(err)
	}

	if len(resp) != len(data) {
		return errors.LenNotMatch.TxErrorf("length of responses (%d) are not same with input (%d)", len(resp), len(data))
	}

	return nil
}

// Set state into chain.
func Set(ctx *processor.Context, addr string, data proto.Message) *processor.InvalidTransactionError {
	tmp := map[string]proto.Message{addr: data}
	return SetAll(ctx, tmp)
}
