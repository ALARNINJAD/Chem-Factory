package shop

type ShopBuyRequest struct {
	Token          string `json:"token"`
	SellerUsername string `json:"seller_username"`
	MaterialName   string `json:"material_name"`
	Number         int    `json:"number"`
	Price          int    `json:"price"`
}

type ShopBuyResponse struct {
	Massage string `json:"massage"`
}
