package mixer

type CheckMixRequest struct {
	Token string `json:"token"`
	ID    int    `json:"id"`
}

type CheckMixResponse struct {
	Time      int  `json:"time"`
	NewStatus bool `json:"new_status"`
}

type PickMixRequest struct {
	Token string `json:"token"`
	ID    int    `json:"id"`
}

type PickNewMixRequest struct {
	Token   string `json:"token"`
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
	MixTime int    `json:"mix_time"`
}

type MixerAddRequest struct {
	Token                string `json:"token"`
	FirstIngredientName  string `json:"first_ingredient_name"`
	SecondIngredientName string `json:"second_ingredient_name"`
	Number               int    `json:"number"`
}

type MixerAddResponse struct {
	ID int `json:"id"`
}