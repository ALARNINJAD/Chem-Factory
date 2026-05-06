package shop

type Shop struct {
	Username     string `json:"username"`
	MaterialName string `json:"material_name"`
	Number       int    `json:"number"`
	Price        int    `json:"price"`
}

type UserIDMatIDPrice struct {
	UserID     int `json:"user_id"`
	MaterialID int `json:"material_id"`
	Price      int `json:"price"`
}

type BaseMatRequest struct {
	Token string `json:"token"`
	ID    int    `json:"id"`
}

type BuyRequest struct {
	Token          string `json:"token"`
	SellerUsername string `json:"seller_username"`
	MaterialName   string `json:"material_name"`
	Number         int    `json:"number"`
	Price          int    `json:"price"`
}

type BuyResponse struct {
	Massage string `json:"massage"`
}

type SetForSellRequest struct {
	Token        string `json:"token"`
	MaterialName string `json:"material_name"`
	Number       int    `json:"number"`
	Price        int    `json:"price"`
}

type SetForSellResponse struct {
	Massage string `json:"massage"`
}

type ItemsResponse struct {
	Items []Shop `json:"items"`
}
