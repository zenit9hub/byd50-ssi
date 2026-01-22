package errors

import "fmt"

// Code represents a stable error category for the core API.
type Code string

const (
	CodeInvalidInput Code = "invalid_input"
	CodeNotFound     Code = "not_found"
	CodeEmptyKey     Code = "empty_key"
	CodeInvalidKey   Code = "invalid_key"
	CodeUpstream     Code = "upstream_error"
	CodeInternal     Code = "internal_error"
)

// Error is a typed error with a stable code and message.
type Error struct {
	code    Code
	message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.code, e.message)
}

func (e *Error) Code() Code {
	return e.code
}

func New(code Code, message string) *Error {
	return &Error{code: code, message: message}
}

func Wrap(code Code, message string, err error) *Error {
	if err == nil {
		return New(code, message)
	}
	return &Error{code: code, message: fmt.Sprintf("%s: %v", message, err)}
}
