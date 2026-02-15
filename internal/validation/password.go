package validation

import "unicode"

var (
	minPasswordLength = 8
	maxPasswordLength = 72
)

var commonPassword = map[string]bool{
	"password": true,
	"12345678": true,
	"qwerty12": true,
	"admin12":  true,
}

func ValidatePassword(password string) error {

	if password == "" {
		return ErrPasswordRequired
	}

	if len(password) < minPasswordLength {
		return ErrPasswordTooShort
	}

	if len(password) > maxPasswordLength {
		return ErrPasswordTooLong
	}

	if commonPassword[password] {
		return ErrPasswordCommon
	}

	var hasLower, hasUpper, hasNumber, hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasLower {
		return ErrPasswordNoLower
	}
	if !hasUpper {
		return ErrPasswordNoUpper
	}
	if !hasNumber {
		return ErrPasswordNoNumber
	}
	if !hasSpecial {
		return ErrPasswordNoSpecial
	}

	return nil
}
