package client

import (
	"encoding/json"
	"net/http"
)

type Response[T any] struct {
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message,omitempty"`
	User    *T     `json:"user,omitempty"`
	Users   *T     `json:"users,omitempty"`
	Meta    *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
}

func WriteSuccess[T any](w http.ResponseWriter, status int, code string, user T) {
	resp := NewSuccessResponse(user, code)
	writeJSON(w, status, resp)
}

func WriteSuccessList[T any](w http.ResponseWriter, status int, code string, users T, meta Meta) {
	resp := NewSuccessResponseList(users, code)
	resp.Meta = &meta
	writeJSON(w, status, resp)
}

func WriteMessage(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, NewMessageResponse(message, code))
}

func WriteFail(w http.ResponseWriter, status int, code, message string) {
	resp := NewFailResponse(message, code)
	writeJSON(w, status, resp)
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(payload)
}

func NewSuccessResponse[T any](user T, code string) Response[T] {
	return Response[T]{
		Status: "success",
		Code:   code,
		User:   &user,
	}
}

func NewMessageResponse(message, code string) Response[any] {
	return Response[any]{
		Status:  "success",
		Code:    code,
		Message: message,
	}
}

func NewSuccessResponseList[T any](users T, code string) Response[T] {
	return Response[T]{
		Status: "success",
		Code:   code,
		Users:  &users,
	}
}

func NewFailResponse(message, code string) Response[any] {
	return Response[any]{
		Status:  "fail",
		Code:    code,
		Message: message,
	}
}
