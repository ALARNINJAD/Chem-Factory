package domain

import "time"

type Mix struct {
	ID                 uint      `json:"id"`
	UserID             uint      `json:"user_id"`
	FirstIngredientID  uint      `json:"first_ingredient_id"`
	SecondIngredientID uint      `json:"second_ingredient_id"`
	Amount             int       `json:"amount"`
	DateTime           time.Time `json:"date_time"`
}

func (mixer Mix) RemainingSeconds(mixTime int) int {
	elapsed := int(time.Since(mixer.DateTime).Seconds())
	remaining := mixTime - elapsed
	if remaining < 0 {
		return 0
	}
	return remaining
}
