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
	//Add handler
	Save(ctx context.Context, user User) error
	FindByID(ctx context.Context, userID string) (User, error)
}
