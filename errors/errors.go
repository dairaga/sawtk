package errors

import (
	"encoding/binary"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
)

// TxErrorf returns a InvaildTransactionError with a formatted string.
func TxErrorf(format string, a ...interface{}) *processor.InvalidTransactionError {
	return &processor.InvalidTransactionError{
		Msg: fmt.Sprintf(format, a...),
	}
}

// TxErrore returns a InvaildTransactionError with an error.
func TxErrore(err error) *processor.InvalidTransactionError {
	return &processor.InvalidTransactionError{
		Msg: err.Error(),
	}
}

// TxErrorp returns a InvalidTransactionError with a protobuf message.
func TxErrorp(msg string, m proto.Message) *processor.InvalidTransactionError {
	databytes, err := proto.Marshal(m)
	if err != nil {
		return TxErrore(err)
	}

	return &processor.InvalidTransactionError{
		Msg:          msg,
		ExtendedData: databytes,
	}
}

// ErrCode ...
type ErrCode uint32

// ToErrCode convert bytes to CtzErrCode
func ToErrCode(raw []byte) ErrCode {
	return ErrCode(binary.LittleEndian.Uint32(raw))
}

// ToBytes converts ErrCode to bytes.
func (c ErrCode) ToBytes() []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, uint32(c))
	return bs
}

// TxErrorf returns a InvalidTransactionError.
func (c ErrCode) TxErrorf(format string, a ...interface{}) *processor.InvalidTransactionError {
	return &processor.InvalidTransactionError{
		Msg:          fmt.Sprintf(format, a...),
		ExtendedData: c.ToBytes(),
	}
}

// TxErrore returns a InvalidTransactionError.
func (c ErrCode) TxErrore(err error) *processor.InvalidTransactionError {
	return &processor.InvalidTransactionError{
		Msg:          err.Error(),
		ExtendedData: c.ToBytes(),
	}
}

func (c ErrCode) String() string {
	return fmt.Sprintf("{err_code: %d}", uint32(c))
}

// Errors of SawTK.
const (
	Internal      ErrCode = 999987 // internal error.
	Wallet        ErrCode = 999988 // generating wallet failure.
	BadParameters ErrCode = 999989 // any fields in request is invalid.
	LenNotMatch   ErrCode = 999990 // data length not match.
	Conflict      ErrCode = 999991 // address conflict.
	NotFound      ErrCode = 999992 // state not found.
	GetState      ErrCode = 999993 // getting state error.
	SetState      ErrCode = 999994 // setting state error.
	Events        ErrCode = 999995 // add event error.
	ReceiptData   ErrCode = 999996 // add event error.
	Unmarshal     ErrCode = 999997 // data (protobuf) unmarshal failure.
	Marshal       ErrCode = 999998 // data (protobuf) marshal failure.
	UnknownCmd    ErrCode = 999999 // unknown command.
)
