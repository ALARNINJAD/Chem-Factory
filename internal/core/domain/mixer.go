package domain

import "time"

type Mixer struct {
	ID                 int       `json:"id"`
	UserID             int       `json:"user_id"`
	FirstIngredientID  int       `json:"first_ingredient_id"`
	SecondIngredientID int       `json:"second_ingredient_id"`
	Number             int       `json:"number"`
	DateTime           time.Time `json:"date_time"`
}
