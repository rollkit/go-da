package grpc

import (
	"errors"

	"google.golang.org/grpc/status"

	"github.com/rollkit/go-da"
	pbda "github.com/rollkit/go-da/types/pb/da"
)

func tryToMapError(err error) error {
	if err == nil {
		return nil
	}

	s, ok := status.FromError(err)
	if ok {
		details := s.Proto().Details
		if len(details) == 1 {
			var errorDetail pbda.ErrorDetails
			unmarshalError := errorDetail.Unmarshal(details[0].Value)
			if unmarshalError != nil {
				return err
			}
			return errorForCode(errorDetail.Code)
		}
	}
	return err
}

func errorForCode(code pbda.ErrorCode) error {
	switch code {
	case pbda.ErrorCode_BlobNotFound:
		return &da.ErrBlobNotFound{}
	case pbda.ErrorCode_BlobSizeOverLimit:
		return &da.ErrBlobSizeOverLimit{}
	case pbda.ErrorCode_TxTimedOut:
		return &da.ErrTxTimedOut{}
	case pbda.ErrorCode_TxAlreadyInMempool:
		return &da.ErrTxAlreadyInMempool{}
	case pbda.ErrorCode_TxIncorrectAccountSequence:
		return &da.ErrTxIncorrectAccountSequence{}
	case pbda.ErrorCode_TxTooLarge:
		return &da.ErrTxTooLarge{}
	case pbda.ErrorCode_ContextDeadline:
		return &da.ErrContextDeadline{}
	case pbda.ErrorCode_FutureHeight:
		return &da.ErrFutureHeight{}
	default:
		return errors.New("unknown error code")
	}
}
