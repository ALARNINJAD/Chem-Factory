package inventory

type InventoryAddRequest struct {
	Token  string `json:"token"`
	Name   string `json:"name"`
	Number int    `json:"number"`
}

type InventoryAddResponse struct {
	Massage string `json:"massage"`
}
