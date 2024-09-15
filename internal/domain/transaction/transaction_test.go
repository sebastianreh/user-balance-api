package transaction_test

import (
	"github.com/sebastianreh/user-balance-api/internal/domain/transaction"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_CreateTransactionByRecord(t *testing.T) {
	t.Run("When record is valid", func(t *testing.T) {
		record := []string{"1", "123", "100.50", "2024-09-13T10:00:00Z"}
		expectedTime, _ := time.Parse(time.RFC3339, "2024-09-13T10:00:00Z")
		expectedAmount, _ := strconv.ParseFloat("100.50", 64)

		transactionEntity, err := transaction.CreateTransactionByRecord(record)

		assert.Nil(t, err)
		assert.Equal(t, "1", transactionEntity.ID)
		assert.Equal(t, "123", transactionEntity.UserID)
		assert.Equal(t, expectedAmount, transactionEntity.Amount)
		assert.Equal(t, expectedTime, *transactionEntity.DateTime)
	})

	t.Run("When amount is not a valid float", func(t *testing.T) {
		record := []string{"1", "123", "invalid_amount", "2024-09-13T10:00:00Z"}

		transactionEntity, err := transaction.CreateTransactionByRecord(record)

		assert.NotNil(t, err)
		assert.Equal(t, "strconv.ParseFloat: parsing \"invalid_amount\": invalid syntax", err.Error())
		assert.Equal(t, transaction.Transaction{}, transactionEntity)
	})

	t.Run("When datetime is not a valid RFC3339 format", func(t *testing.T) {
		record := []string{"1", "123", "100.50", "invalid_datetime"}

		transactionEntity, err := transaction.CreateTransactionByRecord(record)

		assert.NotNil(t, err)
		assert.Equal(t, "parsing time \"invalid_datetime\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"invalid_datetime\" as \"2006\"", err.Error())
		assert.Equal(t, transaction.Transaction{}, transactionEntity)
	})

	t.Run("When both amount and datetime are invalid", func(t *testing.T) {
		record := []string{"1", "123", "invalid_amount", "invalid_datetime"}

		transactionEntity, err := transaction.CreateTransactionByRecord(record)

		assert.NotNil(t, err)
		assert.Equal(t, "strconv.ParseFloat: parsing \"invalid_amount\": invalid syntax", err.Error())
		assert.Equal(t, transaction.Transaction{}, transactionEntity)
	})
}
