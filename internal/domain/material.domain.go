package domain

type Material struct {
	ID                 uint   `json:"id"`
	UserID             uint   `json:"user_id,omitempty"`
	FirstIngredientID  uint   `json:"first_ingredient_id,omitempty"`
	SecondIngredientID uint   `json:"second_ingredient_id,omitempty"`
	Name               string `json:"name"`
	Price              int    `json:"price"`
	MixTime            int    `json:"mix_time,omitempty"`
}

func (m Material) IsRaw() bool {
	return m.FirstIngredientID == 0 && m.SecondIngredientID == 0
}
