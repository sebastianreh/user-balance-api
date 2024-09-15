package services_test

import (
	"context"
	"errors"
	"github.com/sebastianreh/user-balance-api/internal/app/services"
	"github.com/sebastianreh/user-balance-api/internal/domain/user"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	"github.com/sebastianreh/user-balance-api/test/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserService_CreateUser(t *testing.T) {
	ctx := context.TODO()
	log := logger.NewLogger()

	t.Run("When CreateUser succeeds", func(t *testing.T) {
		mockRepo := mocks.NewUserRepositoryMock()
		service := services.NewUserService(log, mockRepo)

		userEntity := user.User{ID: "1", FirstName: "user", LastName: "lastname", Email: "user@email.com"}

		mockRepo.On("Save", ctx, userEntity).Return(userEntity.ID, nil)

		userID, err := service.CreateUser(ctx, userEntity)
		assert.Nil(t, err)
		assert.Equal(t, userEntity.ID, userID)
		mockRepo.AssertCalled(t, "Save", ctx, userEntity)
	})

	t.Run("When CreateUser fails", func(t *testing.T) {
		mockRepo := mocks.NewUserRepositoryMock()
		service := services.NewUserService(log, mockRepo)

		userEntity := user.User{ID: "1", FirstName: "user", LastName: "lastname", Email: "user@email.com"}
		expectedError := errors.New("repository error")

		mockRepo.On("Save", ctx, userEntity).Return("", expectedError)

		id, err := service.CreateUser(ctx, userEntity)
		assert.Equal(t, expectedError, err)
		assert.Empty(t, id)
		mockRepo.AssertCalled(t, "Save", ctx, userEntity)
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	ctx := context.TODO()
	log := logger.NewLogger()

	t.Run("When UpdateUser succeeds", func(t *testing.T) {
		mockRepo := mocks.NewUserRepositoryMock()
		service := services.NewUserService(log, mockRepo)

		userEntity := user.User{ID: "1", FirstName: "user", LastName: "lastname", Email: "user@email.com"}

		mockRepo.On("FindByID", ctx, "1").Return(userEntity, nil)
		mockRepo.On("Update", ctx, userEntity).Return(nil)

		err := service.UpdateUser(ctx, userEntity)
		assert.Nil(t, err)
		mockRepo.AssertCalled(t, "FindByID", ctx, "1")
		mockRepo.AssertCalled(t, "Update", ctx, userEntity)
	})

	t.Run("When FindByID fails in UpdateUser", func(t *testing.T) {
		mockRepo := mocks.NewUserRepositoryMock()
		service := services.NewUserService(log, mockRepo)

		userEntity := user.User{ID: "1", FirstName: "user", LastName: "lastname", Email: "user@email.com"}
		expectedError := errors.New("user not found")

		mockRepo.On("FindByID", ctx, "1").Return(user.User{}, expectedError)

		err := service.UpdateUser(ctx, userEntity)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "FindByID", ctx, "1")
	})

	t.Run("When Save fails in UpdateUser", func(t *testing.T) {
		mockRepo := mocks.NewUserRepositoryMock()
		service := services.NewUserService(log, mockRepo)

		userEntity := user.User{ID: "1", FirstName: "user", LastName: "lastname", Email: "user@email.com"}
		expectedError := errors.New("repository error")

		mockRepo.On("FindByID", ctx, "1").Return(userEntity, nil)
		mockRepo.On("Update", ctx, userEntity).Return(expectedError)

		err := service.UpdateUser(ctx, userEntity)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "FindByID", ctx, "1")
		mockRepo.AssertCalled(t, "Update", ctx, userEntity)
	})
}

func TestUserService_GetUser(t *testing.T) {
	ctx := context.TODO()
	log := logger.NewLogger()

	t.Run("When GetUser succeeds", func(t *testing.T) {
		mockRepo := mocks.NewUserRepositoryMock()
		service := services.NewUserService(log, mockRepo)

		userEntity := user.User{ID: "1", FirstName: "user", LastName: "lastname", Email: "user@email.com"}

		mockRepo.On("FindByID", ctx, "1").Return(userEntity, nil)

		returnedUser, err := service.GetUser(ctx, "1")
		assert.Nil(t, err)
		assert.Equal(t, userEntity, returnedUser)
		mockRepo.AssertCalled(t, "FindByID", ctx, "1")
	})

	t.Run("When GetUser fails with not found error", func(t *testing.T) {
		mockRepo := mocks.NewUserRepositoryMock()
		service := services.NewUserService(log, mockRepo)

		expectedError := errors.New(user.NotFoundError)

		mockRepo.On("FindByID", ctx, "1").Return(user.User{}, expectedError)

		returnedUser, err := service.GetUser(ctx, "1")
		assert.Equal(t, expectedError, err)
		assert.Equal(t, user.User{}, returnedUser)
		mockRepo.AssertCalled(t, "FindByID", ctx, "1")
	})

	t.Run("When GetUser fails with not found error because of logic deletion", func(t *testing.T) {
		mockRepo := mocks.NewUserRepositoryMock()
		service := services.NewUserService(log, mockRepo)

		expectedError := errors.New(user.NotFoundError)

		mockRepo.On("FindByID", ctx, "1").Return(user.User{IsDeleted: true}, expectedError)

		returnedUser, err := service.GetUser(ctx, "1")
		assert.Equal(t, expectedError, err)
		assert.Equal(t, user.User{IsDeleted: true}, returnedUser)
		mockRepo.AssertCalled(t, "FindByID", ctx, "1")
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	ctx := context.TODO()
	log := logger.NewLogger()

	t.Run("When DeleteUser succeeds", func(t *testing.T) {
		mockRepo := mocks.NewUserRepositoryMock()
		service := services.NewUserService(log, mockRepo)

		userEntity := user.User{ID: "1", FirstName: "user", LastName: "lastname", Email: "user@email.com"}
		deletedUserEntity := userEntity
		deletedUserEntity.IsDeleted = true

		mockRepo.On("FindByID", ctx, "1").Return(userEntity, nil)
		mockRepo.On("Update", ctx, deletedUserEntity).Return(nil)

		err := service.DeleteUser(ctx, "1")
		assert.Nil(t, err)
		mockRepo.AssertCalled(t, "FindByID", ctx, "1")
		mockRepo.AssertCalled(t, "Update", ctx, deletedUserEntity)
	})

	t.Run("When FindByID fails in DeleteUser", func(t *testing.T) {
		mockRepo := mocks.NewUserRepositoryMock()
		service := services.NewUserService(log, mockRepo)

		expectedError := errors.New("user not found")

		mockRepo.On("FindByID", ctx, "1").Return(user.User{}, expectedError)

		err := service.DeleteUser(ctx, "1")
		assert.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "FindByID", ctx, "1")
	})

	t.Run("When Save fails in DeleteUser", func(t *testing.T) {
		mockRepo := mocks.NewUserRepositoryMock()
		service := services.NewUserService(log, mockRepo)

		userEntity := user.User{ID: "1", FirstName: "user", LastName: "lastname", Email: "user@email.com"}
		deletedUserEntity := userEntity
		deletedUserEntity.IsDeleted = true
		expectedError := errors.New("repository error")

		mockRepo.On("FindByID", ctx, "1").Return(userEntity, nil)
		mockRepo.On("Update", ctx, deletedUserEntity).Return(expectedError)

		err := service.DeleteUser(ctx, "1")
		assert.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "FindByID", ctx, "1")
		mockRepo.AssertCalled(t, "Update", ctx, deletedUserEntity)
	})
}
