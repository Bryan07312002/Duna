package database

import (
	"duna/internal/database/postgres"
	"duna/internal/hash"

	"duna/internal/models"
)

type Database interface {
	Migrate() error
	// users
	InsertUser(user models.User) error
	GetUserByUsername(username string, hash hash.HashStrategy) (models.User, error)

	// matches
	InsertMatch(match models.Match) error
}

func NewDatabase() (Database, error) {
	config, err := postgres.NewPostgresConfig()
	if err != nil {
		return nil, err
	}

	return postgres.NewPostgresDatabase(*config)
}
