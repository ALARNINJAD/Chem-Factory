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
