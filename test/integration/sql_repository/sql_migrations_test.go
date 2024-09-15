package sql_repository_test

import (
	"github.com/sebastianreh/user-balance-api/internal/infrastructure/postgre_sql"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	"github.com/sebastianreh/user-balance-api/test/integration/sql_repository"
	"github.com/stretchr/testify/assert"

	"testing"
)

func Test_RunMigrations(t *testing.T) {
	repo := sql_repository.SetupTestDB(t)
	defer repo.TeardownTestDB(t)

	t.Run("When RunMigrations executes successfully", func(t *testing.T) {
		migrations := postgre_sql.NewSqlMigrations(logger.NewLogger(), repo.Db)

		err := migrations.RunMigrations()

		assert.Nil(t, err)

		_, err = repo.Db.Exec("SELECT 1 FROM users LIMIT 1;")
		assert.Nil(t, err, "users table should exist")

		_, err = repo.Db.Exec("SELECT 1 FROM transactions LIMIT 1;")
		assert.Nil(t, err, "transactions table should exist")

		_, err = repo.Db.Exec("SELECT indexname FROM pg_indexes WHERE indexname = 'idx_transactions_user_id';")
		assert.Nil(t, err)

		_, err = repo.Db.Exec("SELECT indexname FROM pg_indexes WHERE indexname = 'idx_transactions_date_time';")
		assert.Nil(t, err)
	})
}
