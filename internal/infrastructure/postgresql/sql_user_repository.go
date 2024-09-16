package postgresql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"

	"github.com/sebastianreh/user-balance-api/internal/domain/user"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
)

type sqlUserRepository struct {
	log logger.Logger
	db  *sql.DB
}

func NewSQLUserRepository(log logger.Logger, db *sql.DB) user.Repository {
	return &sqlUserRepository{
		log: log,
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
	err := s.ValidateDeletedUser(ctx, userEntity.ID)
	if err != nil {
		return err
	}

	_, err = s.db.ExecContext(ctx, query, userEntity.ID, userEntity.FirstName, userEntity.LastName, userEntity.Email)
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
	err := row.Scan(&userEntity.ID, &userEntity.FirstName, &userEntity.LastName, &userEntity.Email, &userEntity.IsDeleted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return userEntity, errors.New(user.NotFoundError)
		}

		s.log.ErrorAt(err, user.RepositoryName, "FindByID")
		return userEntity, err
	}

	if userEntity.IsDeleted {
		return userEntity, errors.New(user.NotFoundError)
	}

	return userEntity, nil
}

func (s *sqlUserRepository) Delete(ctx context.Context, userID string) error {
	err := s.ValidateDeletedUser(ctx, userID)
	if err != nil {
		return err
	}

	query := UpdateIsDeletedUser
	_, err = s.db.ExecContext(ctx, query, userID, true)
	if err != nil {
		s.log.ErrorAt(err, transaction.RepositoryName, "Update")
		return err
	}

	return nil
}

func (s *sqlUserRepository) ValidateDeletedUser(ctx context.Context, userID string) error {
	var foundUser user.User
	row := s.db.QueryRowContext(ctx, FindUserByID, userID)
	err := row.Scan(&foundUser.ID, &foundUser.FirstName, &foundUser.LastName, &foundUser.Email, &foundUser.IsDeleted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New(user.NotFoundError)
		}

		s.log.ErrorAt(err, user.RepositoryName, "ValidateDeletedUser")
		return err
	}

	if foundUser.IsDeleted {
		return errors.New(user.NotFoundError)
	}

	return nil
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
	FindUserByID        = "SELECT id, first_name, last_name, email, is_deleted FROM users WHERE id = $1"
	UpdateIsDeletedUser = "UPDATE users SET is_deleted = $2 WHERE id = $1"
)
