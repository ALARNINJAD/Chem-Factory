package inventory

type InventoryAddRequest struct {
	Token  string `json:"token"`
	Name   string `json:"name"`
	Number int    `json:"number"`
}

type InventoryAddResponse struct {
	Massage string `json:"massage"`
}

type InventoryRemoveRequest struct {
	Token  string `json:"token"`
	Name   string `json:"name"`
	Number int    `json:"number"`
}

type InventoryRemoveResponse struct {
	Massage string `json:"massage"`
}

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