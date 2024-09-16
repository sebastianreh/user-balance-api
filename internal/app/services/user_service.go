package services

import (
	"context"

	"github.com/sebastianreh/user-balance-api/internal/domain/user"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
)

type UserService interface {
	CreateUser(ctx context.Context, userEntity user.User) (string, error)
	UpdateUser(ctx context.Context, userEntity user.User) error
	GetUser(ctx context.Context, userID string) (user.User, error)
	DeleteUser(ctx context.Context, userID string) error
}

type userService struct {
	log        logger.Logger
	repository user.Repository
}

func NewUserService(log logger.Logger, repository user.Repository) UserService {
	return &userService{
		log:        log,
		repository: repository,
	}
}

func (u *userService) CreateUser(ctx context.Context, userEntity user.User) (string, error) {
	return u.repository.Save(ctx, userEntity)
}

func (u *userService) UpdateUser(ctx context.Context, userEntity user.User) error {
	_, err := u.repository.FindByID(ctx, userEntity.ID)
	if err != nil {
		return err
	}

	err = u.repository.Update(ctx, userEntity)
	return err
}

func (u *userService) GetUser(ctx context.Context, userID string) (user.User, error) {
	var userEntity user.User
	userEntity, err := u.repository.FindByID(ctx, userID)
	return userEntity, err
}

func (u *userService) DeleteUser(ctx context.Context, userID string) error {
	err := u.repository.Delete(ctx, userID)
	return err
}
