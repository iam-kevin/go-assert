package assert

import (
	"errors"
	"fmt"
)

type assertOptions struct {
	reason AssersionError
}

type AssertOption interface{}

type AssersionError struct {
	err error
}

func (ae *AssersionError) Error() string {
	return fmt.Sprintf("AssersionError: %s", ae.err.Error())
}

func (ae *AssersionError) Unwrap() error {
	return ae.err
}

func ErrorIsNil(err error) {
	if err != nil {
		panic(AssersionError{err})
	}
}

func Capture(deferrer func(err error)) {
	if e := recover(); e != nil {
		if aerr, ok := e.(AssersionError); ok {
			deferrer(aerr.Unwrap())
		}
	}
}

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

// Is panics when the condition resolves to false
func Is(condition bool, options ...AssertOption) {
	if !condition {
		opts := mergeOptions(options)
		panic(opts.reason)
	}
}

// IsNil panics when the input object is not nil
func IsNil(obj interface{}, options ...AssertOption) {
	Is(obj == nil, options...)
}
