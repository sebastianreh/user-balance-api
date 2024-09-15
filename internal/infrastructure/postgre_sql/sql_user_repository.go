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

func (s *sqlUserRepository) Save(ctx context.Context, userEntity user.User) (string, error) {
	query := SaveUser
	var createdID string
	err := s.db.QueryRowContext(ctx, query, userEntity.FirstName, userEntity.LastName, userEntity.Email).Scan(&createdID)
	if err != nil {
		s.log.ErrorAt(err, user.RepositoryName, "Save")
		return "", err
	}

	return createdID, nil
}

func (s *sqlUserRepository) Update(ctx context.Context, userEntity user.User) error {
	query := UpdateUser
	_, err := s.db.ExecContext(ctx, query, userEntity.ID, userEntity.FirstName, userEntity.LastName, userEntity.Email)
	if err != nil {
		s.log.ErrorAt(err, user.RepositoryName, "Update")
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
	INSERT INTO users (first_name, last_name, email) 
	VALUES ($1, $2, $3) 
	RETURNING id;`
	UpdateUser = `
	UPDATE users 
	SET first_name = COALESCE(NULLIF($2, ''), first_name), 
		last_name = COALESCE(NULLIF($3, ''), last_name), 
		email = COALESCE(NULLIF($4, ''), email) 
	WHERE id = $1`
	FindUserByID = "SELECT id, first_name, last_name, email FROM users WHERE id = $1"
)
