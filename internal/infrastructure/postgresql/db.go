package postgresql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
	"github.com/sebastianreh/user-balance-api/internal/infrastructure/config"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
)

func NewPostgresDB(cfg config.Config, log logger.Logger) (*sql.DB, error) {
	pgCfg := cfg.Postgres

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		pgCfg.Host, pgCfg.Port, pgCfg.User, pgCfg.Password, pgCfg.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Error("failed to open database", err)
		time.Sleep(time.Duration(pgCfg.ReconnectIdle) * time.Second)
		return NewPostgresDB(cfg, log)
	}

	if err = db.Ping(); err != nil {
		log.Error("failed to ping database", err)
		time.Sleep(time.Duration(pgCfg.ReconnectIdle) * time.Second)
		return NewPostgresDB(cfg, log)
	}

	log.Info("Successfully connected to the PostgreSQL database", pgCfg.DBName)

	return db, nil
}
