package postgres

import (
	"duna/internal/models"
	"duna/internal/hash"
	"fmt"
)

const USERS_TABLE = "users"

func (p *PostgresDatabase) InsertUser(user models.User) error {
	insertQuery := fmt.Sprintf(
		"INSERT INTO %s (uuid, username, email, password) VALUES($1, $2, $3, $4)",
		USERS_TABLE,
	)

	if _, err := p.ExecSql(
		nil,
		insertQuery,
		user.UUID,
		user.Username,
		user.Email,
	); err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabase) GetUserByUsername(
	queryUsername string, hash hash.HashStrategy) (models.User, error) {
	query := fmt.Sprintf(
		"SELECT uuid, username, email, password FROM %s WHERE username = $1 LIMIT 1",
		USERS_TABLE,
	)

	rows, err := p.QuerySql(nil, query, queryUsername)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to query user: %w", err)
	}
	defer rows.Close()

	if !rows.Next() {
		return models.User{}, fmt.Errorf("user not found")
	}

	var (
		uuid     string
		username string
		email    string
		password string
	)

	if err := rows.Scan(&uuid, &username, &email, &password); err != nil {
		return models.User{}, fmt.Errorf("failed to scan user data: %w", err)
	}

	// assumes password in database is already hashed
	User, err := models.NewUserFromPrimitives(
		uuid, username, email, password, true, hash)
	// errors if is a database error
	if err != nil {
		return models.User{}, fmt.Errorf("invalid data in database: %w", err)
	}

	return User, nil
}
