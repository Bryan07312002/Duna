package postgres

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// PostgresConfig holds PostgreSQL connection configuration
type PostgresConfig struct {
	host     string
	port     int
	user     string
	password string
	dbName   string
	sslMode  string
}

const (
	PG_CREDS_FILE_ENV_VAR = "PG_CREDS_FILE"
	PG_HOST_ENV_VAR       = "PG_HOST"
	PG_PORT_ENV_VAR       = "PG_PORT"
	PG_DB_NAME_ENV_VAR    = "PG_DBNAME"
	PG_SSLMODE_ENV_VAR    = "PG_SSLMODE"
)

// NewPostgresConfig creates a new configuration instance
func NewPostgresConfig() (*PostgresConfig, error) {
	credsPath := getEnv(PG_CREDS_FILE_ENV_VAR, "")
	if credsPath == "" {
		return nil, errors.New(PG_CREDS_FILE_ENV_VAR +
			" environment variable not set")
	}

	user, password, err := readCredentialsFromFile(credsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read credentials: %w", err)
	}

	return &PostgresConfig{
		host:     getEnv(PG_HOST_ENV_VAR, "localhost"),
		port:     getEnvAsInt(PG_PORT_ENV_VAR, 5432),
		user:     user,
		password: password,
		dbName:   getEnv(PG_DB_NAME_ENV_VAR, "postgres"),
		sslMode:  getEnv(PG_SSLMODE_ENV_VAR, "disable"),
	}, nil
}

// ConnectionString returns the formatted connection string
func (c *PostgresConfig) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.host,
		c.port,
		c.user,
		c.password,
		c.dbName,
		c.sslMode)
}

// Getter methods for individual properties (optional)
func (c *PostgresConfig) Host() string    { return c.host }
func (c *PostgresConfig) Port() int       { return c.port }
func (c *PostgresConfig) User() string    { return c.user }
func (c *PostgresConfig) DBName() string  { return c.dbName }
func (c *PostgresConfig) SSLMode() string { return c.sslMode }

func readCredentialsFromFile(path string) (string, string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", "", fmt.Errorf("invalid file path: %w", err)
	}

	// Verify file permissions (optional but recommended)
	// if err := verifyFilePermissions(absPath); err != nil {
	// 	return "", "", err
	// }

	file, err := os.Open(absPath)
	if err != nil {
		return "", "", fmt.Errorf("failed to open credentials file: %w", err)
	}
	defer file.Close()

	var user, password string
	scanner := bufio.NewScanner(file)

	// Read first line for username
	if scanner.Scan() {
		user = strings.TrimSpace(scanner.Text())
	}

	// Read second line for password
	if scanner.Scan() {
		password = strings.TrimSpace(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return "", "", fmt.Errorf("error reading credentials from %s file: %w",
			PG_CREDS_FILE_ENV_VAR, err)
	}

	if user == "" || password == "" {
		return "", "", fmt.Errorf("invalid credentials format in file %s",
			PG_CREDS_FILE_ENV_VAR)
	}

	return user, password, nil
}

// verifyFilePermissions checks if file has secure permissions
func verifyFilePermissions(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.Mode().Perm()&0077 != 0 {
		return fmt.Errorf("insecure file permissions on %s (should be 600)", path)
	}

	return nil
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	strValue := getEnv(key, "")
	if value, err := strconv.Atoi(strValue); err == nil {
		return value
	}
	return defaultValue
}
