package exceptions

import "net/http"

type BadRequestException interface {
	Error() string
	IsBadRequest() bool
	Code() int
}

type badRequestException struct {
	ErrMessage string `json:"error"`
	HttpCode   int
}

func (exception *badRequestException) Error() string {
	return exception.ErrMessage
}

func (exception *badRequestException) Code() int {
	return exception.HttpCode
}

func (exception *badRequestException) IsBadRequest() bool {
	return true
}

func NewBadRequestException(message string) BadRequestException {
	return &badRequestException{ErrMessage: message, HttpCode: http.StatusBadRequest}
}
