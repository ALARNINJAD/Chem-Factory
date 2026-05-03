package material

type Material struct {
	Name                 string `json:"name"`
	Username             string `json:"username"`
	FirstIngredientName  string `json:"first_ingredient_name"`
	SecondIngredientName string `json:"second_ingredient_name"`
	SellPrice            int    `json:"sell_price"`
	BuyPrice             int    `json:"buy_price"`
	MixTime              int    `json:"mix_time"`
}

type BaseMaterial struct {
	Name      string `json:"name"`
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	SellPrice int    `json:"sell_price"`
	BuyPrice  int    `json:"buy_price"`
}
