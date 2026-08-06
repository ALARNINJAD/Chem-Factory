package dto

import "time"

type MarketResponse struct {
	ID           uint      `json:"id"`
	UserID       uint      `json:"user_id,omitempty"`
	MaterialID   uint      `json:"material_id"`
	Username     string    `json:"username"`
	MaterialName string    `json:"material_name"`
	Amount       int       `json:"amount"`
	Price        int       `json:"price"`
	DateTime     time.Time `json:"date_time"`
}

type MarketListResponse struct {
	MarketList []MarketResponse `json:"market_list"`
}

type SetForSellRequest struct {
	MaterialID uint `json:"material_id"`
	Amount     int  `json:"amount"`
}

type BuyRequest struct {
	MarketID uint `json:"market_id"`
	Amount   int  `json:"amount"`
}
