package validation

import "errors"

var (
	ErrEmailRequired = errors.New("email is required")
	ErrEmailInvalid  = errors.New("email is invalid")
	ErrEmailTooLong  = errors.New("email is too long")

	ErrPasswordRequired  = errors.New("password is required")
	ErrPasswordTooShort  = errors.New("password must be at least 8 characters")
	ErrPasswordTooLong   = errors.New("password must be less than 72 characters")
	ErrPasswordNoUpper   = errors.New("password must contain at least one uppercase letter")
	ErrPasswordNoLower   = errors.New("password must contain at least one lowercase letter")
	ErrPasswordNoNumber  = errors.New("password must contains at least one number")
	ErrPasswordNoSpecial = errors.New("password must contains at least one special character")
	ErrPasswordCommon    = errors.New("password is too common")

	ErrAppIDRequired = errors.New("app_id is required")
)
