package da

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbda "github.com/rollkit/go-da/types/pb/da"
)

// Code defines error codes for JSON-RPC.
//
// They are reused for GRPC
type Code int

// Codes are used by JSON-RPC client and server
const (
	CodeBlobNotFound               Code = 32001
	CodeBlobSizeOverLimit          Code = 32002
	CodeTxTimedOut                 Code = 32003
	CodeTxAlreadyInMempool         Code = 32004
	CodeTxIncorrectAccountSequence Code = 32005
	CodeTxTooLarge                 Code = 32006
	CodeContextDeadline            Code = 32007
	CodeFutureHeight               Code = 32008
)

// ErrBlobNotFound is used to indicate that the blob was not found.
type ErrBlobNotFound struct{}

func (e *ErrBlobNotFound) Error() string {
	return "blob: not found"
}

// ErrBlobSizeOverLimit is used to indicate that the blob size is over limit.
type ErrBlobSizeOverLimit struct{}

func (e *ErrBlobSizeOverLimit) Error() string {
	return "blob: over size limit"
}

// ErrTxTimedOut is the error message returned by the DA when mempool is congested.
type ErrTxTimedOut struct{}

func (e *ErrTxTimedOut) Error() string {
	return "timed out waiting for tx to be included in a block"
}

// ErrTxAlreadyInMempool is the error message returned by the DA when tx is already in mempool.
type ErrTxAlreadyInMempool struct{}

func (e *ErrTxAlreadyInMempool) Error() string {
	return "tx already in mempool"
}

// ErrTxIncorrectAccountSequence is the error message returned by the DA when tx has incorrect sequence.
type ErrTxIncorrectAccountSequence struct{}

func (e *ErrTxIncorrectAccountSequence) Error() string {
	return "incorrect account sequence"
}

// ErrTxTooLarge is the err message returned by the DA when tx size is too large.
type ErrTxTooLarge struct{}

func (e *ErrTxTooLarge) Error() string {
	return "tx too large"
}

// ErrContextDeadline is the error message returned by the DA when context deadline exceeds.
type ErrContextDeadline struct{}

func (e *ErrContextDeadline) Error() string {
	return "context deadline"
}

// ErrFutureHeight is returned when requested height is from the future
type ErrFutureHeight struct{}

func (e *ErrFutureHeight) Error() string {
	return "given height is from the future"
}

// gRPC checks for GPRCStatus method on errors to enable advanced error handling.

// getGRPCStatus constructs a gRPC status with error details based on the provided error, gRPC code, and DA error code.
func getGRPCStatus(err error, grpcCode codes.Code, daCode pbda.ErrorCode) *status.Status {
	base := status.New(grpcCode, err.Error())
	detailed, err := base.WithDetails(&pbda.ErrorDetails{Code: daCode})
	if err != nil {
		return base
	}
	return detailed
}

// GRPCStatus returns the gRPC status with details for an ErrBlobNotFound error.
func (e *ErrBlobNotFound) GRPCStatus() *status.Status {
	return getGRPCStatus(e, codes.NotFound, pbda.ErrorCode_BlobNotFound)
}

// GRPCStatus returns the gRPC status with details for an ErrBlobSizeOverLimit error.
func (e *ErrBlobSizeOverLimit) GRPCStatus() *status.Status {
	return getGRPCStatus(e, codes.ResourceExhausted, pbda.ErrorCode_BlobSizeOverLimit)
}

// GRPCStatus returns the gRPC status with details for an ErrTxTimedOut error.
func (e *ErrTxTimedOut) GRPCStatus() *status.Status {
	return getGRPCStatus(e, codes.DeadlineExceeded, pbda.ErrorCode_TxTimedOut)
}

// GRPCStatus returns the gRPC status with details for an ErrTxAlreadyInMempool error.
func (e *ErrTxAlreadyInMempool) GRPCStatus() *status.Status {
	return getGRPCStatus(e, codes.AlreadyExists, pbda.ErrorCode_TxAlreadyInMempool)
}

// GRPCStatus returns the gRPC status with details for an ErrTxIncorrectAccountSequence error.
func (e *ErrTxIncorrectAccountSequence) GRPCStatus() *status.Status {
	return getGRPCStatus(e, codes.InvalidArgument, pbda.ErrorCode_TxIncorrectAccountSequence)
}

// GRPCStatus returns the gRPC status with details for an ErrTxTooLarge error.
func (e *ErrTxTooLarge) GRPCStatus() *status.Status {
	return getGRPCStatus(e, codes.ResourceExhausted, pbda.ErrorCode_TxTooLarge)
}

// GRPCStatus returns the gRPC status with details for an ErrContextDeadline error.
func (e *ErrContextDeadline) GRPCStatus() *status.Status {
	return getGRPCStatus(e, codes.DeadlineExceeded, pbda.ErrorCode_ContextDeadline)
}

// GRPCStatus returns the gRPC status with details for an ErrFutureHeight error.
func (e *ErrFutureHeight) GRPCStatus() *status.Status {
	return getGRPCStatus(e, codes.OutOfRange, pbda.ErrorCode_FutureHeight)
}
