package model

type Inventory struct {
	ID         int `json:"id"`
	UserID     int `json:"user_id"`
	MaterialID int `json:"material_id"`
	Number     int `json:"number"`
}
