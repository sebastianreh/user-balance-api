package sql_repository_test

import (
	"context"
	"github.com/sebastianreh/user-balance-api/internal/domain/user"
	"github.com/sebastianreh/user-balance-api/pkg/strings"
	"testing"
	"time"

	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"github.com/sebastianreh/user-balance-api/internal/infrastructure/postgre_sql"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	"github.com/sebastianreh/user-balance-api/test/integration/sql_repository"
	"github.com/stretchr/testify/assert"
)

const (
	notFoundError = "pq: insert or update on table \"transactions\" violates foreign key constraint \"transactions_user_id_fkey\""
	nilDateError  = "pq: null value in column \"date_time\" of relation \"transactions\" violates not-null constraint"
)

func Test_SqlTransactionRepository_Save(t *testing.T) {
	ctx := context.TODO()
	testDb := sql_repository.SetupTestDB(t)
	testDb.RunMigrations(t)
	log := logger.NewLogger()
	repo := postgre_sql.NewSqlTransactionRepository(log, testDb.Db)
	defer testDb.TeardownTestDB(t)
	userID := testDb.CreateUser(t, user.User{
		FirstName: "user",
		LastName:  "lastname",
		Email:     "user@email.com",
	})
	now := time.Now()

	t.Run("When Save succeeds", func(t *testing.T) {
		defer testDb.CleanTransactions(t)
		tx := transaction.Transaction{
			ID:       "1",
			UserID:   userID,
			Amount:   100.00,
			DateTime: &now,
		}

		err := repo.Save(ctx, tx)
		assert.Nil(t, err)

		savedTransaction, err := repo.FindByID(ctx, "1")
		assert.Nil(t, err)
		assert.Equal(t, tx.ID, savedTransaction.ID)
	})

	t.Run("When Save returns a duplicate error", func(t *testing.T) {
		defer testDb.CleanTransactions(t)
		tx := transaction.Transaction{
			ID:       "1",
			UserID:   userID,
			Amount:   100.00,
			DateTime: &now,
		}

		err := repo.Save(ctx, tx)
		assert.Nil(t, err)

		err = repo.Save(ctx, tx)
		assert.Error(t, err)
		assert.Equal(t, transaction.DuplicateTransactionError, err.Error())
	})

	t.Run("When Save returns a zero amount error", func(t *testing.T) {
		tx := transaction.Transaction{
			ID:     "1",
			UserID: userID,
		}

		err := repo.Save(ctx, tx)
		assert.Error(t, err)
		assert.Equal(t, transaction.ZeroAmountError, err.Error())
	})

	t.Run("When Save returns a date_time nil error", func(t *testing.T) {
		tx := transaction.Transaction{
			ID:     "1",
			UserID: userID,
			Amount: 100.00,
		}

		err := repo.Save(ctx, tx)
		assert.Error(t, err)
		assert.Equal(t, nilDateError, err.Error())
	})

	t.Run("When Save returns a user not found error", func(t *testing.T) {
		tx := transaction.Transaction{
			ID:       "1",
			UserID:   "2000",
			Amount:   100.00,
			DateTime: &now,
		}

		err := repo.Save(ctx, tx)
		assert.Error(t, err)
		assert.Equal(t, notFoundError, err.Error())
	})
}

func Test_SqlTransactionRepository_SaveBatch(t *testing.T) {
	ctx := context.TODO()
	testDb := sql_repository.SetupTestDB(t)
	testDb.RunMigrations(t)
	log := logger.NewLogger()
	repo := postgre_sql.NewSqlTransactionRepository(log, testDb.Db)
	defer testDb.TeardownTestDB(t)
	userID := testDb.CreateUser(t, user.User{
		FirstName: "user",
		LastName:  "lastname",
		Email:     "user@email.com",
	})
	now := time.Now()

	t.Run("When SaveBatch succeeds", func(t *testing.T) {
		defer testDb.CleanTransactions(t)
		transactions := []transaction.Transaction{
			{ID: "1", UserID: userID, Amount: 100.00, DateTime: &now},
			{ID: "2", UserID: userID, Amount: 200.00, DateTime: &now},
		}

		err := repo.SaveBatch(ctx, transactions)
		assert.Nil(t, err)

		savedTransaction, err := repo.FindByID(ctx, "1")
		assert.Nil(t, err)
		assert.Equal(t, "1", savedTransaction.ID)

		savedTransaction, err = repo.FindByID(ctx, "2")
		assert.Nil(t, err)
		assert.Equal(t, "2", savedTransaction.ID)
	})

	t.Run("When SaveBatch returns a duplicate error", func(t *testing.T) {
		transactions := []transaction.Transaction{
			{ID: "1", UserID: userID, Amount: 100.00, DateTime: &now},
			{ID: "1", UserID: userID, Amount: 200.00, DateTime: &now}, // Duplicate ID
		}

		err := repo.SaveBatch(ctx, transactions)
		assert.Error(t, err)
		assert.Equal(t, transaction.DuplicateTransactionError, err.Error())
	})

	t.Run("When SaveBatch returns a zero amount error", func(t *testing.T) {
		transactions := []transaction.Transaction{
			{ID: "1", UserID: "1", DateTime: &now}, // Zero amount
		}

		err := repo.SaveBatch(ctx, transactions)
		assert.Error(t, err)
		assert.Equal(t, transaction.ZeroAmountError, err.Error())
	})

	t.Run("When SaveBatch returns a date_time nil error", func(t *testing.T) {
		transactions := []transaction.Transaction{
			{ID: "1", UserID: userID, Amount: 100.00}, // Missing DateTime
		}

		err := repo.SaveBatch(ctx, transactions)
		assert.Error(t, err)
		assert.Equal(t, nilDateError, err.Error())
	})

	t.Run("When SaveBatch returns a user not found error", func(t *testing.T) {
		transactions := []transaction.Transaction{
			{ID: "1", UserID: "2", Amount: 100.00, DateTime: &now}, // UserID not found
		}

		err := repo.SaveBatch(ctx, transactions)
		assert.Error(t, err)
		assert.Equal(t, notFoundError, err.Error())
	})
}

