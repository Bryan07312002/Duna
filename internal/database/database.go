package database

import "duna/internal/database/postgres"

type Database interface {
	Migrate() error
}

func NewDatabase() (Database, error) {
	config, err := postgres.NewPostgresConfig()
	if err != nil {
		return nil, err
	}

	return postgres.NewPostgresDatabase(*config)
}
