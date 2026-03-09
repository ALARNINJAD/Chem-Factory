package model

type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
	Balance  int    `json:"balace"`
	XP       int    `json:"xp"`
	Level    int    `json:"level"`
}
