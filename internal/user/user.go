package user

type UUIDStrategy interface {
	New() string
}

type User struct {
	UUID     string `json:"uuid"`
	Username string `json:"username"`
	Email    Email  `json:"email"`
	password Password
}

func NewUser(UUID, Username string, Email Email, Password Password) User {
	return User{
		UUID:     UUID,
		Username: Username,
		Email:    Email,
		password: Password,
	}
}

func NewUserFromPrimitives(UUID, Username, Email, Password string,
	passwordAlreadyHashed bool, hash HashStrategy) (User, error) {
	emailObj, err := NewEmail(Email)
	if err != nil {
		return User{}, err
	}

	passwordObj, err := NewPassword(Password, true, nil)
	if err != nil {
		return User{}, err
	}

	return NewUser(UUID, Username, emailObj, passwordObj), nil
}

func (u *User) Password() Password {
	return u.password
}
