package user

type UserRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegisterResponse struct {
	Massage string `json:"massage"`
}
