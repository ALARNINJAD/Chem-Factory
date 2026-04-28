package user

type UserDataRequest struct {
	Token string `json:"token"`
}

type UserDataResponse struct {
	Username string `json:"username"`
	Balance  int    `json:"balace"`
	XP       int    `json:"xp"`
	Level    int    `json:"level"`
}

type UserRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegisterResponse struct {
	Massage string `json:"massage"`
}

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}