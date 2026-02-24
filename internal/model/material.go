package model

type Material struct {
	ID                   int    `db:"id" json:"id,omitempty"`
	Name                 string `db:"name" json:"name"`
	SellPrice            int    `db:"sell_price" json:"sell_price"`
	BuyPrice             int    `db:"buy_price" json:"buy_price"`
	MixTime              int    `db:"mix_time" json:"mix_time"`
	FirstIngredientName  string `json:"first_ingredient_name"`
	SecondIngredientName string `json:"second_ingredient_name"`
	FirstIngredientID    int    `db:"first_ingredient_id,omitempty"`
	SecondIngredientID   int    `db:"second_ingredient_id,omitempty"`
}
