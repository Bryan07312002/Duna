package postgres

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pkg/errors"
)

type PostgresDatabase struct {
	DB     *sql.DB
	config PostgresConfig
}

func NewPostgresDatabase(config PostgresConfig) (*PostgresDatabase, error) {
	defaultDb, err := sql.Open("pgx", config.ConnectionString())
	if err != nil {
		return nil, err
	}

	return &PostgresDatabase{DB: defaultDb, config: config}, nil
}

func (p *PostgresDatabase) BeginTransaction() (*sql.Tx, error) {
	tx, err := p.DB.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "unable to begin transaction")
	}
	return tx, nil
}

func (p *PostgresDatabase) CommitTransaction(tx *sql.Tx) error {
	if tx == nil {
		return errors.New("nil transaction")
	}

	err := tx.Commit()
	if err != nil {
		// Check if the transaction is already done
		if errors.Is(err, sql.ErrTxDone) {
			return errors.Wrap(err, "transaction already committed or rolled back")
		}

		return errors.Wrap(err, "transaction commit failed")
	}

	return nil
}

func (p *PostgresDatabase) RollbackTransaction(tx *sql.Tx) error {
	err := tx.Rollback()
	if err != nil {
		return errors.Wrap(err, "unable to rollback transaction")
	}
	return nil
}

func (p *PostgresDatabase) ExecSql(
	tx *sql.Tx,
	sql string,
	args ...any,
) (sql.Result, error) {
	if tx != nil {
		result, err := tx.Exec(sql, args...)
		if err != nil {
			return nil, errors.Wrap(err, "unable to execute SQL")
		}
		return result, nil
	}

	result, err := p.DB.Exec(sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "unable to execute SQL")
	}

	return result, nil
}

func (p *PostgresDatabase) QuerySql(
	tx *sql.Tx,
	sql string,
	args ...any,
) (*sql.Rows, error) {
	if tx != nil {
		result, err := tx.Query(sql, args...)
		if err != nil {
			return nil, errors.Wrap(err, "unable to execute SQL")
		}
		return result, nil
	}

	rows, err := p.DB.Query(sql, args...)
	if err != nil {
		return nil, errors.Wrap(err, "unable to execute SQL")
	}
	return rows, nil
}

// TODO: should check if already applied migrations match with repository ones
func (p *PostgresDatabase) Migrate() error {
	if err := p.ensureMigrationsTableExits(); err != nil {
		return err
	}

	migrations, err := ReadMigrationDir(
		"internal/database/migrations", DefaultFileSystem{})
	if err != nil {
		return err
	}

	appliedMigrations, err := p.getAppliedMigrations()
	if err != nil {
		return err
	}

	tx, err := p.BeginTransaction()
	if err != nil {
		return nil
	}

	if err := p.execUpMigrations(migrations[len(appliedMigrations):]); err != nil {
		return err
	}

	if err := p.CommitTransaction(tx); err != nil {
		return err
	}

	return nil
}

func (p *PostgresDatabase) execUpMigrations(
	migrations []*Migration,
) error {
	for _, migration := range migrations {
		query, err := migration.GetUpQuery()
		if err != nil {
			return err
		}

		if _, err := p.ExecSql(nil, query); err != nil {
			return err
		}

		tx, err := p.BeginTransaction()
		if err != nil {
			return err
		}

		p.InsertMigration(tx, migration)
		if err := p.CommitTransaction(tx); err != nil {
			return err
		}
	}

	return nil
}

func (p *PostgresDatabase) ensureMigrationsTableExits() error {
	if _, err := p.ExecSql(nil, "CREATE TABLE IF NOT EXISTS migrations "+
		"(name VARCHAR(255) PRIMARY KEY,"+
		" timestamp BIGINT NOT NULL);"); err != nil {
		return err
	}

	return nil
}
