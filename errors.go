package jsonrpc

import (
	"errors"
	"strings"
)

var (
	// ErrInvalidJSONRPCVersion is returned when the jsonrpc version is not "2.0".
	ErrInvalidJSONRPCVersion = errors.New("invalid jsonrpc version")
	// ErrEmptyId is returned when the id value is empty.
	ErrEmptyId = errors.New("empty id value")
	// ErrEmptyMethod is returned when the method value is empty.
	ErrEmptyMethod = errors.New("empty method value")
	// ErrResultAndError is returned when both result and error are set in a response.
	ErrResultAndError = errors.New("result and error cannot be set at the same time")
	// ErrResultOrError is returned when neither result nor error are set in a response.
	ErrResultOrError = errors.New("result or error should be set")
	// ErrInvalidParams is returned when the params value is not an array or object.
	ErrInvalidParams = errors.New("invalid params: must be an array or object")
)

// ValidationErrors is a collection of errors encountered during validation.
type ValidationErrors struct {
	Errors []error
}

// Append adds an error to the collection.
func (v *ValidationErrors) Append(err error) {
	v.Errors = append(v.Errors, err)
}

// Error returns a string representation of all errors in the collection.
func (v ValidationErrors) Error() string {
	msgs := make([]string, len(v.Errors))

	for i, err := range v.Errors {
		msgs[i] = err.Error()
	}

	return strings.Join(msgs, ", ")
}

// HasErrors returns true if the collection contains any errors.
func (v ValidationErrors) HasErrors() bool {
	return len(v.Errors) > 0
}

// Unwrap returns the list of errors in the collection.
func (v ValidationErrors) Unwrap() []error {
	return v.Errors
}
