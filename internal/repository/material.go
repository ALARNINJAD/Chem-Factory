package repository

import (
	"fmt"
)

type material struct {
	ID                   int    `json:"id,omitempty"`
	UserID               int    `json:"user_id,omitempty"`
	FirstIngredientID    int    `json:"first_ingredient_id,omitempty"`
	SecondIngredientID   int    `json:"second_ingredient_id,omitempty"`
	Username             string `json:"username"`
	FirstIngredientName  string `json:"first_ingredient_name"`
	SecondIngredientName string `json:"second_ingredient_name"`
	Name                 string `json:"name"`
	SellPrice            int    `json:"sell_price"`
	BuyPrice             int    `json:"buy_price"`
	MixTime              int    `json:"mix_time"`
}

func createMaterialTable(r *repositoryManager) {
	query := `
	CREATE TABLE IF NOT EXISTS material (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		first_ingredient_id INTEGER,
		second_ingredient_id INTEGER,
		username TEXT NOT NULL,
		name TEXT NOT NULL UNIQUE,
		first_ingredient_name TEXT NOT NULL,
		second_ingredient_name TEXT NOT NULL,
		sell_price INTEGER NOT NULL CHECK("sell_price" >= 0),
		buy_price INTEGER NOT NULL CHECK("buy_price" >= 0),
		mix_time INTEGER NOT NULL CHECK("mix_time" >= 0),
		UNIQUE("first_ingredient_id","second_ingredient_id"),
		UNIQUE("first_ingredient_name","second_ingredient_name"),
		FOREIGN KEY("user_id") REFERENCES "user"("id")
		FOREIGN KEY("first_ingredient_id") REFERENCES "material"("id"),
		FOREIGN KEY("second_ingredient_id") REFERENCES "material"("id")
	)`
	if _, err := r.db.Exec(query); err != nil {
		panic(fmt.Errorf("Repository material, crate material table: %w", err))
	}
}

func (r *repositoryManager) SaveMaterial(m material) error {

	_, err := r.db.Exec(`
		INSERT material(user_id, first_ingredient_id, second_ingredient_id,
		username, name, first_ingredient_name, second_ingredient_name,
		sell_price, buy_price, mix_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		m.UserID, m.FirstIngredientID, m.SecondIngredientID,
		m.Username, m.Name, m.FirstIngredientName, m.SecondIngredientName,
		m.SellPrice, m.BuyPrice, m.MixTime)
	if err != nil {
		return fmt.Errorf("Repository material, save material: %w", err)
	}

	return nil
}

func (r *repositoryManager) FindIDbyMaterialName(name string) (int, error) {

	var id int
	err := r.db.QueryRow("SELECT id FROM material WHERE name = ?", name).Scan(&id)
	return id, fmt.Errorf("Repository material, find id by material name: %w", err)
}

func (r *repositoryManager) FindBaseMaterials() ([]material, error) {

	rows, err := r.db.Query(`
        SELECT id, user_id, first_ingredient_id, second_ingredient_id,
		username, name, first_ingredient_name, second_ingredient_name,
		sell_price, buy_price, mix_time FROM material 
        WHERE second_ingredient_id = first_ingredient_id`)
	if err != nil {
		return nil, fmt.Errorf("Repository material, export base materials query: %w", err)
	}
	defer rows.Close()

	var list []material

	for rows.Next() {

		var m material
		err := rows.Scan(
			&m.ID,
			&m.UserID,
			&m.FirstIngredientID,
			&m.SecondIngredientID,
			&m.Username,
			&m.Name,
			&m.FirstIngredientName,
			&m.SecondIngredientName,
			&m.SellPrice,
			&m.BuyPrice,
			&m.MixTime,
		)
		if err != nil {
			return nil, fmt.Errorf("Repository material, export base materials scan: %w", err)
		}

		list = append(list, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Repository material, export base materials rows: %w", err)
	}

	return list, nil
}

func (r *repositoryManager) FindMaterialByID(id int) (*material, error) {

	var m material

	err := r.db.QueryRow(`
		SELECT id, user_id, first_ingredient_id, second_ingredient_id,
		username, name, first_ingredient_name, second_ingredient_name,
		sell_price, buy_price, mix_time
		FROM material WHERE id = ? `, id,
	).Scan(&m.ID, &m.UserID, &m.FirstIngredientID, &m.SecondIngredientID,
		&m.Username, &m.Name, &m.FirstIngredientName, &m.SecondIngredientName,
		&m.SellPrice, &m.BuyPrice, &m.MixTime)
	if err != nil {
		return nil, fmt.Errorf("Repository material, find material by id: %w", err)
	}

	return &m, nil
}

func (r *repositoryManager) FindMaterialByIngrID(firstID int, secondID int) (*material, error) {

	var m material

	err := r.db.QueryRow(`
		SELECT id, user_id, first_ingredient_id, second_ingredient_id,
		username, name, first_ingredient_name, second_ingredient_name,
		sell_price, buy_price, mix_time
		FROM material  WHERE first_ingredient_id = ? AND second_ingredient_id = ?`, firstID, secondID,
	).Scan(&m.ID, &m.UserID, &m.FirstIngredientID, &m.SecondIngredientID,
		&m.Username, &m.Name, &m.FirstIngredientName, &m.SecondIngredientName,
		&m.SellPrice, &m.BuyPrice, &m.MixTime)
	if err != nil {
		return nil, fmt.Errorf("Repository material, find material by ingredient id: %w", err)
	}

	return &m, nil
}

func (r *repositoryManager) FindMaterialIDByIngrID(firstID int, secondID int) (int, error) {

	var id int

	err := r.db.QueryRow(`
		SELECT id FROM material WHERE first_ingredient_id = ? AND second_ingredient_id = ?`,
		firstID, secondID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Repository material, find material id by ingredient id: %w", err)
	}

	return id, nil
}

func (r *repositoryManager) FindMatMixTimeByID(id int) (int, error) {

	var time int
	if err := r.db.QueryRow(`SELECT mix_time FROM material WHERE id = ?`, id).Scan(&time); err != nil {
		return 0, fmt.Errorf("Repository material, find material mix time by id: %w", err)
	}
	return time, nil
}
