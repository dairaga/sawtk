package tp

import "github.com/dairaga/sawtk/errors"

// error definition.
const (
	Internal      errors.ErrCode = 999988 // internal error.
	Conflict      errors.ErrCode = 999989 // address conflict.
	LenNotMatch   errors.ErrCode = 999990 // data length not match.
	BadParameters errors.ErrCode = 999991 // any fields in request is invalid.
	NotFound      errors.ErrCode = 999992 // state not found.
	GetState      errors.ErrCode = 999993 // getting state error.
	SetState      errors.ErrCode = 999994 // setting state error.
	Events        errors.ErrCode = 999995 // add event error.
	ReceiptData   errors.ErrCode = 999996 // add event error.
	Unmarshal     errors.ErrCode = 999997 // data (protobuf) unmarshal failure.
	Marshal       errors.ErrCode = 999998 // data (protobuf) marshal failure.
	UnknownCmd    errors.ErrCode = 999999 // unknown command.
)
