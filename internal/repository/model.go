package repository

type material struct {
	id                 int    `db:"id" json:"id"`
	name               string `db:"name" json:"name"`
	sellPrice          int    `db:"sell_price" json:"sell_price"`
	buyPrice           int    `db:"buy_price" json:"buy_price"`
	firstIngredientID  *int   `db:"first_ingredient_id" json:"first_ingredient_id"`
	secondIngredientID *int   `db:"second_ingredient_id" json:"second_ingredient_id"`
	mixTime            int    `db:"mix_time" json:"mix_time"`
}

type inventory struct {
	id         int `db:"id" json:"id"`
	userID     int `db:"user_id" json:"user_id"`
	materialID int `db:"material_id" json:"material_id"`
	number     int `db:"number" json:"number"`
}
