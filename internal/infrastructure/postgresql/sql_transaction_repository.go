package postgresql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"github.com/sebastianreh/user-balance-api/internal/domain/user"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
)

type sqlTransactionRepository struct {
	log logger.Logger
	db  *sql.DB
}

func NewSQLTransactionRepository(log logger.Logger, db *sql.DB) transaction.Repository {
	return &sqlTransactionRepository{
		log: log,
		db:  db,
	}
}

func (s *sqlTransactionRepository) Save(ctx context.Context, userTransaction transaction.Transaction) error {
	var userFound user.User
	var oldTransaction transaction.Transaction
	if userTransaction.Amount == 0 {
		return errors.New(transaction.ZeroAmountError)
	}

	row := s.db.QueryRowContext(ctx, FindByID, userTransaction.ID)
	err := row.Scan(&oldTransaction.ID, &oldTransaction.UserID,
		&oldTransaction.Amount, &oldTransaction.DateTime, &oldTransaction.IsDeleted)

	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	if oldTransaction.IsDeleted {
		_, err = s.db.ExecContext(ctx, UpdateIsDeletedTransaction, userTransaction.ID, false)
		if err != nil {
			s.log.ErrorAt(err, transaction.RepositoryName, "Update")
			return err
		}
		return nil
	}

	err = s.db.QueryRowContext(ctx, FindUserByID, userTransaction.UserID).Scan(&userFound.ID, &userFound.FirstName,
		&userFound.LastName, &userFound.Email, &userFound.IsDeleted)
	if userFound.IsDeleted {
		return errors.New(user.NotFoundError)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New(user.NotFoundError)
		}

		s.log.ErrorAt(err, user.RepositoryName, "Save")
		return err
	}

	if userFound.IsDeleted {
		return errors.New(user.NotFoundError)
	}

	query := SaveByUserID
	_, err = s.db.ExecContext(ctx, query, userTransaction.ID, userTransaction.UserID,
		userTransaction.Amount, userTransaction.DateTime)
	if err != nil {
		s.log.ErrorAt(err, transaction.RepositoryName, "Save")
		duplicateErr := handleDuplicateError(err)
		if duplicateErr != nil {
			err = duplicateErr
		}
		return err
	}

	return nil
}

func (s *sqlTransactionRepository) Delete(ctx context.Context, transactionID string) error {
	query := UpdateIsDeletedTransaction
	_, err := s.db.ExecContext(ctx, query, transactionID, true)
	if err != nil {
		s.log.ErrorAt(err, transaction.RepositoryName, "Update")
		return err
	}

	return nil
}

func (s *sqlTransactionRepository) SaveBatch(ctx context.Context, transactions []transaction.Transaction) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		s.log.ErrorAt(err, transaction.RepositoryName, "SaveBatch")
		return err
	}

	query := SaveByUserID
	stmt, err := tx.PrepareContext(ctx, query)
	if err != nil {
		s.log.ErrorAt(err, transaction.RepositoryName, "SaveBatch")
		_ = tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, transactionEntity := range transactions {
		if transactionEntity.Amount == 0 {
			_ = tx.Rollback()
			return errors.New(transaction.ZeroAmountError)
		}

		_, err = stmt.ExecContext(ctx, transactionEntity.ID, transactionEntity.UserID,
			transactionEntity.Amount, transactionEntity.DateTime)
		if err != nil {
			s.log.ErrorAt(err, transaction.RepositoryName, "SaveBatch")
			foreignKeyErr := handleForeignKeyError(err)
			if foreignKeyErr != nil {
				err = foreignKeyErr
			}

			duplicateErr := handleDuplicateError(err)
			if duplicateErr != nil {
				err = duplicateErr
			}
			_ = tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		s.log.ErrorAt(err, transaction.RepositoryName, "SaveBatch")
		return err
	}

	return nil
}

func (s *sqlTransactionRepository) Update(ctx context.Context, userTransaction transaction.Transaction) error {
	query := UpdateTransaction
	if userTransaction.Amount == 0 {
		return errors.New(transaction.ZeroAmountError)
	}

	_, err := s.db.ExecContext(ctx, query, userTransaction.ID, userTransaction.UserID,
		userTransaction.Amount, userTransaction.DateTime, userTransaction.IsDeleted)
	if err != nil {
		s.log.ErrorAt(err, transaction.RepositoryName, "Update")
		return err
	}

	return nil
}

func (s *sqlTransactionRepository) FindByID(ctx context.Context, transactionID string) (transaction.Transaction, error) {
	var transactionEntity transaction.Transaction
	query := FindByID
	row := s.db.QueryRowContext(ctx, query, transactionID)
	err := row.Scan(&transactionEntity.ID, &transactionEntity.UserID,
		&transactionEntity.Amount, &transactionEntity.DateTime, &transactionEntity.IsDeleted)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return transactionEntity, errors.New(transaction.NotFoundError)
		}

		s.log.ErrorAt(err, transaction.RepositoryName, "FindByID")
		return transactionEntity, err
	}

	if transactionEntity.IsDeleted {
		return transactionEntity, errors.New(transaction.NotFoundError)
	}

	return transactionEntity, nil
}

func (s *sqlTransactionRepository) FindByUserIDWithOptions(ctx context.Context, userID, fromDate,
	toDate string) ([]transaction.Transaction, error) {
	query := findByUserIDOptionalDateRangeQuery(fromDate, toDate)
	args := []interface{}{userID}

	if fromDate != "" && toDate != "" {
		args = append(args, fromDate, toDate)
	}

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		s.log.ErrorAt(err, transaction.RepositoryName, "FindByUserIDWithOptions")
		return nil, err
	}

	defer rows.Close()

	var transactions []transaction.Transaction
	for rows.Next() {
		var transactionEntity transaction.Transaction
		if err = rows.Scan(&transactionEntity.ID, &transactionEntity.UserID,
			&transactionEntity.Amount, &transactionEntity.DateTime, &transactionEntity.IsDeleted); err != nil {
			s.log.ErrorAt(err, transaction.RepositoryName, "FindByUserIDWithOptions")
			return nil, err
		}
		if !transactionEntity.IsDeleted {
			transactions = append(transactions, transactionEntity)
		}
	}

	return transactions, nil
}

func findByUserIDOptionalDateRangeQuery(fromDate, toDate string) string {
	query := GetAllByUserID

	if fromDate != "" && toDate != "" {
		query += FromToDateOption
	}

	return query
}

func handleDuplicateError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pqErr.Code == "23505" {
			return errors.New(transaction.DuplicateTransactionError)
		}
	}
	return nil
}

func handleForeignKeyError(err error) error {
	var pqErr *pq.Error
	if errors.As(err, &pqErr) {
		if pqErr.Code == "23503" && pqErr.Constraint == "transactions_user_id_fkey" {
			return errors.New(user.NotFoundError)
		}
	}
	return nil
}

const (
	SaveByUserID               = "INSERT INTO transactions (id, user_id, amount, date_time) VALUES ($1, $2, $3, $4)"
	UpdateIsDeletedTransaction = "UPDATE transactions SET is_deleted = $2 WHERE id = $1"
	UpdateTransaction          = "UPDATE transactions SET user_id = $2, amount = $3, date_time = $4 WHERE id = $1"
	GetAllByUserID             = "SELECT * FROM transactions WHERE user_id = $1"
	FindByID                   = "SELECT * FROM transactions WHERE id = $1"
	FromToDateOption           = ` AND date_time >= CAST($2 AS timestamptz) AND date_time <= CAST($3 AS timestamptz)`
)
