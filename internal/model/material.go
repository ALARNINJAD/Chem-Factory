package model

type Material struct {
	ID                   int    `json:"id"`
	Name                 string `json:"name"`
	SellPrice            int    `json:"sell_price"`
	BuyPrice             int    `json:"buy_price"`
	FirstIngredientName  string `json:"first_ingredient_name"`
	SecondIngredientName string `json:"second_ingredient_name"`
	MixTime              int    `json:"mix_time"`
}
