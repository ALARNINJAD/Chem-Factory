package model

type Mixer struct {
	Username             string `json:"username"`
	FirstIngredientName  string `json:"first_ingredient_name"`
	SecondIngredientName string `json:"second_ingredient_name"`
	Number               int    `json:"number"`
}
