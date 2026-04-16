package model

type material struct {
	Name                 string `json:"name"`
	FirstIngredientName  string `json:"first_ingredient_name"`
	SecondIngredientName string `json:"second_ingredient_name"`
	SellPrice            int    `json:"sell_price"`
	BuyPrice             int    `json:"buy_price"`
	MixTime              int    `json:"mix_time"`
}
