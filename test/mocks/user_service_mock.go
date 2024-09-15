package mocks

import (
	"context"

	"github.com/sebastianreh/user-balance-api/internal/domain/user"
	"github.com/stretchr/testify/mock"
)

type UserServiceMock struct {
	mock.Mock
}

func NewUserServiceMock() *UserServiceMock {
	return new(UserServiceMock)
}

func (m *UserServiceMock) CreateUser(ctx context.Context, userEntity user.User) error {
	args := m.Called(ctx, userEntity)
	return args.Error(0)
}

func (m *UserServiceMock) UpdateUser(ctx context.Context, userEntity user.User) error {
	args := m.Called(ctx, userEntity)
	return args.Error(0)
}

func (m *UserServiceMock) GetUser(ctx context.Context, userID string) (user.User, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).(user.User), args.Error(1)
}

func (m *UserServiceMock) DeleteUser(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}
