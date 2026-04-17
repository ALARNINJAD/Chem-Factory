package inventory

import "chem-factory/internal/model"

type InventoryExportRequest struct {
	Token string `json:"token"`
}

type InventoryExportResponse struct {
	InventoryList []model.Inventory `json:"inventory_list"`
}
