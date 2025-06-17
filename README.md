# go-assert

A simple assertion utility for Go that provides panic-based assertions with error recovery.

## Installation

```go
go get github.com/iam-kevin/go-assert
```

## Usage

```go
import "github.com/iam-kevin/go-assert"

// Basic assertion
assert.Is(condition, "Custom error message")

// Nil checking
assert.IsNil(obj, "Should be nil")

// Error checking
assert.ErrorIsNil(err)

// Recovery pattern
func safeFunction() (err error) {
    defer assert.Capture(func(capturedErr error) {
        err = capturedErr
    })
    
    assert.Is(someCondition, "Something went wrong")
    return nil
}
```

## Functions

- `Is(condition bool, options ...AssertOption)` - Assert condition is true
- `IsNil(obj interface{}, options ...AssertOption)` - Assert object is nil
- `ErrorIsNil(err error)` - Assert error is nil
- `Capture(func(err error))` - Recover from assertion panics

All assertions panic on failure. Use `Capture` with `defer` to handle failures gracefully.