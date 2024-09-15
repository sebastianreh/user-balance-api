package postgre_sql

import (
	"context"
	"database/sql"
	"errors"
	"github.com/sebastianreh/user-balance-api/internal/domain/user"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
)

type sqlUserRepository struct {
	log logger.Logger
	db  *sql.DB
}

func NewSqlUserRepository(logger logger.Logger, db *sql.DB) user.Repository {
	return &sqlUserRepository{
		log: logger,
		db:  db,
	}
}

func (s *sqlUserRepository) Save(ctx context.Context, userEntity user.User) error {
	query := SaveUser
	_, err := s.db.ExecContext(ctx, query, userEntity.ID, userEntity.FirstName, userEntity.LastName, userEntity.Email)
	if err != nil {
		s.log.ErrorAt(err, user.RepositoryName, "Save")
		return err
	}

	return nil
}

func (s *sqlUserRepository) FindByID(ctx context.Context, userID string) (user.User, error) {
	var userEntity user.User
	query := FindUserByID
	row := s.db.QueryRowContext(ctx, query, userID)
	err := row.Scan(&userEntity.ID, &userEntity.FirstName, &userEntity.LastName, &userEntity.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userEntity, errors.New(user.NotFoundError)
		}

		s.log.ErrorAt(err, user.RepositoryName, "FindByID")
		return userEntity, err
	}

	return userEntity, nil
}

const (
	SaveUser = `
	INSERT INTO users (id, first_name, last_name, email) 
	VALUES ($1, $2, $3, $4) 
	ON CONFLICT (id) DO UPDATE 
	SET first_name = COALESCE(NULLIF(EXCLUDED.first_name, ''), users.first_name), 
		last_name = COALESCE(NULLIF(EXCLUDED.last_name, ''), users.last_name), 
		email = COALESCE(NULLIF(EXCLUDED.email, ''), users.email)`
	FindUserByID = "SELECT id, first_name, last_name, email FROM users WHERE id = $1"
)
