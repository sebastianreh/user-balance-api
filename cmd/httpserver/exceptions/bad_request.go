package exceptions

import "net/http"

type BadRequestException struct {
	HTTPCode   int    `json:"code" default:"400"`
	ErrMessage string `json:"error" default:"error message"`
}

func (exception BadRequestException) Error() string {
	return exception.ErrMessage
}

func (exception BadRequestException) Code() int {
	return exception.HTTPCode
}

func (exception BadRequestException) IsBadRequest() bool {
	return true
}

func NewBadRequestException(message string) BadRequestException {
	return BadRequestException{ErrMessage: message, HTTPCode: http.StatusBadRequest}
}
