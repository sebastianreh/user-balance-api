package postgre_sql

import (
	"context"
	"database/sql"
	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
)

type sqlTransactionRepository struct {
	log logger.Logger
	db  *sql.DB
}

func NewSqlTransactionRepository(logger logger.Logger, db *sql.DB) transaction.Repository {
	return &sqlTransactionRepository{
		log: logger,
		db:  db,
	}
}

func (s *sqlTransactionRepository) Save(ctx context.Context, userTransaction transaction.Transaction) error {
	query := SaveByUserID
	_, err := s.db.ExecContext(ctx, query, userTransaction.ID, userTransaction.UserID, userTransaction.Amount, userTransaction.DateTime)
	if err != nil {
		s.log.ErrorAt(err, transaction.RepositoryName, "Save")
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
		if _, err = stmt.ExecContext(ctx, transactionEntity.ID, transactionEntity.UserID,
			transactionEntity.Amount, transactionEntity.DateTime); err != nil {
			s.log.ErrorAt(err, transaction.RepositoryName, "SaveBatch")
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

func (s *sqlTransactionRepository) FindByID(ctx context.Context, transactionID string) ([]transaction.Transaction, error) {
	query := FindByID
	rows, err := s.db.QueryContext(ctx, query, transactionID)
	if err != nil {
		s.log.ErrorAt(err, transaction.RepositoryName, "FindByID")
		return nil, err
	}
	defer rows.Close()

	var transactions []transaction.Transaction
	for rows.Next() {
		var transactionEntity transaction.Transaction
		if err = rows.Scan(&transactionEntity.ID, &transactionEntity.UserID,
			&transactionEntity.Amount, &transactionEntity.DateTime); err != nil {
			s.log.ErrorAt(err, transaction.RepositoryName, "FindByID")
			return nil, err
		}
		transactions = append(transactions, transactionEntity)
	}

	return transactions, nil
}

func (s *sqlTransactionRepository) FindByUserIDWithOptions(ctx context.Context, userID, fromDate, toDate string) ([]transaction.Transaction, error) {
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
			&transactionEntity.Amount, &transactionEntity.DateTime); err != nil {
			s.log.ErrorAt(err, transaction.RepositoryName, "FindByUserIDWithOptions")
			return nil, err
		}
		transactions = append(transactions, transactionEntity)
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

const (
	SaveByUserID     = "INSERT INTO transactions (id, user_id, amount, date_time) VALUES ($1, $2, $3, $4)"
	GetAllByUserID   = "SELECT * FROM transactions WHERE user_id = $1"
	FindByID         = "SELECT * FROM transactions WHERE id = $1"
	FromToDateOption = ` AND date_time >= CAST($2 AS timestamptz) AND date_time <= CAST($3 AS timestamptz)`
)
