package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

const MIGRATIONS_TABLE = "migrations"

func (p *PostgresDatabase) getAppliedMigrations() ([]*Migration, error) {
	selectQuery := fmt.Sprintf(
		"SELECT name, timestamp FROM %s ORDER BY timestamp",
		MIGRATIONS_TABLE,
	)

	rows, err := p.QuerySql(nil, selectQuery)
	if err != nil {
		return nil, err
	}

	var migrations []*Migration
	for rows.Next() {
		var name string
		var timestamp int64
		if err := rows.Scan(&name, &timestamp); err != nil {
			return nil, err
		}

		migration := NewMigration(name, timestamp, "", DefaultFileSystem{})
		migrations = append(migrations, migration)
	}

	return migrations, nil
}

func (p *PostgresDatabase) InsertMigration(tx *sql.Tx, migration *Migration) error {
	insertMigrationQuery := fmt.Sprintf(
		"INSERT INTO %s (name, timestamp) VALUES($1, $2)",
		MIGRATIONS_TABLE,
	)

	if _, err := p.ExecSql(
		tx,
		insertMigrationQuery,
		migration.name,
		migration.timestamp,
	); err != nil {
		return err
	}

	return nil
}

// FileSystem interface for file operations (makes it mockable)
type FileSystem interface {
	ReadFile(name string) ([]byte, error)
	ReadDir(name string) ([]os.DirEntry, error)
}

// DefaultFileSystem implements FileSystem using os package
type DefaultFileSystem struct{}

func (DefaultFileSystem) ReadFile(name string) ([]byte, error) {
	return os.ReadFile(name)
}

func (DefaultFileSystem) ReadDir(name string) ([]os.DirEntry, error) {
	return os.ReadDir(name)
}

type Migration struct {
	name      string
	timestamp int64
	path      string
	fs        FileSystem
}

// Read a directory searching for migrations folders, a migration folder
// is defined by <number>-<name>, all that does not match this style will be
// ignored
func ReadMigrationDir(dir string, fs FileSystem) ([]*Migration, error) {
	entries, err := fs.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var migrations []*Migration
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		splited := strings.Split(entry.Name(), "-")
		if len(splited) <= 1 {
			continue
		}

		timestamp, err := strconv.ParseInt(splited[0], 10, 64)
		if err != nil {
			return nil, err
		}

		migrations = append(migrations, NewMigration(
			strings.Join(splited[1:], "-"),
			timestamp,
			dir,
			fs,
		))
	}

	// Sort descending (newest first)
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].timestamp < migrations[j].timestamp
	})

	return migrations, nil
}

// Creates a new Migration instance with dependencies
func NewMigration(name string, timestamp int64, path string, fs FileSystem) *Migration {
	if fs == nil {
		fs = DefaultFileSystem{}
	}

	return &Migration{
		name:      name,
		timestamp: timestamp,
		path:      path,
		fs:        fs,
	}
}

func (m *Migration) FullName() string {
	return fmt.Sprintf("%d", m.timestamp) + "-" + m.name
}

func (m *Migration) GetUpQuery() (string, error) {
	query, err := m.fs.ReadFile(filepath.Join(m.path, m.FullName(), "up.sql"))
	if err != nil {
		return "", fmt.Errorf("failed to read up.sql: %w", err)
	}

	return string(query), nil
}

func (m *Migration) GetDownQuery() (string, error) {
	query, err := m.fs.ReadFile(filepath.Join(m.path, m.FullName(), "down.sql"))
	if err != nil {
		return "", fmt.Errorf("failed to read up.sql: %w", err)
	}

	return string(query), nil
}
