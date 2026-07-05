package dto

type ProfileResponse struct {
	Username string `json:"username"`
	Balance  int    `json:"balance"`
	XP       int    `json:"xp"`
	Level    int    `json:"level"`
}