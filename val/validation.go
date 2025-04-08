package val

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidUsername = regexp.MustCompile(`^[a-z0-9_]+$`).MatchString
	isValidFullname = regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString
)

func ValidateString(value string, minLength, maxLength int) error {
	n := len(value)
	if n < minLength || n > maxLength {
		return fmt.Errorf("length must be between %d and %d", minLength, maxLength)
	}
	return nil
}

func ValidateUsername(value string) error {
	if err := ValidateString(value, 3, 50); err != nil {
		return err
	}
	if !isValidUsername(value) {
		return fmt.Errorf("username must contain only letters, numbers and underscores")
	}
	return nil
}

func ValidateFullname(value string) error {
	if err := ValidateString(value, 3, 50); err != nil {
		return err
	}
	if !isValidFullname(value) {
		return fmt.Errorf("fullname must contain only letters and spaces")
	}
	return nil
}

func ValidatePassword(value string) error {
	return ValidateString(value, 6, 100)
}

func ValidateEmail(value string) error {
	if err := ValidateString(value, 3, 320); err != nil {
		return err
	}
	_, err := mail.ParseAddress(value)
	if err != nil {
		return fmt.Errorf("invalid email address")
	}
	return nil
}

func ValidateEmailId(value int64) error {
	if value <= 0 {
		return fmt.Errorf("email id must be positive integer")
	}
	return nil
}

func ValidateSecrectCode(value string) error {
	return ValidateString(value, 32, 128)
}
