package dto

import "time"

type MixRequest struct {
	FirstIngredientID  uint `json:"first_ingredient_id"`
	SecondIngredientID uint `json:"second_ingredient_id"`
	Amount             int  `json:"amount"`
}
type MixResponse struct {
	ID                   uint      `json:"id"`
	UserID               uint      `json:"user_id"`
	FirstIngredientID    uint      `json:"first_ingredient_id"`
	SecondIngredientID   uint      `json:"second_ingredient_id"`
	Username             string    `json:"username"`
	FirstIngredientName  string    `json:"first_ingredient_name"`
	SecondIngredientName string    `json:"second_ingredient_name"`
	Amount               int       `json:"amount"`
	DateTime             time.Time `json:"date_time"`
	RemainingSeconds     int       `json:"remaining_seconds"`
}

type CheckRequest struct {
	ID uint `json:"id"`
}

// type CheckResponse struct {}

// type GetMixesRequest struct {}
type MixesResponse struct {
	Mixes MixResponse `json:"mixes"`
}
