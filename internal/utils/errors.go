package utils

import "errors"

const (
	CREATE_USER_INVALID_PASSWORD_ERROR = iota
	CREATE_USER_INVALID_EMAIL_ERROR
	CANNOT_GET_USER_WITH_EMAIL_ERROR
	USER_ALREADY_EXISTS_IN_DB_ERROR
	CREATE_USER_FAILED_TO_CREATE_ERROR
	INTERNAL_SERVER_ERROR
	CREATE_USER_INVALID_BODY_REQUEST_ERROR
	USER_NOT_FOUND_ERROR
	USER_INVALID_PASSWORD_ERROR
	CREATE_POST_INVALID_BODY_ERROR
	USER_NOT_AUTHORIZED_ERROR
	BAD_REQUEST_ERROR
	NOT_FOUND_ERROR
)

var ErrInvalidUUID = errors.New("invalid UUID format")
var ErrResourceNotFoundInDB = errors.New("resource not found in db")
var ErrAuthTokenExpired = errors.New("jwt token expired")
var ErrAuthTokenInvalid = errors.New("jwt token is invalid")
var ErrAuthTokenMissing = errors.New("jwt token is missing")

type AppError struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

func (e AppError) Error() string {
	return e.Msg
}

func NewCustomError(code int, msg string) *AppError {
	return &AppError{
		Code: code,
		Msg:  msg,
	}
}
