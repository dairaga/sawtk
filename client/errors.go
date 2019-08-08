package client

import (
	"encoding/json"
	"net"
)

// SawtoothError is the error defined in sawtooth resful api
type SawtoothError struct {
	Code    int32  `json:"code,omitempty"`
	Title   string `json:"title,omitempty"`
	Message string `json:"message,omitempty"`
}

func (se *SawtoothError) Error() string {
	return marshalToString(se)
}

// Error ...
type Error struct {
	HTTPCode int            `json:"http_code,omitempty"`
	Data     *SawtoothError `json:"data,omitempty"`
}

func (e *Error) Error() string {
	return marshalToString(e)
}

// NewError returns an error with http status code and sawtooth error.
func NewError(httpCode int, raw []byte) error {
	data := new(SawtoothError)
	if err := json.Unmarshal(raw, data); err != nil {
		return err
	}

	return &Error{
		HTTPCode: httpCode,
		Data:     data,
	}
}

// IsTimeError return err is a time out error or not
func IsTimeError(err error) bool {
	v, ok := err.(net.Error)
	return ok && v != nil && v.Timeout()
}

// IsSawtkError returns err is sawtk error or not.
func IsSawtkError(err error) bool {
	v, ok := err.(*Error)
	return ok && v != nil
}
