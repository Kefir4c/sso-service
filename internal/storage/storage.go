package storage

import "errors"

var (
	ErrLoginExists  = errors.New("login already exists")
	ErrUserNotFound = errors.New("user not found")

	ErrAppNotFound = errors.New("app not found")
)
