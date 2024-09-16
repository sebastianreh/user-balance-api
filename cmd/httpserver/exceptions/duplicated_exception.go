package exceptions

import "net/http"

type DuplicatedException struct {
	HTTPCode   int    `json:"code" default:"409"`
	ErrMessage string `json:"error" default:"error message"`
}

type Causes struct {
	Code string `json:"code"`
}

func (exception DuplicatedException) Error() string {
	return exception.ErrMessage
}

func (exception DuplicatedException) Code() int {
	return exception.HTTPCode
}

func NewDuplicatedException(message string) DuplicatedException {
	return DuplicatedException{ErrMessage: message, HTTPCode: http.StatusConflict}
}
