package exceptions

import "net/http"

type NotFoundException interface {
	Error() string
	IsNotFoundError() bool
	Code() int
}

type notFoundException struct {
	ErrMessage string `json:"error"`
	HttpCode   int
}

func (exception *notFoundException) Error() string {
	return exception.ErrMessage
}

func (exception *notFoundException) Code() int {
	return exception.HttpCode
}

func (exception *notFoundException) IsNotFoundError() bool {
	return true
}

func NewNotFoundException(message string) NotFoundException {
	return &notFoundException{ErrMessage: message, HttpCode: http.StatusNotFound}
}
