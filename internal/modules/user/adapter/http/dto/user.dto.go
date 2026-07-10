package dto

type ProfileResponse struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username"`
	Balance  int    `json:"balance"`
	XP       int    `json:"xp"`
	Level    int    `json:"level"`
}