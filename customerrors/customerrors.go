package customerrors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Error struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

var (
	ErrDecodeBody = Error{Message: "Decode request body error", Code: http.StatusBadRequest}
	ErrParseQuery = Error{Message: "Parse query params error", Code: http.StatusBadRequest}
	ErrNotFound   = Error{Message: "Resource not found ", Code: http.StatusNotFound}
	ErrInternal   = Error{Message: "Internal server error", Code: http.StatusInternalServerError}
)

func New(err error, code int) Error {
	return Error{Message: err.Error(), Code: code}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s (code %d)", e.Message, e.Code)
}

func (e *Error) GetJSON() []byte {
	var b, _ = json.Marshal(e)
	return b
}
