package sqlrepository_test

import (
	"github.com/sebastianreh/user-balance-api/internal/infrastructure/postgresql"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	"github.com/sebastianreh/user-balance-api/test/integration/sqlrepository"
	"github.com/stretchr/testify/assert"

	"testing"
)

func Test_RunMigrations(t *testing.T) {
	repo := sqlrepository.SetupTestDB(t)
	defer repo.TeardownTestDB(t)

	t.Run("When RunMigrations executes successfully", func(t *testing.T) {
		migrations := postgresql.NewSQLMigrations(logger.NewLogger(), repo.DB)

		err := migrations.RunMigrations()

		assert.Nil(t, err)

		_, err = repo.DB.Exec("SELECT 1 FROM users LIMIT 1;")
		assert.Nil(t, err, "users table should exist")

		_, err = repo.DB.Exec("SELECT 1 FROM transactions LIMIT 1;")
		assert.Nil(t, err, "transactions table should exist")

		_, err = repo.DB.Exec("SELECT indexname FROM pg_indexes WHERE indexname = 'idx_transactions_user_id';")
		assert.Nil(t, err)

		_, err = repo.DB.Exec("SELECT indexname FROM pg_indexes WHERE indexname = 'idx_transactions_date_time';")
		assert.Nil(t, err)
	})
}
