package errs

import "net/http"

var (
	ErrOk                   = NewErrorType(http.StatusOK)
	ErrInternal             = NewErrorType(http.StatusInternalServerError)
	ErrNotFound             = NewErrorType(http.StatusNotFound)
	ErrBadRequest           = NewErrorType(http.StatusBadRequest)
	ErrUnauthorized         = NewErrorType(http.StatusUnauthorized)
	ErrForbidden            = NewErrorType(http.StatusForbidden)
	ErrMethodNotAllowed     = NewErrorType(http.StatusMethodNotAllowed)
	ErrConflict             = NewErrorType(http.StatusConflict)
	ErrNocontent            = NewErrorType(http.StatusNoContent)
	ErrUnsupportedMediaType = NewErrorType(http.StatusUnsupportedMediaType)
	ErrTooManyRequests      = NewErrorType(http.StatusTooManyRequests)
	ErrTimeout              = NewErrorType(http.StatusRequestTimeout)
	ErrPaymentRequired      = NewErrorType(http.StatusPaymentRequired)
	ErrUnprocessableEntity  = NewErrorType(http.StatusUnprocessableEntity)
	ErrNotImplemented       = NewErrorType(http.StatusNotImplemented)
)

type ErrorType struct {
	HTTPStatus int
}

func NewErrorType(httpStatus int) *ErrorType {
	return &ErrorType{HTTPStatus: httpStatus}
}

type AppError struct {
	Type    *ErrorType
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(errorType *ErrorType, message string) *AppError {
	return &AppError{
		Type:    errorType,
		Message: message,
	}
}
