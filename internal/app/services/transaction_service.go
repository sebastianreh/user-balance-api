package services

import (
	"context"
	"errors"
	. "github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, transactionEntity Transaction) error
	UpdateTransaction(ctx context.Context, transactionEntity Transaction) error
	GetTransaction(ctx context.Context, transactionID string) (Transaction, error)
	DeleteTransaction(ctx context.Context, transactionID string) error
}

type transactionService struct {
	log        logger.Logger
	repository Repository
}

func NewTransactionService(log logger.Logger, repository Repository) TransactionService {
	return &transactionService{
		log:        log,
		repository: repository,
	}
}

func (t *transactionService) CreateTransaction(ctx context.Context, transactionEntity Transaction) error {
	return t.repository.Save(ctx, transactionEntity)
}

func (t *transactionService) UpdateTransaction(ctx context.Context, transactionEntity Transaction) error {
	_, err := t.repository.FindByID(ctx, transactionEntity.ID)
	if err != nil {
		return err
	}

	return t.repository.Save(ctx, transactionEntity)
}

func (t *transactionService) GetTransaction(ctx context.Context, transactionID string) (Transaction, error) {
	var transactionEntity Transaction
	transactionEntity, err := t.repository.FindByID(ctx, transactionID)
	if transactionEntity.IsDeleted {
		return transactionEntity, errors.New(NotFoundError)
	}
	return transactionEntity, err
}

func (t *transactionService) DeleteTransaction(ctx context.Context, transactionID string) error {
	transactionEntity, err := t.repository.FindByID(ctx, transactionID)
	if err != nil {
		return err
	}

	transactionEntity.IsDeleted = true
	return t.repository.Save(ctx, transactionEntity)
}
