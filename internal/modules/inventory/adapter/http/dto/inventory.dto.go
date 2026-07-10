package dto

import "time"

type Inventory struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	MaterialID   int       `json:"material_id"`
	MaterialName string    `json:"material_name"`
	Amount       int       `json:"amount"`
	DateTime     time.Time `json:"date_time"`
}

type ExportResponse struct {
	InventoryList []Inventory `json:"inventory_list"`
}
