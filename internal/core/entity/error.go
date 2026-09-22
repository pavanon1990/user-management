package entity

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrNoFieldsToUpdate   = errors.New("no fields to update")
	ErrInvalidCredentials = errors.New("invalid email or password")
)
