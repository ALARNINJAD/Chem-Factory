package model

type user struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Balance  int    `json:"balance"`
	XP       int    `json:"xp"`
	Level    int    `json:"level"`
}
