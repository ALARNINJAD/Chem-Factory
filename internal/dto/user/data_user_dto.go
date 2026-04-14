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
