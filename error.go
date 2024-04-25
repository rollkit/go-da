package da

import "fmt"

// Reason is the reason for the error.
type Reason uint

const (
	ReasonGasFee Reason = iota
	ReasonBlobSize
	ReasonUnknown
)

func (r Reason) String() string {
	switch r {
	case ReasonGasFee:
		return "ErrGasFee"
	case ReasonBlobSize:
		return "ErrBlobSize"
	default:
		return fmt.Sprintf("unknown(%d)", r)
	}
}

// Error is a wrapper for error, description  and reason.
type Error struct {
	err   error
	reason Reason
}

// Error satisfies the error interface.
func (e Error) Error() string {
	if e.err != nil {
		return fmt.Sprintf("%s: %v", e.reason, e.err)
	}
	return e.reason.String()
}

// Unwrap satisfies the Is/As interface.
func (e Error) Unwrap() error {
	return e.err
}

// Is satisfies the error Unwrap interface.
func (e Error) Is(target error) bool {
	if target == nil {
		return e == target
	}
	err, ok := target.(Error)
	if !ok {
		return false
	}
	return e.reason == err.reason
}

// NewError returns a custom Error.
func NewError(err error, reason Reason) error {
	return Error{
		err:   err,
		reason: reason,
	}
}

// NewGasFeeError returns a gas fee error.
func NewGasFeeError(err error) error {
	return NewError(err, ReasonGasFee)
}

// NewBlobSizeError returns a blob size limit error.
func NewBlobSizeError(err error) error {
	return NewError(err, ReasonBlobSize)
}

var (
	ErrGasFee = NewGasFeeError(nil)
	ErrBlobSize = NewBlobSizeError(nil)
)
