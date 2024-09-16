package mocks

import (
	"context"

	"github.com/sebastianreh/user-balance-api/internal/domain/user"
	"github.com/stretchr/testify/mock"
)

type UserRepositoryMock struct {
	mock.Mock
}

func NewUserRepositoryMock() *UserRepositoryMock {
	return new(UserRepositoryMock)
}

func (m *UserRepositoryMock) Save(ctx context.Context, userEntity user.User) (string, error) {
	args := m.Called(ctx, userEntity)
	return args.Get(0).(string), args.Error(1)
}

func (m *UserRepositoryMock) Update(ctx context.Context, userEntity user.User) error {
	args := m.Called(ctx, userEntity)
	return args.Error(0)
}

func (m *UserRepositoryMock) Delete(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *UserRepositoryMock) FindByID(ctx context.Context, userID string) (user.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *UserRepositoryMock) FindByTransactionID(ctx context.Context, transactionID string) (user.User, error) {
	args := m.Called(ctx, transactionID)
	return args.Get(0).(user.User), args.Error(1)
}
