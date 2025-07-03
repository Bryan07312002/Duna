package auth

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

func (u *User) Password() (Password, error) {
	if u.password.alreadyHashed {
		return u.password, nil
	}

	return u.password.Hash()
}
