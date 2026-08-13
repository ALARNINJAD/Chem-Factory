package domain

import (
	"time"
)

type Inventory struct {
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	MaterialID uint      `json:"material_id"`
	Amount     int       `json:"amount"`
	DateTime   time.Time `json:"date_time"`
}
