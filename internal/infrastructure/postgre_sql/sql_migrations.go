package postgre_sql

import (
	"database/sql"
	"fmt"
	"github.com/labstack/gommon/log"

	"github.com/sebastianreh/user-balance-api/pkg/logger"
)

const (
	RunMigrationsName = "RunMigrations"
)

type SqlMigrations interface {
	RunMigrations() error
}

type sqlMigrations struct {
	log logger.Logger
	db  *sql.DB
}

func NewSqlMigrations(log logger.Logger, db *sql.DB) SqlMigrations {
	return &sqlMigrations{
		log: log,
		db:  db,
	}
}

func (s *sqlMigrations) RunMigrations() error {
	if _, err := s.db.Exec(createUsersTable); err != nil {
		s.log.ErrorAt(fmt.Errorf("failed to create users table: %w", err),
			RunMigrationsName, "createUsersTable")
		return err
	}

	if _, err := s.db.Exec(createTransactionsTable); err != nil {
		s.log.ErrorAt(fmt.Errorf("failed to create transactions table: %w", err),
			RunMigrationsName, "createTransactionsTable")
		return err
	}

	if _, err := s.db.Exec(createUserIDIndex); err != nil {
		s.log.ErrorAt(fmt.Errorf("failed to create user_id index: %w", err),
			RunMigrationsName, "createUserIDIndex")
		return err
	}

	if _, err := s.db.Exec(createDateTimeIndex); err != nil {
		s.log.ErrorAt(fmt.Errorf("failed to create date_time index: %w", err),
			RunMigrationsName, "createDateTimeIndex")
		return err
	}

	log.Info("Database migrations executed successfully.")
	return nil
}

const (
	createUsersTable = `
	CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL PRIMARY KEY,
	first_name VARCHAR(255),
	last_name VARCHAR(255),
	email VARCHAR(255)
	);`

	createTransactionsTable = `
	CREATE TABLE IF NOT EXISTS transactions (
	id VARCHAR(255) PRIMARY KEY,
	user_id BIGINT REFERENCES users(id) ON DELETE CASCADE,
	amount DECIMAL(10, 2) NOT NULL,
	date_time TIMESTAMPTZ NOT NULL
	);`

	createUserIDIndex = `
	CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);`

	createDateTimeIndex = `
	CREATE INDEX IF NOT EXISTS idx_transactions_date_time ON transactions(date_time);`
)
