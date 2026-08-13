package dto

import (
	"chem-factory/pkg/dto"
	"time"
)

type MixesResponse struct {
	Mixes []MixResponse `json:"mixes"`
}

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
	MaterialName         string    `json:"material_name"`
	FirstIngredientName  string    `json:"first_ingredient_name"`
	SecondIngredientName string    `json:"second_ingredient_name"`
	Amount               int       `json:"amount"`
	DateTime             time.Time `json:"date_time"`
	RemainingSeconds     int       `json:"remaining_seconds"`
	IsNew                bool      `json:"is_new"`
	dto.MessageResponse
}

// response is mix response
type CheckRequest struct {
	ID uint `json:"id"`
}

type PickRequest struct {
	ID uint `json:"id"`
}
type PickResponse struct {
	IsPicked         bool `json:"is_picked"`
	IsNew            bool `json:"is_new"`
	RemainingSeconds int  `json:"remaining_seconds"`
	dto.MessageResponse
}

// response is mix response
type NewMaterialRequest struct {
	MixID   uint   `json:"id"`
	Name    string `json:"name"`
	Price   int    `json:"price"`
	MixTime int    `json:"mix_time,omitempty"`
}
