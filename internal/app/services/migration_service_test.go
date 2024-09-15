package services_test

import (
	"context"
	"errors"
	"github.com/sebastianreh/user-balance-api/internal/infrastructure/config"
	"mime/multipart"
	"testing"

	"github.com/sebastianreh/user-balance-api/internal/app/services"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	"github.com/sebastianreh/user-balance-api/test/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_MigrationService_ProcessBalance(t *testing.T) {
	ctx := context.TODO()
	cfg := config.NewConfig()
	loggerMock := logger.NewLogger()
	fileHeader := &multipart.FileHeader{}

	t.Run("When ReadFile returns an error", func(t *testing.T) {
		expectedError := errors.New(services.ReadFileError)

		csvProcessor := mocks.NewCsvProcessorMock()
		csvProcessor.On("ReadFile", fileHeader, mock.Anything).Return([][]string{}, expectedError)

		userRepo := mocks.NewUserRepositoryMock()
		transactionRepo := mocks.NewTransactionRepositoryMock()

		service := services.NewMigrationService(cfg, loggerMock, userRepo, transactionRepo, csvProcessor)

		err := service.ProcessBalance(ctx, fileHeader)

		assert.NotNil(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("When transaction.CreateTransactionByRecord returns an error", func(t *testing.T) {
		records := [][]string{{"1", "test_user", "100.00", "2024-09-13"}}
		expectedError := errors.New("error creating transaction by record: parsing time \"2024-09-13\"" +
			" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"\" as \"T\"")

		csvProcessor := mocks.NewCsvProcessorMock()
		csvProcessor.On("ReadFile", fileHeader, mock.Anything).Return(records, nil)

		userRepo := mocks.NewUserRepositoryMock()
		userRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

		transactionRepo := mocks.NewTransactionRepositoryMock()

		service := services.NewMigrationService(cfg, loggerMock, userRepo, transactionRepo, csvProcessor)

		err := service.ProcessBalance(ctx, fileHeader)

		assert.Error(t, err)
		assert.Equal(t, err, expectedError)
	})

	t.Run("When userRepository.Save returns an error", func(t *testing.T) {
		records := [][]string{{"1", "test_user", "100.00", "2024-09-13T10:00:00Z"}}
		expectedError := errors.New("error saving user: repository error")

		csvProcessor := mocks.NewCsvProcessorMock()
		csvProcessor.On("ReadFile", fileHeader, mock.Anything).Return(records, nil)

		userRepo := mocks.NewUserRepositoryMock()
		userRepo.On("Save", mock.Anything, mock.Anything).Return(errors.New("repository error"))

		transactionRepo := mocks.NewTransactionRepositoryMock()

		service := services.NewMigrationService(cfg, loggerMock, userRepo, transactionRepo, csvProcessor)

		err := service.ProcessBalance(ctx, fileHeader)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
	})

	t.Run("When transactionRepository.SaveBatch returns an error", func(t *testing.T) {
		records := [][]string{{"1", "test_user", "100.00", "2024-09-13T10:00:00Z"}}
		expectedError := errors.New("error saving transaction batch: repository error")

		csvProcessor := mocks.NewCsvProcessorMock()
		csvProcessor.On("ReadFile", fileHeader, mock.Anything).Return(records, nil)

		userRepo := mocks.NewUserRepositoryMock()
		userRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

		transactionRepo := mocks.NewTransactionRepositoryMock()
		transactionRepo.On("SaveBatch", mock.Anything, mock.Anything).Return(errors.New("repository error"))

		service := services.NewMigrationService(cfg, loggerMock, userRepo, transactionRepo, csvProcessor)

		err := service.ProcessBalance(ctx, fileHeader)

		assert.Error(t, err)
		assert.Equal(t, err, expectedError)
	})

	t.Run("When ProcessBalance completes successfully", func(t *testing.T) {
		records := [][]string{{"1", "test_user", "100.00", "2024-09-13T10:00:00Z"}}

		csvProcessor := mocks.NewCsvProcessorMock()
		csvProcessor.On("ReadFile", fileHeader, mock.Anything).Return(records, nil)

		userRepo := mocks.NewUserRepositoryMock()
		userRepo.On("Save", mock.Anything, mock.Anything).Return(nil)

		transactionRepo := mocks.NewTransactionRepositoryMock()
		transactionRepo.On("SaveBatch", mock.Anything, mock.Anything).Return(nil)

		service := services.NewMigrationService(cfg, loggerMock, userRepo, transactionRepo, csvProcessor)

		err := service.ProcessBalance(ctx, fileHeader)

		assert.Nil(t, err)
	})
}
