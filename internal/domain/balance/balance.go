package balance

import "math"

const (
	roundToNumber = 100
)

type UserBalance struct {
	Balance      float64 `json:"balance"`
	TotalDebits  int     `json:"total_debits"`
	TotalCredits int     `json:"total_credits"`
}

func (u *UserBalance) RoundBalanceToTwoDecimalPlaces() {
	u.Balance = math.Round(u.Balance*roundToNumber) / roundToNumber
}
