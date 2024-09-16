package sqlrepository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"testing"

	"github.com/sebastianreh/user-balance-api/internal/domain/user"

	"github.com/sebastianreh/user-balance-api/internal/infrastructure/config"
	"github.com/sebastianreh/user-balance-api/internal/infrastructure/postgresql"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
)

const (
	testDBName         = "test_db"
	deleteUsers        = "TRUNCATE TABLE users RESTART IDENTITY CASCADE"
	deleteTransactions = "TRUNCATE TABLE transactions RESTART IDENTITY CASCADE"
)

type TestSQLRepository struct {
	DB  *sql.DB
	log logger.Logger
}

func SetupTestDB(t *testing.T) *TestSQLRepository {
	log := logger.NewLogger()
	cfg := config.NewConfig()

	dataSource := fmt.Sprintf("host=%s port=%s user=%s password=%s sslmode=disable",
		cfg.Postgres.Host, cfg.Postgres.Port, cfg.Postgres.User, cfg.Postgres.Password)
	db, err := sql.Open("postgres", dataSource)
	if err != nil {
		t.Fatalf("Failed to connect to PostgreSQL server: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", testDBName))
	if err != nil && !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("Failed to create test database: %v", err)
	}

	cfg.Postgres.DBName = testDBName
	db, err = postgresql.NewPostgresDB(cfg, log)
	if err != nil {
		t.Fatalf("Database initialization error, shutting down server: %v", err)
	}

	return &TestSQLRepository{
		DB:  db,
		log: log,
	}
}

func (r *TestSQLRepository) RunMigrations(t *testing.T) {
	migrations := postgresql.NewSQLMigrations(r.log, r.DB)
	if err := migrations.RunMigrations(); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}
}

func (r *TestSQLRepository) CleanUsers(t *testing.T) {
	r.cleanDatabase(t, deleteUsers)
}

func (r *TestSQLRepository) CleanTransactions(t *testing.T) {
	r.cleanDatabase(t, deleteTransactions)
}

func (r *TestSQLRepository) cleanDatabase(t *testing.T, query string) {
	_, err := r.DB.Exec(query)
	if err != nil {
		t.Fatalf("Failed to delete transactions: %v", err)
	}
}

func (r *TestSQLRepository) CreateUser(t *testing.T, input user.User) string {
	repo := postgresql.NewSQLUserRepository(r.log, r.DB)
	userID, err := repo.Save(context.TODO(), input)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	return userID
}

func (r *TestSQLRepository) TeardownTestDB(t *testing.T) {
	if r.DB != nil {
		err := r.DB.Close()
		if err != nil {
			t.Errorf("Failed to close the test database: %v", err)
		}
	}
}
