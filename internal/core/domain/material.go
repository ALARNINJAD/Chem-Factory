package domain

type Material struct {
	ID                   int    `json:"id"`
	UserID               int    `json:"user_id"`
	FirstIngredientID    int    `json:"first_ingredient_id"`
	SecondIngredientID   int    `json:"second_ingredient_id"`
	Name                 string `json:"name"`
	Price                int    `json:"price"`
	MixTime              int    `json:"mix_time"`
}