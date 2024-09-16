package services_test

import (
	"context"
	"errors"
	"testing"

	"github.com/sebastianreh/user-balance-api/test/mocks"

	"github.com/sebastianreh/user-balance-api/internal/app/services"
	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	"github.com/stretchr/testify/assert"
)

func TestTransactionService_CreateTransaction(t *testing.T) {
	ctx := context.TODO()
	log := logger.NewLogger()

	t.Run("When CreateTransaction succeeds", func(t *testing.T) {
		mockRepo := mocks.NewTransactionRepositoryMock()
		service := services.NewTransactionService(log, mockRepo)

		transactionEntity := transaction.Transaction{ID: "1", UserID: "1", Amount: 100}

		mockRepo.On("Save", ctx, transactionEntity).Return(nil)

		err := service.CreateTransaction(ctx, transactionEntity)
		assert.Nil(t, err)
		mockRepo.AssertCalled(t, "Save", ctx, transactionEntity)
	})

	t.Run("When CreateTransaction fails", func(t *testing.T) {
		mockRepo := mocks.NewTransactionRepositoryMock()
		service := services.NewTransactionService(log, mockRepo)

		transactionEntity := transaction.Transaction{ID: "1", UserID: "1", Amount: 100}
		expectedError := errors.New("repository error")

		mockRepo.On("Save", ctx, transactionEntity).Return(expectedError)

		err := service.CreateTransaction(ctx, transactionEntity)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "Save", ctx, transactionEntity)
	})
}

func TestTransactionService_UpdateTransaction(t *testing.T) {
	ctx := context.TODO()
	log := logger.NewLogger()

	t.Run("When UpdateTransaction succeeds", func(t *testing.T) {
		mockRepo := mocks.NewTransactionRepositoryMock()
		service := services.NewTransactionService(log, mockRepo)

		transactionEntity := transaction.Transaction{ID: "1", UserID: "1", Amount: 100}

		mockRepo.On("FindByID", ctx, "1").Return(transactionEntity, nil)
		mockRepo.On("Update", ctx, transactionEntity).Return(nil)

		err := service.UpdateTransaction(ctx, transactionEntity)
		assert.Nil(t, err)
		mockRepo.AssertCalled(t, "FindByID", ctx, "1")
		mockRepo.AssertCalled(t, "Update", ctx, transactionEntity)
	})

	t.Run("When FindByID fails in UpdateTransaction", func(t *testing.T) {
		mockRepo := mocks.NewTransactionRepositoryMock()
		service := services.NewTransactionService(log, mockRepo)

		transactionEntity := transaction.Transaction{ID: "1", UserID: "1", Amount: 100}
		expectedError := errors.New("transaction not found")

		mockRepo.On("FindByID", ctx, "1").Return(transaction.Transaction{}, expectedError)

		err := service.UpdateTransaction(ctx, transactionEntity)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "FindByID", ctx, "1")
	})

	t.Run("When Update fails in UpdateTransaction", func(t *testing.T) {
		mockRepo := mocks.NewTransactionRepositoryMock()
		service := services.NewTransactionService(log, mockRepo)

		transactionEntity := transaction.Transaction{ID: "1", UserID: "1", Amount: 100}
		expectedError := errors.New("repository error")

		mockRepo.On("FindByID", ctx, "1").Return(transactionEntity, nil)
		mockRepo.On("Update", ctx, transactionEntity).Return(expectedError)

		err := service.UpdateTransaction(ctx, transactionEntity)
		assert.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "FindByID", ctx, "1")
		mockRepo.AssertCalled(t, "Update", ctx, transactionEntity)
	})
}

func TestTransactionService_GetTransaction(t *testing.T) {
	ctx := context.TODO()
	log := logger.NewLogger()

	t.Run("When GetTransaction succeeds", func(t *testing.T) {
		mockRepo := mocks.NewTransactionRepositoryMock()
		service := services.NewTransactionService(log, mockRepo)

		transactionEntity := transaction.Transaction{ID: "1", UserID: "1", Amount: 100}

		mockRepo.On("FindByID", ctx, "1").Return(transactionEntity, nil)

		transactionEntity, err := service.GetTransaction(ctx, "1")
		assert.Nil(t, err)
		assert.Equal(t, transactionEntity, transactionEntity)
		mockRepo.AssertCalled(t, "FindByID", ctx, "1")
	})

	t.Run("When GetTransaction fails with not found error", func(t *testing.T) {
		mockRepo := mocks.NewTransactionRepositoryMock()
		service := services.NewTransactionService(log, mockRepo)

		expectedError := errors.New(transaction.NotFoundError)

		mockRepo.On("FindByID", ctx, "1").Return(transaction.Transaction{}, expectedError)

		transactionEntity, err := service.GetTransaction(ctx, "1")
		assert.Equal(t, expectedError, err)
		assert.Equal(t, transaction.Transaction{}, transactionEntity)
		mockRepo.AssertCalled(t, "FindByID", ctx, "1")
	})

	t.Run("When GetTransaction fails with not found error because of logic deletion", func(t *testing.T) {
		mockRepo := mocks.NewTransactionRepositoryMock()
		service := services.NewTransactionService(log, mockRepo)

		expectedError := errors.New(transaction.NotFoundError)

		mockRepo.On("FindByID", ctx, "1").Return(transaction.Transaction{IsDeleted: true}, expectedError)

		transactionEntity, err := service.GetTransaction(ctx, "1")
		assert.Equal(t, expectedError, err)
		assert.Equal(t, transaction.Transaction{IsDeleted: true}, transactionEntity)
		mockRepo.AssertCalled(t, "FindByID", ctx, "1")
	})
}

func TestTransactionService_DeleteTransaction(t *testing.T) {
	ctx := context.TODO()
	log := logger.NewLogger()

	t.Run("When DeleteTransaction succeeds", func(t *testing.T) {
		mockRepo := mocks.NewTransactionRepositoryMock()
		service := services.NewTransactionService(log, mockRepo)

		mockRepo.On("Delete", ctx, "1").Return(nil)

		err := service.DeleteTransaction(ctx, "1")
		assert.Nil(t, err)
		mockRepo.AssertCalled(t, "Delete", ctx, "1")
	})

	t.Run("When FindByID fails in DeleteTransaction", func(t *testing.T) {
		mockRepo := mocks.NewTransactionRepositoryMock()
		service := services.NewTransactionService(log, mockRepo)

		expectedError := errors.New("transaction not found")

		mockRepo.On("Delete", ctx, "1").Return(expectedError)

		err := service.DeleteTransaction(ctx, "1")
		assert.Equal(t, expectedError, err)
		mockRepo.AssertCalled(t, "Delete", ctx, "1")
	})
}
