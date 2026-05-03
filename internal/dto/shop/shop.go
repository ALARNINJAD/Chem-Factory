package shop

type Shop struct {
	Username     string `json:"username"`
	MaterialName string `json:"material_name"`
	Number       int    `json:"number"`
	Price        int    `json:"price"`
}

type ShopUserIDMatIDPrice struct {
	UserID     int `json:"user_id"`
	MaterialID int `json:"material_id"`
	Price      int `json:"price"`
}

// route ---------------------------------------------

// buy basd material

type ShopBaseMatRequest struct {
	Token string `json:"token"`
	ID    int    `json:"id"`
}

// buy material

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

// set for sell

type ShopSetForSellRequest struct {
	Token        string `json:"token"`
	MaterialName string `json:"material_name"`
	Number       int    `json:"number"`
	Price        int    `json:"price"`
}

type ShopSetForSellResponse struct {
	Massage string `json:"massage"`
}

// get all items

type ShopItemsResponse struct {
	Items []Shop `json:"items"`
}
