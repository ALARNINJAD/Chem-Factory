package repository

import (
	"database/sql"
	"fmt"
	"log"
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

func (r *repositoryManager) EmptyMixerStruct() *mixer {
	return &mixer{}
}

func (r *repositoryManager) FindMixIDByUserIDIngrID(userID, firstID, secID int) (int, error) {

	var id int
	err := r.db.QueryRow(`
		SELECT id FROM mixer
		WHERE user_id = ? AND first_ingredient_id = ? AND second_ingredient_id = ?`,
		userID, firstID, secID).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Repository mixer, find mix id by user&ingr id: %w ", err)
	}
	return id, nil
}

func (r *repositoryManager) AddToMixer(tx *sql.Tx, m mixer) error {

	_, err := tx.Exec(`
		INSERT INTO mixer(
		user_id, first_ingredient_id, second_ingredient_id,
		username, first_ingredient_name, second_ingredient_name, number)
		VALUES(?,?,?,?,?,?,?)`,
		m.UserID, m.FirstIngredientID, m.SecondIngredientID,
		m.Username, m.FirstIngredientName, m.SecondIngredientName, m.Number)
	if err != nil {
		return fmt.Errorf("Repository mixer, add to mixer: %w ", err)
	}

	return nil
}

func (r *repositoryManager) FindMixRowByID(id int) (*mixer, error) {

	var m mixer

	err := r.db.QueryRow(`
		SELECT id, user_id, first_ingredient_id, second_ingredient_id,
		username, first_ingredient_name, second_ingredient_name,
		number, date_time FROM mixer WHERE id = ?`, id,
	).Scan(&m.ID, &m.UserID, &m.FirstIngredientID, &m.SecondIngredientID,
		&m.Username, &m.FirstIngredientName, &m.SecondIngredientName,
		&m.Number, &m.DateTime)
	if err != nil {
		log.Println(id)
		log.Println(m)
		return nil, fmt.Errorf("Repository mixer, find mix row by id: %w ", err)
	}

	return &m, nil
}

func (r *repositoryManager) DeleteMixByID(tx *sql.Tx, id int) error {

	_, err := tx.Exec("DELETE FROM mixer WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("Repository mixer, delelte mixer by id: %w ", err)
	}
	return nil
}
