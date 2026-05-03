package repository

import (
	"database/sql"
	"fmt"
	"time"
)

type mixer struct {
	ID                   int       `json:"id,omitempty"`
	UserID               int       `json:"user_id,omitempty"`
	FirstIngredientID    int       `json:"first_ingredient_id,omitempty"`
	SecondIngredientID   int       `json:"second_ingredient_id,omitempty"`
	Username             string    `json:"username"`
	FirstIngredientName  string    `json:"first_ingredient_name"`
	SecondIngredientName string    `json:"second_ingredient_name"`
	Number               int       `json:"number"`
	DateTime             time.Time `json:"date_time"`
}

type mixerManager struct{ db *sql.DB }

func NewMixerManager(db *sql.DB) *mixerManager { return &mixerManager{db: db} }

func (m *mixerManager) EmptyStruct() *mixer { return &mixer{} }

func (m *mixerManager) FindIDByUserIDIngrID(userID, firstID, secID int) (int, error) {
	var id int
	err := m.db.QueryRow(`
		SELECT id FROM mixer
		WHERE user_id = ? AND first_ingredient_id = ? AND second_ingredient_id = ?`,
		userID, firstID, secID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Repository mixer, find mix id by user id and ingredients id: %w ", err)
	}
	return id, nil
}

func (m *mixerManager) Add(tx *sql.Tx, mxr mixer) error {
	_, err := tx.Exec(`
		INSERT INTO mixer(
		user_id, first_ingredient_id, second_ingredient_id,
		username, first_ingredient_name, second_ingredient_name, number)
		VALUES(?,?,?,?,?,?,?)`,
		mxr.UserID, mxr.FirstIngredientID, mxr.SecondIngredientID,
		mxr.Username, mxr.FirstIngredientName, mxr.SecondIngredientName, mxr.Number)
	if err != nil {
		return fmt.Errorf("Repository mixer, add to mixer: %w ", err)
	}
	return nil
}

func (m *mixerManager) FindByID(id int) (*mixer, error) {
	var mxr mixer
	err := m.db.QueryRow(`
		SELECT id, user_id, first_ingredient_id, second_ingredient_id,
		username, first_ingredient_name, second_ingredient_name,
		number, date_time FROM mixer WHERE id = ?`, id,
	).Scan(&mxr.ID, &mxr.UserID, &mxr.FirstIngredientID, &mxr.SecondIngredientID,
		&mxr.Username, &mxr.FirstIngredientName, &mxr.SecondIngredientName,
		&mxr.Number, &mxr.DateTime)
	if err != nil {
		return nil, fmt.Errorf("Repository mixer, find mix row by id: %w ", err)
	}
	return &mxr, nil
}

func (m *mixerManager) DeleteByID(tx *sql.Tx, id int) error {
	_, err := tx.Exec("DELETE FROM mixer WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("Repository mixer, delelte mixer by id: %w ", err)
	}
	return nil
}
