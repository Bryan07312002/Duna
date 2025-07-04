package database

import "duna/internal/user"

type MockDatabase struct {
	FuncMigrate           func() error
	FuncInsertUser        func(user user.User) error
	FuncGetUserByUsername func(username string, hash user.HashStrategy) (user.User, error)
}

func (m *MockDatabase) Migrate() error {
	return m.FuncMigrate()
}

func (m *MockDatabase) InsertUser(user user.User) error {
	return m.FuncInsertUser(user)
}

func (m *MockDatabase) GetUserByUsername(username string, hash user.HashStrategy) (user.User, error) {
	return m.FuncGetUserByUsername(username, hash)
}
