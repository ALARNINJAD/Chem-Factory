package mixer

type CheckRequest struct {
	Token string `json:"token"`
	ID    int    `json:"id"`
}

type CheckResponse struct {
	Time      int  `json:"time"`
	NewStatus bool `json:"new_status"`
}

type PickRequest struct {
	Token string `json:"token"`
	ID    int    `json:"id"`
}

type PickNewRequest struct {
	Token   string `json:"token"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
	MixTime int    `json:"mix_time"`
}

type AddRequest struct {
	Token                string `json:"token"`
	FirstIngredientName  string `json:"first_ingredient_name"`
	SecondIngredientName string `json:"second_ingredient_name"`
	Number               int    `json:"number"`
}

type MixerAddResponse struct {
	ID int `json:"id"`
}