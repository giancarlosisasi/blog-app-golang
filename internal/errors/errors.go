package errors

const (
	CREATE_USER_INVALID_PASSWORD_ERROR = iota
	CREATE_USER_INVALID_EMAIL_ERROR
	CANNOT_GET_USER_WITH_EMAIL_ERROR
	USER_ALREADY_EXISTS_IN_DB_ERROR
	CREATE_USER_FAILED_TO_CREATE_ERROR
)

type CustomError struct {
	Code int    `json:"code"`
	Msg  string `json:"message"`
}

func (e CustomError) Error() string {
	return e.Msg
}

func NewCustomError(code int, msg string) *CustomError {
	return &CustomError{
		Code: code,
		Msg:  msg,
	}
}
