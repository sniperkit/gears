package gears

import "fmt"

// StatusError is an interface used in case
// an http status code needs to be passed
// with the error message
type StatusError interface {
	Error() string
	Status() int
}

// Error is an interface used in cases
// when an error-code string and a longer
// description needs to be passed with the
// original error message
type Error interface {
	StatusError
	Code() string
	Description() string
}

// DetailedError is an interface for being used
// for cases where more context is required to be
// provided with the error message.
type DetailedError interface {
	Error
	Details() interface{}
}

// JSONError is an interface used for handling
// errors as JSON encoded bytes.
type JSONError interface {
	Error() string
	MarshalJSON() ([]byte, error)
}

type detailedError struct {
	status           int
	ErrorCode        string      `json:"code"`
	ErrorDescription string      `json:"description"`
	ErrorDetails     interface{} `json:"details,omitempty"`
}

// NewError returns an implementation of Error interface.
// Use NewError to set the value for the "error" key in the
// context of a middleware. The details can be a nil interface.
func NewError(status int, code, description string, details interface{}) Error {
	return detailedError{
		status,
		code,
		description,
		details,
	}
}

func (e detailedError) Error() string {
	return fmt.Sprintf("%s: %s.", e.ErrorCode, e.ErrorDescription)
}

func (e detailedError) Code() string {
	return e.ErrorCode
}

func (e detailedError) Description() string {
	return e.ErrorDescription
}

func (e detailedError) Details() interface{} {
	return e.ErrorDetails
}

func (e detailedError) Status() int {
	return e.status
}
