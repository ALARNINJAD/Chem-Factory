package material

type Material struct {
	Name                 string `json:"name"`
	Username             string `json:"username"`
	FirstIngredientName  string `json:"first_ingredient_name"`
	SecondIngredientName string `json:"second_ingredient_name"`
	Price int    `json:"price"`
	MixTime              int    `json:"mix_time"`
}

type Base struct {
	Name      string `json:"name"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	Price int    `json:"price"`
}
