package postgre_sql

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/sebastianreh/user-balance-api/internal/infrastructure/config"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
)

func NewPostgresDB(config config.Config, log logger.Logger) (*sql.DB, error) {
	cfg := config.Postgres

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error("failed to open database: %w", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		log.Error("failed to ping database: %w", err)
		return nil, err
	}

	log.Info("Successfully connected to the PostgreSQL database")

	return db, nil
}
