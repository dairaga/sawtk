package tp

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/setting_pb2"
)

func mapKeys(data map[string]proto.Message) []string {
	size := len(data)
	if size <= 0 {
		return []string{}
	}

	tmp := make([]string, size, size)
	i := 0
	for k := range data {
		tmp[i] = k
		i++
	}

	return tmp
}

// ----------------------------------------------------------------------------

// Attribute returns a sawtooth event attribute.
func Attribute(key, value string) processor.Attribute {
	return processor.Attribute{
		Key:   key,
		Value: value,
	}
}

// ----------------------------------------------------------------------------

// Context wraps sawtooth processor context.
type Context struct {
	ref    *processor.Context
	cmd    int32
	signer string
}

func (ctx *Context) String() string {
	return fmt.Sprintf(`{cmd: %d, signer: "%s"}`, ctx.Cmd(), ctx.signer)
}

// Cmd returns current transaction command.
func (ctx *Context) Cmd() int32 {
	return ctx.cmd
}

// SignerPublicKey returns transaction signer public key.
func (ctx *Context) SignerPublicKey() string {
	return ctx.signer
}

// Context returns internal *process.Context.
func (ctx *Context) Context() *processor.Context {
	return ctx.ref
}

// ----------------------------------------------------------------------------

// GetAll returns states with multiple addresses.
func (ctx *Context) GetAll(data map[string]proto.Message) *processor.InvalidTransactionError {
	keys := mapKeys(data)

	result, err := ctx.ref.GetState(keys)
	if err != nil {
		return GetState.TxErrore(err)
	}

	for k, v := range data {
		b, ok := result[k]
		if !ok {
			// data not found and remove input map with key.
			delete(data, k)
			continue
		}
		if v != nil {
			// value in input data and umarshal result bytes from chain.
			if err := proto.Unmarshal(b, v); err != nil {
				return Unmarshal.TxErrore(err)
			}
		}
	}

	return nil
}

// Get returns a state with address.
func (ctx *Context) Get(address string, data proto.Message) (bool, *processor.InvalidTransactionError) {
	m := map[string]proto.Message{address: data}

	err := ctx.GetAll(m)
	if err != nil {
		return false, err
	}

	_, ok := m[address]
	return ok, nil
}

// ----------------------------------------------------------------------------

// SetAll sets all data into chain.
func (ctx *Context) SetAll(data map[string]proto.Message) *processor.InvalidTransactionError {
	tmp := make(map[string][]byte)

	for k, v := range data {
		dataBytes, err := proto.Marshal(v)
		if err != nil {
			return Marshal.TxErrore(err)
		}
		tmp[k] = dataBytes
	}

	resp, err := ctx.ref.SetState(tmp)
	if err != nil {
		return SetState.TxErrore(err)
	}

	if len(resp) != len(data) {
		return LenNotMatch.TxErrorf("length of responses (%d) are not same with input (%d)", len(resp), len(data))
	}

	return nil
}

// Set state into chain.
func (ctx *Context) Set(address string, data proto.Message) *processor.InvalidTransactionError {
	tmp := map[string]proto.Message{address: data}
	return ctx.SetAll(tmp)
}

// ----------------------------------------------------------------------------

// Del remove state from chain.
func (ctx *Context) Del(addrs []string) ([]string, error) {
	return ctx.ref.DeleteState(addrs)
}

// ----------------------------------------------------------------------------

// AddEvent adds event to chain.
func (ctx *Context) AddEvent(typ string, data []byte, attributes ...processor.Attribute) *processor.InvalidTransactionError {

	if err := ctx.ref.AddEvent(typ, attributes, data); err != nil {
		return Events.TxErrore(err)
	}

	return nil
}

// AddEventMessage adds event to chain.
func (ctx *Context) AddEventMessage(typ string, data proto.Message, attributes ...processor.Attribute) *processor.InvalidTransactionError {
	var dataBytes []byte
	var err error

	if data != nil {
		dataBytes, err = proto.Marshal(data)
		if err != nil {
			return Marshal.TxErrore(err)
		}
	}
	return ctx.AddEvent(typ, dataBytes, attributes...)
}

// ----------------------------------------------------------------------------

// AddReceiptData adds receipt data.
func (ctx *Context) AddReceiptData(data proto.Message) *processor.InvalidTransactionError {

	databyes, err := proto.Marshal(data)
	if err != nil {
		return Marshal.TxErrore(err)
	}
	if err := ctx.ref.AddReceiptData(databyes); err != nil {
		return ReceiptData.TxErrore(err)
	}
	return nil
}

// ----------------------------------------------------------------------------

// Setting returns value in setting-tp.
func (ctx *Context) Setting(address string) (string, *processor.InvalidTransactionError) {
	setting := new(setting_pb2.Setting)
	if ok, err := ctx.Get(address, setting); err != nil {
		return "", err
	} else if !ok {
		return "", NotFound.TxErrorf("settings not found: %s", address)
	}

	return setting.Entries[0].Value, nil
}
