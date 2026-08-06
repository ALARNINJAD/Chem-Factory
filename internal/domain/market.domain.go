package domain

import (
	"errors"
	"time"
)

type Market struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	MaterialID uint      `json:"material_id"`
	Amount     int       `json:"amount"`
	DateTime   time.Time `json:"date_time"`
}

func (market *Market) IncreaseAmount(amount int) error {
	if amount < 0 {
		return errors.New("increase amount must be positive")
	}
	market.Amount += amount
	return nil
}

func (market *Market) ReduceAmount(amount int) error {
	if amount < 0 {
		return errors.New("market reduce amount must be positive")
	}
	if market.Amount < amount {
		return errors.New("market insufficient")
	}
	market.Amount -= amount
	return nil
}

func (market *Market) IsAdmin() bool { return market.UserID == 0 }
