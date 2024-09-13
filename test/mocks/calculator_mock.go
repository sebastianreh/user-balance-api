package mocks

import (
	"github.com/sebastianreh/user-balance-api/internal/domain/balance"
	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"github.com/stretchr/testify/mock"
)

type CalculatorMock struct {
	mock.Mock
}

func NewCalculatorMock() *CalculatorMock {
	return new(CalculatorMock)
}

func (m *CalculatorMock) CalculateBalanceByUser(transactions []transaction.Transaction) balance.UserBalance {
	args := m.Called(transactions)
	return args.Get(0).(balance.UserBalance)
}
