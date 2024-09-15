package mocks

import (
	"context"
	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"github.com/stretchr/testify/mock"
)

type TransactionServiceMock struct {
	mock.Mock
}

func NewTransactionServiceMock() *TransactionServiceMock {
	return new(TransactionServiceMock)
}

func (m *TransactionServiceMock) CreateTransaction(ctx context.Context, transactionEntity transaction.Transaction) error {
	args := m.Called(ctx, transactionEntity)
	return args.Error(0)
}

func (m *TransactionServiceMock) UpdateTransaction(ctx context.Context, transactionEntity transaction.Transaction) error {
	args := m.Called(ctx, transactionEntity)
	return args.Error(0)
}

func (m *TransactionServiceMock) GetTransaction(ctx context.Context, transactionID string) (transaction.Transaction, error) {
	args := m.Called(ctx, transactionID)
	return args.Get(0).(transaction.Transaction), args.Error(1)
}

func (m *TransactionServiceMock) DeleteTransaction(ctx context.Context, transactionID string) error {
	args := m.Called(ctx, transactionID)
	return args.Error(0)
}
