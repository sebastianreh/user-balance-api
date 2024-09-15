package user

import (
	"context"
)

const (
	RepositoryName = "UserRepository"
	NotFoundError  = "user not found"
	DuplicateError = "duplicate user"
)

type Repository interface {
	Save(ctx context.Context, user User) (string, error)
	Update(ctx context.Context, user User) error
	FindByID(ctx context.Context, userID string) (User, error)
}
