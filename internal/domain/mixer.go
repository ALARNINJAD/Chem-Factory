package domain

import "time"

type Mixer struct {
	ID                 int       `json:"id"`
	UserID             int       `json:"user_id"`
	FirstIngredientID  int       `json:"first_ingredient_id"`
	SecondIngredientID int       `json:"second_ingredient_id"`
	Amount             int       `json:"amount"`
	DateTime           time.Time `json:"date_time"`
}

func (mixer Mixer) RemainingSeconds(mixTime int) int {
    elapsed := int(time.Since(mixer.DateTime).Seconds())
    remaining := mixTime - elapsed
    if remaining < 0 {
        return 0
    }
    return remaining
}