package model

type shop struct {
	Username     string `json:"username"`
	MaterialName string `json:"material_name"`
	Number       int    `json:"number"`
	SellPrice    int    `json:"sell_price"`
}
