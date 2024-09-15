package transaction

import "context"

const (
	RepositoryName            = "TransactionRepository"
	NotFoundError             = "transaction not found"
	DuplicateTransactionError = "duplicated transaction"
	ZeroAmountError           = "amount must be different from zero"
)

type Repository interface {
	//Add handler
	Save(ctx context.Context, transaction Transaction) error
	SaveBatch(ctx context.Context, transactions []Transaction) error
	//Add handler
	FindByID(ctx context.Context, transactionID string) (Transaction, error)
	FindByUserIDWithOptions(ctx context.Context, userID, fromDate, toDate string) ([]Transaction, error)
}
