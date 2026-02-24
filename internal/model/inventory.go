package model

type Inventory struct {
	ID         int `db:"id" json:"id,omitempty"`
	UserID     int `db:"user_id" json:"user_id"`
	MaterialID int `db:"material_id" json:"material_id"`
	Number     int `db:"number" json:"number"`
}
