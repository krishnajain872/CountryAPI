package error

import (
	"net/http"
	"time"

	"github.com/krishnajain872/country-search-api-cache/pkg/types"
)

/* =====================
   Constructors
   ===================== */

func New(
	errType types.ErrorType,
	message string,
	status int,
	err error,
) *types.AppError {
	return &types.AppError{
		Type:       errType,
		Message:    message,
		StatusCode: status,
		Err:        err,
		Context:    map[string]any{},
		Timestamp:  time.Now(),
	}
}

/* =====================
   Helpers
   ===================== */

func Validation(field, message string) *types.AppError {
	return New(
		types.Validation,
		field+" "+message,
		http.StatusBadRequest,
		types.ErrInvalidInput,
	).WithContext("field", field)
}

func ResourceNotFound(resource string) *types.AppError {
	return New(
		types.NotFound,
		resource+" not found",
		http.StatusNotFound,
		types.ErrNotFound,
	)
}

func External(message string, err error) *types.AppError {
	return New(types.External, message, http.StatusBadGateway, err)
}

func Timeout(message string) *types.AppError {
	return New(types.Timeout, message, http.StatusGatewayTimeout, types.ErrTimeout)
}

func InternalErr(message string, err error) *types.AppError {
	return New(types.Internal, message, http.StatusInternalServerError, err)
}

/* =====================
   Normalization
   ===================== */

func Get(err error) *types.AppError {
	if err == nil {
		return nil
	}
	if appErr, ok := err.(*types.AppError); ok {
		return appErr
	}
	return InternalErr("an unexpected error occurred", err)
}
