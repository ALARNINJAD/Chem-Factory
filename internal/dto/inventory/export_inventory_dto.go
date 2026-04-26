package inventory

type Inventory struct {
	Username     string `json:"username"`
	MaterialName string `json:"material_name"`
	Number       int    `json:"number"`
}

type InventoryExportRequest struct {
	Token string `json:"token"`
}

type InventoryExportResponse struct {
	InventoryList []Inventory `json:"inventory_list"`
}
