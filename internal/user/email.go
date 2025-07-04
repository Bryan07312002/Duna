package user

import (
	"errors"
	"regexp"
)

type Email string

func checkEmailValid(email string) error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(string(email)) {
		return errors.New("invalid email format")
	}

	return nil
}

func NewEmail(email string) (Email, error) {
	if err := checkEmailValid(email); err != nil {
		return "", err
	}

	return Email(email), nil
}
