package services

import (
	"context"
	"errors"
	. "github.com/sebastianreh/user-balance-api/internal/domain/user"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
)

type UserService interface {
	CreateUser(ctx context.Context, userEntity User) error
	UpdateUser(ctx context.Context, userEntity User) error
	GetUser(ctx context.Context, userID string) (User, error)
	DeleteUser(ctx context.Context, userID string) error
}

type userService struct {
	log        logger.Logger
	repository Repository
}

func NewUserService(log logger.Logger, repository Repository) UserService {
	return &userService{
		log:        log,
		repository: repository,
	}
}

func (u *userService) CreateUser(ctx context.Context, userEntity User) error {
	return u.repository.Save(ctx, userEntity)
}

func (u *userService) UpdateUser(ctx context.Context, userEntity User) error {
	_, err := u.repository.FindByID(ctx, userEntity.ID)
	if err != nil {
		return err
	}

	return u.repository.Save(ctx, userEntity)
}

func (u *userService) GetUser(ctx context.Context, userID string) (User, error) {
	var userEntity User
	userEntity, err := u.repository.FindByID(ctx, userID)
	if userEntity.IsDeleted {
		return userEntity, errors.New(NotFoundError)
	}
	return userEntity, err
}

func (u *userService) DeleteUser(ctx context.Context, userID string) error {
	userEntity, err := u.repository.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	userEntity.IsDeleted = true
	return u.repository.Save(ctx, userEntity)
}
