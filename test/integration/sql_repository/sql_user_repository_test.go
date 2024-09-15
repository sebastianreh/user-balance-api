package sql_repository_test

import (
	"context"
	"testing"

	"github.com/sebastianreh/user-balance-api/internal/domain/user"
	"github.com/sebastianreh/user-balance-api/internal/infrastructure/postgre_sql"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	"github.com/sebastianreh/user-balance-api/test/integration/sql_repository"
	"github.com/stretchr/testify/assert"
)

func Test_SqlUSerRepository_Save(t *testing.T) {
	ctx := context.TODO()
	testDb := sql_repository.SetupTestDB(t)
	testDb.RunMigrations(t)
	log := logger.NewLogger()
	repo := postgre_sql.NewSqlUserRepository(log, testDb.Db)
	defer testDb.TeardownTestDB(t)

	t.Run("When Save succeeds", func(t *testing.T) {
		userEntity := user.User{
			FirstName: "user",
			LastName:  "lastname",
			Email:     "user@email.com",
		}

		id, err := repo.Save(ctx, userEntity)

		assert.Nil(t, err)
		assert.NotEmpty(t, id)
	})

	t.Run("When Save returns an error", func(t *testing.T) {
		userEntity := user.User{
			ID:        "1",
			FirstName: "user",
			LastName:  "lastname",
			Email:     "user@email.com",
		}

		err := testDb.Db.Close()
		if err != nil {
			t.Fatalf(err.Error())
		}

		id, err := repo.Save(ctx, userEntity)

		assert.Error(t, err)
		assert.Empty(t, id)
		assert.Equal(t, "sql: database is closed", err.Error())
	})
}

func Test_SqlUSerRepository_Update(t *testing.T) {
	ctx := context.TODO()
	testDb := sql_repository.SetupTestDB(t)
	testDb.RunMigrations(t)
	log := logger.NewLogger()
	repo := postgre_sql.NewSqlUserRepository(log, testDb.Db)
	defer testDb.TeardownTestDB(t)

	t.Run("When Update succeeds", func(t *testing.T) {
		defer testDb.CleanUsers(t)
		userEntity := user.User{
			FirstName: "user",
			LastName:  "lastname",
			Email:     "user@email.com",
		}
		testDb.CreateUser(t, userEntity)

		userEntity.ID = "1"
		err := repo.Update(ctx, userEntity)

		assert.Nil(t, err)
	})

	t.Run("When Update returns an error", func(t *testing.T) {
		userEntity := user.User{
			ID:        "1",
			FirstName: "user",
			LastName:  "lastname",
			Email:     "user@email.com",
		}

		err := testDb.Db.Close()
		if err != nil {
			t.Fatalf(err.Error())
		}

		id, err := repo.Save(ctx, userEntity)

		assert.Error(t, err)
		assert.Empty(t, id)
		assert.Equal(t, "sql: database is closed", err.Error())
	})
}

func Test_SqlUserRepository_FindByID(t *testing.T) {
	ctx := context.TODO()
	testDb := sql_repository.SetupTestDB(t)
	testDb.RunMigrations(t)
	log := logger.NewLogger()
	repo := postgre_sql.NewSqlUserRepository(log, testDb.Db)
	defer testDb.TeardownTestDB(t)

	t.Run("When FindByUserID succeeds", func(t *testing.T) {
		defer testDb.CleanUsers(t)
		userEntity := user.User{
			ID:        "1",
			FirstName: "user",
			LastName:  "lastname",
			Email:     "user@email.com",
		}

		userID := testDb.CreateUser(t, userEntity)

		foundUser, err := repo.FindByID(ctx, userID)
		assert.Nil(t, err)
		assert.Equal(t, userID, foundUser.ID)
		assert.Equal(t, "user", foundUser.FirstName)
		assert.Equal(t, "lastname", foundUser.LastName)
		assert.Equal(t, "user@email.com", foundUser.Email)
	})

	t.Run("When FindByUserIDWithOptions fails due to DB closure", func(t *testing.T) {
		err := testDb.Db.Close()
		if err != nil {
			t.Fatalf(err.Error())
		}

		_, err = repo.FindByID(ctx, "1")
		assert.Error(t, err)
		assert.Equal(t, "sql: database is closed", err.Error())
	})
}
