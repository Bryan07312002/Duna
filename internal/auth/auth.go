package auth

import "errors"

func Authenticate(username, password string,
	user User) (string, string, error) {
	if username != user.Username {
		return "", "", errors.New("usernames do not match")
	}

	userPassword, err := user.Password()
	if err != nil {
		return "", "", err
	}

	if !userPassword.Compare(password) {
		return "", "", errors.New("passwords not match")
	}
}
