package model

import "time"

type Shop struct {
	ID         int       `json:"id,omitempty"`
	UserID     int       `json:"user_id,omitempty"`
	MaterialID int       `json:"material_id,omitempty"`
	Name       string    `json:"name"`
	Number     int       `json:"number"`
	SellPrice  int       `json:"sell_price"`
	DateTime   time.Time `json:"date_time,omitempty"`
}
