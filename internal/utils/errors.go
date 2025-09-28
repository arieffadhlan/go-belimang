package utils

type AppError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func (e AppError) Error() string {
	return e.Message
}

func NewNotFound(msg string) AppError {
	return AppError{StatusCode: 404, Message: msg}
}

func NewConflict(msg string) AppError {
	return AppError{StatusCode: 409, Message: msg}
}

func NewInternal(msg string) AppError {
	return AppError{StatusCode: 500, Message: msg}
}

func NewBadRequest(msg string) AppError {
	return AppError{StatusCode: 400, Message: msg}
}

func NewTooManyReq(msg string) AppError {
	return AppError{StatusCode: 429, Message: msg}
}
