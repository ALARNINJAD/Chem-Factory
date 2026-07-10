package dto

import "time"

type Market struct {
	ID           int       `json:"id"`
	UserID       int       `json:"user_id"`
	MaterialID   int       `json:"material_id"`
	Username     string    `json:"username"`
	MaterialName string    `json:"material_name"`
	Amount       int       `json:"amount"`
	Price        int       `json:"price"`
	DateTime     time.Time `json:"date_time"`
}

type MarketListResponse struct {
	MarketList []Market `json:"market_list"`
}

type SetForSellRequest struct {
	MaterialID  int `json:"material_id"`
	Amount      int `json:"amount"`
	Price       int `json:"price"`
}

type BuyRequest struct {
	MarketID   int `json:"market_id"`
	Amount     int `json:"amount"`
}
