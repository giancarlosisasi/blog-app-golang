package models

type ErrorResponse struct {
	Msg  string `json:"message"`
	Code string `json:"code"`
}
