package exceptions

import "net/http"

type InternalServerException struct {
	HTTPCode   int    `json:"code" default:"500"`
	ErrMessage string `json:"error" default:"error message"`
}

func (exception InternalServerException) Error() string {
	return exception.ErrMessage
}

func (exception InternalServerException) Code() int {
	return exception.HTTPCode
}

func NewInternalServerException(message string) InternalServerException {
	return InternalServerException{ErrMessage: message, HTTPCode: http.StatusInternalServerError}
}
