package storage

import "errors"

var (
	ErrLoginExists  = errors.New("email already exists")
	ErrUserExists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")

	ErrAppNotFound = errors.New("app not found")
)
