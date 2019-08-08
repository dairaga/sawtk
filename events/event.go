package events

import (
	"github.com/dairaga/sawtk/errors"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	"github.com/hyperledger/sawtooth-sdk-go/protobuf/events_pb2"
)

// Add adds event to chain.
func Add(ctx *processor.Context, typ string, data []byte, attrs ...string) *processor.InvalidTransactionError {
	size := len(attrs)

	if size%2 != 0 {
		return errors.BadParameters.TxErrorf("length of attributes are not event: %d", size)
	}

	attributes := make([]processor.Attribute, size/2+1)

	for i := 0; i < size; i = i + 2 {
		k := attrs[i]
		v := attrs[i+1]
		attributes[i/2] = processor.Attribute{
			Key:   k,
			Value: v,
		}
	}

	if err := ctx.AddEvent(typ, attributes, data); err != nil {
		return errors.Events.TxErrore(err)
	}

	return nil
}

// AddMessage adds event to chain.
func AddMessage(ctx *processor.Context, typ string, data proto.Message, attrs ...string) *processor.InvalidTransactionError {

	var dataBytes []byte
	var err error

	if data != nil {
		dataBytes, err = proto.Marshal(data)
		if err != nil {
			return errors.Marshal.TxErrore(err)
		}
	}
	return Add(ctx, typ, dataBytes, attrs...)
}

// Attr returns index and value in event attributes with key. Returns index -1 if key is not found.
func Attr(evt *events_pb2.Event, key string) (int, string) {

	for i, x := range evt.Attributes {
		if x.Key == key {
			return i, x.Value
		}
	}

	return -1, ""
}
