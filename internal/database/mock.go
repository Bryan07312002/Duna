package database

import ( 
	"duna/internal/models"
	"duna/internal/hash"
)

type MockDatabase struct {
	FuncMigrate           func() error
	FuncInsertUser        func(user models.User) error
	FuncGetUserByUsername func(username string,
		hash hash.HashStrategy) (models.User, error)
}

func (m *MockDatabase) Migrate() error {
	return m.FuncMigrate()
}

func (m *MockDatabase) InsertUser(user models.User) error {
	return m.FuncInsertUser(user)
}

func (m *MockDatabase) GetUserByUsername(username string,
	hash hash.HashStrategy) (models.User, error) {
	return m.FuncGetUserByUsername(username, hash)
}
