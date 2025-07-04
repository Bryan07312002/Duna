package database

import (
	"duna/internal/database/postgres"

	"duna/internal/user"
)

type Database interface {
	Migrate() error
	InsertUser(user user.User) error
	GetUserByUsername(username string, hash user.HashStrategy) (user.User, error)
}

func NewDatabase() (Database, error) {
	config, err := postgres.NewPostgresConfig()
	if err != nil {
		return nil, err
	}

	return postgres.NewPostgresDatabase(*config)
}
