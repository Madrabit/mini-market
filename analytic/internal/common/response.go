package common

import (
	"encoding/json"
	"net/http"
)

type Response[T any] struct {
	Success bool    `json:"success"`
	Message *string `json:"message,omitempty"`
	Data    T       `json:"data"`
}

func ErrResponse(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(Response[any]{
		Success: false,
		Message: &msg,
		Data:    nil,
	})
}

func okResponseInternal[T any](w http.ResponseWriter, data T, msg *string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(Response[T]{
		Success: true,
		Message: msg,
		Data:    data,
	})
}

func OkResponse[T any](w http.ResponseWriter, data T) {
	okResponseInternal(w, data, nil)
}

func OkResponseMsg[T any](w http.ResponseWriter, data T, msg string) {
	okResponseInternal(w, data, &msg)
}
