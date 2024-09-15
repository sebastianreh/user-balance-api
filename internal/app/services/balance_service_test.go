package services_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/sebastianreh/user-balance-api/internal/app/services"
	"github.com/sebastianreh/user-balance-api/internal/domain/balance"
	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"github.com/sebastianreh/user-balance-api/internal/domain/user"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	"github.com/sebastianreh/user-balance-api/test/mocks"
	"github.com/stretchr/testify/assert"
)

func Test_BalanceService_GetBalanceByUserIDWithOptions(t *testing.T) {
	ctx := context.TODO()
	userID := "123"
	fromDate := "2024-01-01"
	toDate := "2024-12-31"
	now := time.Now()

	t.Run("When GetBalanceByUserIDWithOptions success", func(t *testing.T) {
		transactions := []transaction.Transaction{
			{
				ID:       "1",
				UserID:   "1",
				Amount:   100,
				DateTime: &now,
			},
			{
				ID:       "2",
				UserID:   "1",
				Amount:   -200,
				DateTime: &now,
			},
		}
		userEntity := user.User{ID: userID}
		expectedBalance := balance.UserBalance{
			Balance:      -100,
			TotalDebits:  1,
			TotalCredits: 1,
		}

		userRepo := mocks.NewUserRepositoryMock()
		userRepo.On("FindByID", ctx, userID).Return(userEntity, nil)

		transactionRepo := mocks.NewTransactionRepositoryMock()
		transactionRepo.On("FindByUserIDWithOptions", ctx, userID, fromDate, toDate).Return(transactions, nil)

		calculator := mocks.NewCalculatorMock()
		calculator.On("CalculateBalanceByUser", transactions).Return(expectedBalance)

		service := services.NewBalanceService(logger.NewLogger(), userRepo, transactionRepo, calculator)
		userBalance, err := service.GetBalanceByUserIDWithOptions(ctx, userID, fromDate, toDate)

		assert.Nil(t, err)
		assert.NotNil(t, userBalance)
	})

	t.Run("When GetBalanceByUserIDWithOptions user not found", func(t *testing.T) {
		expectedError := errors.New("user not found")

		userRepo := mocks.NewUserRepositoryMock()
		userRepo.On("FindByID", ctx, userID).Return(user.User{}, expectedError)

		transactionRepo := mocks.NewTransactionRepositoryMock()

		calculator := mocks.NewCalculatorMock()

		service := services.NewBalanceService(logger.NewLogger(), userRepo, transactionRepo, calculator)
		userBalance, err := service.GetBalanceByUserIDWithOptions(ctx, userID, fromDate, toDate)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, balance.UserBalance{}, userBalance)
	})

	t.Run("When GetBalanceByUserIDWithOptions transaction repository returns error", func(t *testing.T) {
		userEntity := user.User{ID: userID}

		expectedError := errors.New("transaction repository error")

		userRepo := mocks.NewUserRepositoryMock()
		userRepo.On("FindByID", ctx, userID).Return(userEntity, nil)

		transactionRepo := mocks.NewTransactionRepositoryMock()
		transactionRepo.On("FindByUserIDWithOptions", ctx, userID, fromDate, toDate).Return([]transaction.Transaction{}, expectedError)

		calculator := mocks.NewCalculatorMock()

		service := services.NewBalanceService(logger.NewLogger(), userRepo, transactionRepo, calculator)
		userBalance, err := service.GetBalanceByUserIDWithOptions(ctx, userID, fromDate, toDate)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, balance.UserBalance{}, userBalance)
	})
}

func Test_BalanceService_GetBalanceByUserID(t *testing.T) {
	ctx := context.TODO()
	userID := "123"
	now := time.Now()

	t.Run("When GetBalance success", func(t *testing.T) {
		transactions := []transaction.Transaction{
			{
				ID:       "1",
				UserID:   "1",
				Amount:   100,
				DateTime: &now,
			},
			{
				ID:       "2",
				UserID:   "1",
				Amount:   -200,
				DateTime: &now,
			},
		}
		userEntity := user.User{ID: userID}
		expectedBalance := balance.UserBalance{
			Balance:      -100,
			TotalDebits:  1,
			TotalCredits: 1,
		}

		userRepo := mocks.NewUserRepositoryMock()
		userRepo.On("FindByID", ctx, userID).Return(userEntity, nil)

		transactionRepo := mocks.NewTransactionRepositoryMock()
		transactionRepo.On("FindByUserIDWithOptions", ctx, userID, "", "").Return(transactions, nil)

		calculator := mocks.NewCalculatorMock()
		calculator.On("CalculateBalanceByUser", transactions).Return(expectedBalance)

		service := services.NewBalanceService(logger.NewLogger(), userRepo, transactionRepo, calculator)
		userBalance, err := service.GetBalanceByUserID(ctx, userID)

		assert.Nil(t, err)
		assert.NotNil(t, userBalance)
	})

	t.Run("When GetBalance user not found", func(t *testing.T) {
		expectedError := errors.New("user not found")

		userRepo := mocks.NewUserRepositoryMock()
		userRepo.On("FindByID", ctx, userID).Return(user.User{}, expectedError)

		transactionRepo := mocks.NewTransactionRepositoryMock()

		calculator := mocks.NewCalculatorMock()

		service := services.NewBalanceService(logger.NewLogger(), userRepo, transactionRepo, calculator)
		userBalance, err := service.GetBalanceByUserID(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, balance.UserBalance{}, userBalance)
	})

	t.Run("When GetBalance transaction repository returns error", func(t *testing.T) {
		userEntity := user.User{ID: userID}

		expectedError := errors.New("transaction repository error")

		userRepo := mocks.NewUserRepositoryMock()
		userRepo.On("FindByID", ctx, userID).Return(userEntity, nil)

		transactionRepo := mocks.NewTransactionRepositoryMock()
		transactionRepo.On("FindByUserIDWithOptions", ctx, userID, "", "").Return([]transaction.Transaction{}, expectedError)

		calculator := mocks.NewCalculatorMock()

		service := services.NewBalanceService(logger.NewLogger(), userRepo, transactionRepo, calculator)
		userBalance, err := service.GetBalanceByUserID(ctx, userID)

		assert.Error(t, err)
		assert.Equal(t, expectedError, err)
		assert.Equal(t, balance.UserBalance{}, userBalance)
	})
}
