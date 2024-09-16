package services

import (
	"context"

	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, transactionEntity transaction.Transaction) error
	UpdateTransaction(ctx context.Context, transactionEntity transaction.Transaction) error
	GetTransaction(ctx context.Context, transactionID string) (transaction.Transaction, error)
	DeleteTransaction(ctx context.Context, transactionID string) error
}

type transactionService struct {
	log        logger.Logger
	repository transaction.Repository
}

func NewTransactionService(log logger.Logger, repository transaction.Repository) TransactionService {
	return &transactionService{
		log:        log,
		repository: repository,
	}
}

func (t *transactionService) CreateTransaction(ctx context.Context, transactionEntity transaction.Transaction) error {
	return t.repository.Save(ctx, transactionEntity)
}

func (t *transactionService) UpdateTransaction(ctx context.Context, transactionEntity transaction.Transaction) error {
	_, err := t.repository.FindByID(ctx, transactionEntity.ID)
	if err != nil {
		return err
	}

	return t.repository.Update(ctx, transactionEntity)
}

func (t *transactionService) GetTransaction(ctx context.Context, transactionID string) (transaction.Transaction, error) {
	var transactionEntity transaction.Transaction
	transactionEntity, err := t.repository.FindByID(ctx, transactionID)
	return transactionEntity, err
}

func (t *transactionService) DeleteTransaction(ctx context.Context, transactionID string) error {
	return t.repository.Delete(ctx, transactionID)
}
