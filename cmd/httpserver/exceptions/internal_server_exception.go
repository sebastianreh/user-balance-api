package exceptions

import "net/http"

type InternalServerException interface {
	Error() string
	isInternalServerException() bool
	Code() int
}

type internalServerException struct {
	ErrMessage string `json:"error"`
	HttpCode   int
}

func (exception *internalServerException) Error() string {
	return exception.ErrMessage
}

func (exception *internalServerException) Code() int {
	return exception.HttpCode
}

func (exception *internalServerException) isInternalServerException() bool {
	return true
}

func NewInternalServerException(message string) InternalServerException {
	return &internalServerException{ErrMessage: message, HttpCode: http.StatusInternalServerError}
}
