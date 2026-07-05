package domain

type Material struct {
	ID                 int    `json:"id"`
	UserID             int    `json:"user_id"`
	FirstIngredientID  int    `json:"first_ingredient_id,omitempty"`
	SecondIngredientID int    `json:"second_ingredient_id,omitempty"`
	Name               string `json:"name"`
	Price              int    `json:"price"`
	MixTime            int    `json:"mix_time,omitempty"`
}

