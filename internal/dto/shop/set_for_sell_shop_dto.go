package shop

type ShopSetForSellRequest struct {
	Token        string `json:"token"`
	MaterialName string `json:"material_name"`
	Number       int    `json:"number"`
	Price        int    `json:"price"`
}

type ShopSetForSellResponse struct {
	Massage string `json:"massage"`
}
