package domain

import "time"

type Shop struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	MaterialID int       `json:"material_id"`
	Username   string    `json:"username"`
	Number     int       `json:"number"`
	Price      int       `json:"price"`
	DateTime   time.Time `json:"date_time"`
}
