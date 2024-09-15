package mocks

import (
	"context"
	"github.com/sebastianreh/user-balance-api/internal/domain/balance"
	"github.com/stretchr/testify/mock"
)

type BalanceServiceMock struct {
	mock.Mock
}

func NewBalanceServiceMock() *BalanceServiceMock {
	return new(BalanceServiceMock)
}

func (m *BalanceServiceMock) GetBalanceByUserIDWithOptions(ctx context.Context, userID string, fromDate, toDate string) (balance.UserBalance, error) {
	args := m.Called(ctx, userID, fromDate, toDate)
	return args.Get(0).(balance.UserBalance), args.Error(1)
}

func (m *BalanceServiceMock) GetBalanceByUserID(ctx context.Context, userID string) (balance.UserBalance, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(balance.UserBalance), args.Error(1)
}
