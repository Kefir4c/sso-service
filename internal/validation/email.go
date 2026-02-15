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
		return ErrEmailInvilid
	}

	local, domain, _ := strings.Cut(addr.Address, "@")

	if len(local) > 64 || len(local) == 0 {
		return ErrEmailInvilid
	}

	if !strings.Contains(domain, ".") ||
		strings.HasPrefix(domain, ".") ||
		strings.HasSuffix(domain, ".") {
		return ErrEmailInvilid
	}

	return nil
}
