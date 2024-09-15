package balance_test

import (
	"testing"

	"github.com/sebastianreh/user-balance-api/internal/domain/balance"
	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"github.com/stretchr/testify/assert"
)

func Test_CalculateBalanceByUser(t *testing.T) {
	calculator := balance.NewBalanceCalculator()

	t.Run("When transactions include both debits and credits", func(t *testing.T) {
		transactions := []transaction.Transaction{
			{Amount: 100.00},
			{Amount: -50.00},
			{Amount: 25.00},
		}

		expectedBalance := balance.UserBalance{
			Balance:      75.00,
			TotalDebits:  1,
			TotalCredits: 2,
		}

		result := calculator.CalculateBalanceByUser(transactions)

		assert.Equal(t, expectedBalance.Balance, result.Balance)
		assert.Equal(t, expectedBalance.TotalDebits, result.TotalDebits)
		assert.Equal(t, expectedBalance.TotalCredits, result.TotalCredits)
	})

	t.Run("When transactions only have debits", func(t *testing.T) {
		transactions := []transaction.Transaction{
			{Amount: -100.00},
			{Amount: -50.00},
		}

		expectedBalance := balance.UserBalance{
			Balance:      -150.00,
			TotalDebits:  2,
			TotalCredits: 0,
		}

		result := calculator.CalculateBalanceByUser(transactions)

		assert.Equal(t, expectedBalance.Balance, result.Balance)
		assert.Equal(t, expectedBalance.TotalDebits, result.TotalDebits)
		assert.Equal(t, expectedBalance.TotalCredits, result.TotalCredits)
	})

	t.Run("When transactions only have credits", func(t *testing.T) {
		transactions := []transaction.Transaction{
			{Amount: 100.00},
			{Amount: 50.00},
		}

		expectedBalance := balance.UserBalance{
			Balance:      150.00,
			TotalDebits:  0,
			TotalCredits: 2,
		}

		result := calculator.CalculateBalanceByUser(transactions)

		assert.Equal(t, expectedBalance.Balance, result.Balance)
		assert.Equal(t, expectedBalance.TotalDebits, result.TotalDebits)
		assert.Equal(t, expectedBalance.TotalCredits, result.TotalCredits)
	})

	t.Run("When transactions result in zero balance", func(t *testing.T) {
		transactions := []transaction.Transaction{
			{Amount: 100.00},
			{Amount: -100.00},
		}

		expectedBalance := balance.UserBalance{
			Balance:      0.00,
			TotalDebits:  1,
			TotalCredits: 1,
		}

		result := calculator.CalculateBalanceByUser(transactions)

		assert.Equal(t, expectedBalance.Balance, result.Balance)
		assert.Equal(t, expectedBalance.TotalDebits, result.TotalDebits)
		assert.Equal(t, expectedBalance.TotalCredits, result.TotalCredits)
	})

	t.Run("When transactions have fractional amounts and need rounding", func(t *testing.T) {
		transactions := []transaction.Transaction{
			{Amount: 100.555},
			{Amount: -50.555},
		}

		expectedBalance := balance.UserBalance{
			Balance:      50.00,
			TotalDebits:  1,
			TotalCredits: 1,
		}

		result := calculator.CalculateBalanceByUser(transactions)

		assert.Equal(t, expectedBalance.Balance, result.Balance)
		assert.Equal(t, expectedBalance.TotalDebits, result.TotalDebits)
		assert.Equal(t, expectedBalance.TotalCredits, result.TotalCredits)
	})

	t.Run("When there are no transactions", func(t *testing.T) {
		transactions := []transaction.Transaction{}

		expectedBalance := balance.UserBalance{
			Balance:      0.00,
			TotalDebits:  0,
			TotalCredits: 0,
		}

		result := calculator.CalculateBalanceByUser(transactions)

		assert.Equal(t, expectedBalance.Balance, result.Balance)
		assert.Equal(t, expectedBalance.TotalDebits, result.TotalDebits)
		assert.Equal(t, expectedBalance.TotalCredits, result.TotalCredits)
	})
}
