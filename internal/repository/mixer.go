package repository

import (
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

func createMixerTable(r *repositoryManager) {
	query := `
	CREATE TABLE IF NOT EXISTS mixer (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		first_ingredient_id INTEGER,
		second_ingredient_id INTEGER,
		username TEXT NOT NULL,
		first_ingredient_name TEXT NOT NULL,
		second_ingredient_name TEXT NOT NULL,
		number INTEGER NOT NULL CHECK("number" >= 1),
		date_time DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		UNIQUE("user_id","first_ingredient_id","second_ingredient_id"),
		UNIQUE("username","first_ingredient_name","second_ingredient_name"),
		FOREIGN KEY("user_id") REFERENCES "user"("id")
		FOREIGN KEY("first_ingredient_id") REFERENCES "material"("id"),
		FOREIGN KEY("second_ingredient_id") REFERENCES "material"("id")
	)`
	if _, err := r.db.Exec(query); err != nil {
		panic(fmt.Errorf("Repository mixer, create mixer table: %w ", err))
	}
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

func (r *repositoryManager) AddToMixer(m mixer) error {

	_, err := r.db.Exec(`
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
		number, date_time
		FROM mixer WHERE id = ?`, id,
	).Scan(&m.ID, &m.UserID, &m.FirstIngredientID, &m.SecondIngredientID,
		&m.Username, &m.FirstIngredientName, &m.SecondIngredientName,
		&m.Number, &m.DateTime)
	if err != nil {
		return nil, fmt.Errorf("Repository mixer, find mix row by id: %w ", err)
	}

	return &m, nil
}
