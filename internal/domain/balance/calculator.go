package balance

import (
	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
)

type Calculator interface {
	CalculateBalanceByUser(transactions []transaction.Transaction) UserBalance
}

type calculator struct {
}

func NewBalanceCalculator() Calculator {
	return &calculator{}
}

func (c calculator) CalculateBalanceByUser(transactions []transaction.Transaction) UserBalance {
	var userBalance UserBalance
	for _, userTransaction := range transactions {
		if userTransaction.Amount < 0 {
			userBalance.TotalDebits++
		}

		if userTransaction.Amount > 0 {
			userBalance.TotalCredits++
		}

		userBalance.Balance += userTransaction.Amount
	}

	userBalance.RoundBalanceToTwoDecimalPlaces()

	return userBalance
}
