package services

import (
	"context"
	"github.com/sebastianreh/user-balance-api/internal/domain/balance"
	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"github.com/sebastianreh/user-balance-api/internal/domain/user"
	"github.com/sebastianreh/user-balance-api/pkg/logger"
	"github.com/sebastianreh/user-balance-api/pkg/strings"
)

const (
	UserNotFound = "user not found"
)

type BalanceService interface {
	GetBalanceByUserIDWithOptions(ctx context.Context, userID, fromDate, toDate string) (balance.UserBalance, error)
	GetBalanceByUserID(ctx context.Context, userID string) (balance.UserBalance, error)
}

type balanceService struct {
	log                   logger.Logger
	userRepository        user.Repository
	transactionRepository transaction.Repository
	balanceCalculator     balance.Calculator
}

func NewBalanceService(logger logger.Logger, userRepository user.Repository,
	transactionRepository transaction.Repository, balanceCalculator balance.Calculator) BalanceService {
	return &balanceService{
		log:                   logger,
		userRepository:        userRepository,
		transactionRepository: transactionRepository,
		balanceCalculator:     balanceCalculator,
	}
}

func (s balanceService) GetBalanceByUserIDWithOptions(ctx context.Context, userID, fromDate, toDate string) (balance.UserBalance, error) {
	var userBalance balance.UserBalance
	_, err := s.userRepository.FindByID(ctx, userID)
	if err != nil {
		return userBalance, err
	}

	transactions, err := s.transactionRepository.FindByUserIDWithOptions(ctx, userID, fromDate, toDate)
	if err != nil {
		return userBalance, err
	}

	userBalance = s.balanceCalculator.CalculateBalanceByUser(transactions)

	return userBalance, nil
}

func (s balanceService) GetBalanceByUserID(ctx context.Context, userID string) (balance.UserBalance, error) {
	return s.GetBalanceByUserIDWithOptions(ctx, userID, customStr.Empty, customStr.Empty)
}
