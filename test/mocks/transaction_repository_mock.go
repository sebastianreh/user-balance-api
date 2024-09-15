package mocks

import (
	"context"
	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"github.com/stretchr/testify/mock"
)

type TransactionRepositoryMock struct {
	mock.Mock
}

func NewTransactionRepositoryMock() *TransactionRepositoryMock {
	return new(TransactionRepositoryMock)
}

func (m *TransactionRepositoryMock) Save(ctx context.Context, tx transaction.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *TransactionRepositoryMock) SaveBatch(ctx context.Context, txs []transaction.Transaction) error {
	args := m.Called(ctx, txs)
	return args.Error(0)
}

func (m *TransactionRepositoryMock) FindByID(ctx context.Context, transactionID string) (transaction.Transaction, error) {
	args := m.Called(ctx, transactionID)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *TransactionRepositoryMock) FindByUserIDWithOptions(ctx context.Context, userID, fromDate, toDate string) ([]transaction.Transaction, error) {
	args := m.Called(ctx, userID, fromDate, toDate)
	return args.Get(0).([]transaction.Transaction), args.Error(1)
}
