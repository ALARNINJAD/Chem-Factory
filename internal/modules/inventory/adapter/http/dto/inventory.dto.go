package dto

import "time"

type InventoryResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id"`
	MaterialID   uint      `json:"material_id"`
	MaterialName string    `json:"material_name"`
	Amount       int       `json:"amount"`
	DateTime     time.Time `json:"date_time"`
}

type ExportResponse struct {
	InventoryList []InventoryResponse `json:"inventory_list"`
}
