package postgres

import (
	"duna/internal/auth"
	"fmt"
)

const USERS_TABLE = "users"

func (p *PostgresDatabase) InsertUser(user auth.User) {
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
