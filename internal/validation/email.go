package validation

import (
	"net/mail"
	"strings"
)

const maxEmailLength = 254

func ValidateEmail(email string) error {
	if email == "" {
		return ErrEmailRequired
	}
	if len(email) > maxEmailLength {
		return ErrEmailTooLong
	}

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return ErrEmailInvalid
	}

	local, domain, found := strings.Cut(addr.Address, "@")
	if !found {
		return ErrEmailInvalid
	}

	if len(local) > 64 || len(local) == 0 {
		return ErrEmailInvalid
	}

	if !strings.Contains(domain, ".") ||
		strings.HasPrefix(domain, ".") ||
		strings.HasSuffix(domain, ".") {
		return ErrEmailInvalid
	}

	return nil
}
