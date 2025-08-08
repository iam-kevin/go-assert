// Package assert provides utility functions for assertions in Go programs.
// It offers a panic-based assertion mechanism that can be used for testing
// and runtime validation with optional error recovery.
package assert

import (
	"errors"
	"fmt"
)

// assertOptions holds configuration options for assertion functions.
type assertOptions struct {
	reason AssersionError
}

// AssertOption is an interface for assertion configuration options.
// It can be a string (converted to an error message) or a function
// that modifies assertOptions.
type AssertOption interface{}

// AssersionError wraps an error to provide assertion-specific error handling.
// It implements the error interface and supports error unwrapping.
type AssersionError struct {
	err error
}

// Error returns a formatted string representation of the assertion error.
func (ae *AssersionError) Error() string {
	return fmt.Sprintf("AssersionError: %s", ae.err.Error())
}

// Unwrap returns the underlying error, supporting Go's error unwrapping.
func (ae *AssersionError) Unwrap() error {
	return ae.err
}

// ErrorIsNil panics with an AssersionError if the provided error is not nil.
// This is useful for asserting that operations complete without errors.
func ErrorIsNil(err error) {
	if err != nil {
		panic(AssersionError{err})
	}
}

// Capture recovers from assertion panics and converts them back to errors.
// It should be used with defer to handle assertion failures gracefully.
// The deferrer function is called with the underlying error if an AssersionError
// was caught, otherwise it's not called.
func Capture(deferrer func(err error)) {
	if e := recover(); e != nil {
		if aerr, ok := e.(AssersionError); ok {
			deferrer(aerr.Unwrap())
		} else {
			panic(e)
		}
	}
}

// normError normalizes various error types into an AssersionError.
// It handles string messages, existing errors, and unknown types.
func normError(e interface{}) AssersionError {
	switch v := e.(type) {
	case string:
		{
			return AssersionError{errors.New(v)}
		}
	case error:
		{
			return AssersionError{v}
		}
	default:
		{
			return AssersionError{errors.New("unknown error object used")}
		}
	}
}

// mergeOptions processes and combines assertion options into a single configuration.
// It handles string options (converted to error messages) and function options
// that modify the assertOptions struct.
func mergeOptions(opts []AssertOption) *assertOptions {
	o := new(assertOptions)
	for _, w := range opts {
		switch l := w.(type) {
		case func(opts *assertOptions):
			{
				l(o)
			}
		case string:
			{
				o.reason = normError(l)
			}
		}

	}
	return o
}

// Is panics when the condition resolves to false.
// Optional AssertOption parameters can be provided to customize the error message.
func Is(condition bool, messageWithArgs ...AssertOption) {
	if !condition {
		opts := mergeOptions(messageWithArgs)
		panic(opts.reason)
	}
}

// IsNil panics when the input object is not nil.
// Optional AssertOption parameters can be provided to customize the error message.
func IsNil(obj interface{}, messageWithArgs ...AssertOption) {
	Is(obj == nil, messageWithArgs...)
}
