package domain

import (
	"errors"
	"time"
)

type market struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	MaterialID int       `json:"material_id"`
	Amount     int       `json:"amount"`
	Price      int       `json:"price"`
	DateTime   time.Time `json:"date_time"`
}

func (market *market) IncreaseAmount(amount int) error {
	if amount < 0 {
		return errors.New("increase amount must be positive")
	}
	market.Amount += amount
	return nil
}

func (market *market) ReduceAmount(amount int) error {
	if amount < 0 {
		return errors.New("market reduce amount must be positive")
	}
	if market.Amount < amount {
		return errors.New("market insufficient")
	}
	market.Amount -= amount
	return nil
}
