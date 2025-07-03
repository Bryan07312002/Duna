package auth

import (
	"errors"
	"unicode"
)

type HashStrategy interface {
	Encode(str string) (string, error)
	Compare(enconded, str string) bool
}

type Password struct {
	str           string
	hash          HashStrategy
	alreadyHashed bool
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
	}

	return Password{
		str:           password,
		hash:          hash,
		alreadyHashed: alreadyHashed,
	}, nil
}

func (p *Password) Hash() (Password, error) {
	hashed, err := p.hash.Encode(p.str)
	if err != nil {
		return Password{}, err
	}

	return NewPassword(hashed, p.alreadyHashed, p.hash)
}

func (p *Password) Compare(incomingPassword string) bool {
	if !p.alreadyHashed {
		return false
	}

	return p.hash.Compare(p.str, incomingPassword)
}

// Always obfuscate when converted to string
func (p Password) String() string {
	return "********"
}

// Obfuscate in %#v formatting
func (p Password) GoString() string {
	return "********"
}
