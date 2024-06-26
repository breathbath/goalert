package validate

import "github.com/breathbath/goalert/validation"

// Username will validate a username to ensure it is between 3 and 24 characters,
// and only contains lower-case ASCII letters, numbers, and '-', '_', and '.'.
func Username(fname, value string) error {
	b := []byte(value)
	l := len(b)
	if l < 3 {
		return validation.NewFieldError(fname, "must be at least 3 characters")
	}
	if l > 24 {
		return validation.NewFieldError(fname, "cannot be more than 24 characters")
	}

	for _, c := range value {
		if c >= 'a' && c <= 'z' {
			continue
		}
		if c >= '0' && c <= '9' {
			continue
		}
		if c == '-' || c == '_' || c == '.' {
			continue
		}

		return validation.NewFieldError(fname, "can only contain lower-case letters, digits, and the characters '-', '_', and '.'")
	}

	return nil
}
