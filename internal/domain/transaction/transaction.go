package transaction

import (
	"strconv"
	"time"
)

type Transaction struct {
	ID        string     `json:"id"`
	UserID    string     `json:"user_id"`
	Amount    float64    `json:"amount"`
	DateTime  *time.Time `json:"date_time"`
	IsDeleted bool       `json:"-"`
}

func CreateTransactionByRecord(record []string) (Transaction, error) {
	var transaction Transaction
	amount, err := strconv.ParseFloat(record[2], 64)
	if err != nil {
		return transaction, err
	}

	parsedTime, err := time.Parse(time.RFC3339, record[3])
	if err != nil {
		return transaction, err
	}

	transaction = Transaction{
		ID:       record[0],
		UserID:   record[1],
		Amount:   amount,
		DateTime: &parsedTime,
	}

	return transaction, nil
}
