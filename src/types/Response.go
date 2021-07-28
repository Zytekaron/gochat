package types

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success   bool        `json:"success"`
	Error     string      `json:"error,omitempty"`
	ErrorCode int         `json:"error_code,omitempty"`
	Data      interface{} `json:"data,omitempty"`
}

func newSuccess(payload interface{}) *Response {
	return &Response{
		true,
		"",
		0,
		payload,
	}
}

func newError(err string, errCode int) *Response {
	return &Response{
		false,
		err,
		errCode,
		nil,
	}
}

func WriteSuccessJson(w http.ResponseWriter, payload interface{}) error {
	res := newSuccess(payload)
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(res)
}

func WriteErrorJson(w http.ResponseWriter, statusCode int, errCode int, err string) error {
	res := newError(err, errCode)
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(res)
}

func WriteSuccess(w http.ResponseWriter, text string) error {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(text))
	return err
}

func WriteError(w http.ResponseWriter, statusCode int, text string) error {
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(text))
	return err
}
