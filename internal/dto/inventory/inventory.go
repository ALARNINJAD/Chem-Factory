package inventory

type Inventory struct {
	Username     string `json:"username"`
	MaterialName string `json:"material_name"`
	Number       int    `json:"number"`
}

type AddRequest struct {
	Token  string `json:"token"`
	Name   string `json:"name"`
	Number int    `json:"number"`
}

type AddResponse struct {
	Massage string `json:"massage"`
}

type RemoveRequest struct {
	Token  string `json:"token"`
	Name   string `json:"name"`
	Number int    `json:"number"`
}

type RemoveResponse struct {
	Massage string `json:"massage"`
}

type ExportRequest struct {
	Token string `json:"token"`
}

type ExportResponse struct {
	InventoryList []Inventory `json:"inventory_list"`
}