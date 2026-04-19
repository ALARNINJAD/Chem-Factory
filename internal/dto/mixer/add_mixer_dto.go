package mixer

type MixerAddRequest struct {
	Token                string `json:"token"`
	FirstIngredientName  string `json:"first_ingredient_name"`
	SecondIngredientName string `json:"second_ingredient_name"`
	Number               int    `json:"number"`
}

type MixerAddResponse struct {
	ID int `json:"id"`
}
