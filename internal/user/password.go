package user

import (
	"errors"
	"unicode"
)

type HashStrategy interface {
	Encode(str string) (string, error)
	Compare(enconded, str string) bool
}

type Password struct {
	value string
	hash  HashStrategy
}

func isPasswordValid(p string) error {
	if len(p) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	hasUpper := false
	hasLower := false
	hasNumber := false
	hasSpecial := false

	for _, c := range p {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	if !hasUpper {
		return errors.New("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return errors.New("password must contain at least one lowercase letter")
	}
	if !hasNumber {
		return errors.New("password must contain at least one number")
	}
	if !hasSpecial {
		return errors.New("password must contain at least one special character")
	}

	return nil
}

func NewPassword(password string, alreadyHashed bool,
	hash HashStrategy) (Password, error) {
	if !alreadyHashed {
		if err := isPasswordValid(password); err != nil {
			return Password{}, err
		}

		hasedPassword, err := hash.Encode(password)
		if err != nil {
			return Password{}, err
		}
		password = hasedPassword
	}

	return Password{
		value:         password,
		hash:          hash,
	}, nil
}

func (p Password) Compare(incomingPassword string) bool {
	return p.hash.Compare(p.value, incomingPassword)
}

// Always obfuscate when converted to string
func (p Password) String() string {
	return "********"
}

// Obfuscate in %#v formatting
func (p Password) GoString() string {
	return "********"
}
