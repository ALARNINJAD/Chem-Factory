package domain

import (
	"errors"
	"time"
)

type Inventory struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	MaterialID int       `json:"material_id"`
	Amount     int       `json:"amount"`
	DateTime   time.Time `json:"date_time"`
}

func (inventory *Inventory) IncreaseAmount(amount int) error {
	if amount < 0 {
		return errors.New("increase amount must be positive")
	}
	inventory.Amount += amount
	return nil
}

func (inventory *Inventory) ReduceAmount(amount int) error {
	if amount < 0 {
		return errors.New("inventory reduce amount must be positive")
	}
	if inventory.Amount < amount {
		return errors.New("inventory insufficient")
	}
	inventory.Amount -= amount
	return nil
}
