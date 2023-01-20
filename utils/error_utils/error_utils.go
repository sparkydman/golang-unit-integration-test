package error_utils

import (
	"encoding/json"
	"net/http"
)

type messageErr struct {
	ErrMessage string `json:"message"`
	ErrStatus  int    `json:"status"`
	ErrError   string `json:"error"`
}

type MessageErr interface {
	Message() string
	Status() int
	Error() string
}

func (e *messageErr) Error() string {
	return e.ErrError
}

func (e *messageErr) Status() int {
	return e.ErrStatus
}

func (e *messageErr) Message() string {
	return e.ErrMessage
}

func NewNotFoundErrorMessage(message string) MessageErr {
	return &messageErr{
		ErrMessage: message,
		ErrStatus:  http.StatusNotFound,
		ErrError: "not_found",
	}
}

func NewBadRequestErrorMessage(message string) MessageErr {
	return &messageErr{
		ErrMessage: message,
		ErrStatus: http.StatusBadRequest,
		ErrError: "bad_request",
	}
}

func NewUnprocessibleEntityErrorMessages(message string) MessageErr { 
	return &messageErr{
		ErrMessage: message,
		ErrStatus: http.StatusUnprocessableEntity,
		ErrError: "invalid_request",
	}
}

func NewApiErrorMessageFromBytes(body []byte) (MessageErr, error){
	var result messageErr
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func NewInternalServerErrorMessage(message string) MessageErr {
	return &messageErr{
		ErrMessage: message,
		ErrStatus: http.StatusInternalServerError,
		ErrError: "server_error",
	}
}