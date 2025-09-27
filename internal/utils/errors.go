package utils

import (
	"net/http"
)

type AppError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func (e AppError) Error() string {
	return e.Message
}

func NewNotFound(msg string) AppError {
	return AppError{StatusCode: http.StatusNotFound, Message: msg}
}

func NewConflict(msg string) AppError {
	return AppError{StatusCode: http.StatusConflict, Message: msg}
}

func NewInternal(msg string) AppError {
	return AppError{StatusCode: http.StatusInternalServerError, Message: msg}
}

func NewBadRequest(msg string) AppError {
	return AppError{StatusCode: http.StatusBadRequest, Message: msg}
}
