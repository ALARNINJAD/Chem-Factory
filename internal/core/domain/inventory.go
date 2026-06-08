package domain

import "time"

type Inventory struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	MaterialID   int       `json:"material_id"`
	Username     string    `json:"username"`
	MaterialName string    `json:"material_name"`
	Number       int       `json:"number"`
	DateTime     time.Time `json:"date_time"`
}
