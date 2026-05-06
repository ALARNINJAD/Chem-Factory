package user

type DataRequest struct {
	Token string `json:"token"`
}

type DataResponse struct {
	Username string `json:"username"`
	Balance  int    `json:"balace"`
	XP       int    `json:"xp"`
	Level    int    `json:"level"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Massage string `json:"massage"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}