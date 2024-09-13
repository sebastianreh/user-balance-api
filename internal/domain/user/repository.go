package user

import (
	"context"
)

const (
	RepositoryName = "UserRepository"
)

type Repository interface {
	Save(ctx context.Context, user User) error
	FindByID(ctx context.Context, userID string) (User, error)
	FindByTransactionID(ctx context.Context, transactionID string) (User, error)
}
