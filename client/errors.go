package client

import (
	"encoding/json"
	"fmt"
	"net"
)

// SawtoothError is the error defined in sawtooth resful api.
type SawtoothError struct {
	Code    int32  `json:"code,omitempty"`
	Title   string `json:"title,omitempty"`
	Message string `json:"message,omitempty"`
}

func (se *SawtoothError) Error() string {
	return fmt.Sprintf(`{"code": %d, "title": %q, "message": %q}`, se.Code, se.Title, se.Message)
}

// Error wraps Sawtooth error and http response code.
type Error struct {
	code int            // http response code
	err  *SawtoothError // sawtooth error
}

// HTTPStatusCode returns http status code
func (e *Error) HTTPStatusCode() int {
	return e.code
}

// Unwrap implements error unwrap interface from go1.13
func (e *Error) Unwrap() error {
	return e.err
}

func (e *Error) Error() string {
	return fmt.Sprintf(`{"http_code": %d, "error": %s}`, e.code, e.err.Error())
}

// NewError returns an error with http status code and sawtooth error.
func NewError(httpCode int, raw []byte) error {
	se := new(SawtoothError)
	if err := json.Unmarshal(raw, se); err != nil {
		return err
	}

	return &Error{
		code: httpCode,
		err:  se,
	}
}

// IsTimeError return err is a time out error or not
func IsTimeError(err error) bool {
	//x := &net.OpError{}
	//return errors.As(err, &x) && x.Timeout()
	v, ok := err.(net.Error)
	return ok && v != nil && v.Timeout()
}

// IsSawtkError returns err is sawtk error or not.
func IsSawtkError(err error) bool {
	v, ok := err.(*Error)
	return ok && v != nil
}
