package transaction

import "context"

const (
	RepositoryName = "TransactionRepository"
)

type Repository interface {
	Save(ctx context.Context, transaction Transaction) error
	SaveBatch(ctx context.Context, transactions []Transaction) error
	FindByID(ctx context.Context, transactionID string) ([]Transaction, error)
	FindByUserIDWithOptions(ctx context.Context, userID, fromDate, toDate string) ([]Transaction, error)
}
