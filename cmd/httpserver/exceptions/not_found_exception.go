package exceptions

import "net/http"

type NotFoundException struct {
	HTTPCode   int    `json:"code" default:"404"`
	ErrMessage string `json:"error" default:"error message"`
}

func (exception NotFoundException) Error() string {
	return exception.ErrMessage
}

func (exception NotFoundException) Code() int {
	return exception.HTTPCode
}

func NewNotFoundException(message string) NotFoundException {
	return NotFoundException{ErrMessage: message, HTTPCode: http.StatusNotFound}
}
