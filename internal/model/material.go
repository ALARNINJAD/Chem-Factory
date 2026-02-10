package model

type Material struct {
	MaterialID         int    `db:"material_id" json:"material_id"`
	Name               string `db:"name" json:"name"`
	SellPrice          int    `db:"sell_price" json:"sell_price"`
	BuyPrice           int    `db:"buy_price" json:"buy_price"`
	FirstIngredientID  int    `db:"first_ingredient_id" json:"first_ingredient_id"`
	SecondIngredientID int    `db:"second_ingredient_id" json:"second_ingredient_id"`
	MixTime            int    `db:"mix_time" json:"mix_time"`
}