func Test_SqlTransactionRepository_FindByID(t *testing.T) {
	ctx := context.TODO()
	testDb := sql_repository.SetupTestDB(t)
	testDb.RunMigrations(t)
	log := logger.NewLogger()
	repo := postgre_sql.NewSqlTransactionRepository(log, testDb.Db)
	defer testDb.TeardownTestDB(t)
	userID := testDb.CreateUser(t, user.User{
		FirstName: "user",
		LastName:  "lastname",
		Email:     "user@email.com",
	})
	now := time.Now()

	t.Run("When FindByID succeeds", func(t *testing.T) {
		defer testDb.CleanTransactions(t)
		tx := transaction.Transaction{
			ID:       "1",
			UserID:   userID,
			Amount:   100.00,
			DateTime: &now,
		}

		err := repo.Save(ctx, tx)
		assert.Nil(t, err)

		savedTransaction, err := repo.FindByID(ctx, "1")
		assert.Nil(t, err)
		assert.Equal(t, "1", savedTransaction.ID)
		assert.Equal(t, tx.UserID, savedTransaction.UserID)
		assert.Equal(t, tx.Amount, savedTransaction.Amount)
		assert.Equal(t, tx.DateTime.UTC().Format(time.RFC3339), savedTransaction.DateTime.Format(time.RFC3339))
	})

	t.Run("When FindByID returns no results", func(t *testing.T) {
		defer testDb.CleanTransactions(t)
		savedTransaction, err := repo.FindByID(ctx, "nonexistent")
		assert.Error(t, err)
		assert.Empty(t, savedTransaction, 0)
	})

	t.Run("When FindByID encounters a database error", func(t *testing.T) {
		testDb.CleanTransactions(t)
		err := testDb.Db.Close()
		if err != nil {
			t.Fatalf(err.Error())
		}

		_, err = repo.FindByID(ctx, "1")
		assert.Error(t, err)
		assert.Equal(t, "sql: database is closed", err.Error())
	})
}

func Test_SqlTransactionRepository_FindByUserIDWithOptions(t *testing.T) {
	ctx := context.TODO()
	testDb := sql_repository.SetupTestDB(t)
	testDb.RunMigrations(t)
	log := logger.NewLogger()
	repo := postgre_sql.NewSqlTransactionRepository(log, testDb.Db)
	defer testDb.TeardownTestDB(t)

	userID := testDb.CreateUser(t, user.User{
		ID:        "1",
		FirstName: "Test",
		LastName:  "User",
		Email:     "testuser@email.com",
	})

	now := time.Now()

	t.Run("When FindByUserIDWithOptions succeeds without date range", func(t *testing.T) {
		defer testDb.CleanTransactions(t)
		transaction1 := transaction.Transaction{
			ID:       "1",
			UserID:   userID,
			Amount:   100.00,
			DateTime: &now,
		}

		transaction2 := transaction.Transaction{
			ID:       "2",
			UserID:   userID,
			Amount:   200.00,
			DateTime: &now,
		}

		err := repo.Save(ctx, transaction1)
		assert.Nil(t, err)

		err = repo.Save(ctx, transaction2)
		assert.Nil(t, err)

		transactions, err := repo.FindByUserIDWithOptions(ctx, userID, customStr.Empty, customStr.Empty)
		assert.Nil(t, err)
		assert.Len(t, transactions, 2)
		assert.Equal(t, transaction1.ID, transactions[0].ID)
		assert.Equal(t, transaction2.ID, transactions[1].ID)
	})

	t.Run("When FindByUserIDWithOptions succeeds with date range", func(t *testing.T) {
		defer testDb.CleanTransactions(t)
		pastTime := now.AddDate(0, 0, -1)
		transaction1 := transaction.Transaction{
			ID:       "2",
			UserID:   userID,
			Amount:   100.00,
			DateTime: &pastTime,
		}

		pastTime = pastTime.AddDate(0, -6, 0)
		transaction2 := transaction.Transaction{
			ID:       "1",
			UserID:   userID,
			Amount:   200.00,
			DateTime: &pastTime,
		}

		fromDate := now.AddDate(-1, 0, 0).Format(time.RFC3339)

		err := repo.Save(ctx, transaction1)
		assert.Nil(t, err)

		err = repo.Save(ctx, transaction2)
		assert.Nil(t, err)

		transactions, err := repo.FindByUserIDWithOptions(ctx, userID, fromDate, now.Format(time.RFC3339))
		assert.Nil(t, err)
		assert.Len(t, transactions, 2)
	})

	t.Run("When FindByUserIDWithOptions returns no results", func(t *testing.T) {
		transactions, err := repo.FindByUserIDWithOptions(ctx, "1231412", customStr.Empty, customStr.Empty)
		assert.Nil(t, err)
		assert.Len(t, transactions, 0)
	})

	t.Run("When FindByUserIDWithOptions encounters a database error", func(t *testing.T) {
		err := testDb.Db.Close()
		assert.Nil(t, err)

		_, err = repo.FindByUserIDWithOptions(ctx, "1", customStr.Empty, customStr.Empty)
		assert.Error(t, err)
		assert.Equal(t, "sql: database is closed", err.Error())
	})
}
