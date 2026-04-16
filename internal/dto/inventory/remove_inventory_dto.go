package inventory

type InventoryRemoveRequest struct {
	Token  string `json:"token"`
	Name   string `json:"name"`
	Number int    `json:"number"`
}

type InventoryRemoveResponse struct {
	Massage string `json:"massage"`
}
