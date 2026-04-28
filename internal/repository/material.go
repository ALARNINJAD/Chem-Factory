package repository

import (
	mat "chem-factory/internal/dto/material"
	"fmt"
)

type material struct {
	ID                   int    `json:"id,omitempty"`
	UserID               int    `json:"user_id,omitempty"`
	FirstIngredientID    int    `json:"first_ingredient_id,omitempty"`
	SecondIngredientID   int    `json:"second_ingredient_id,omitempty"`
	Name                 string `json:"name"`
	Username             string `json:"username"`
	FirstIngredientName  string `json:"first_ingredient_name"`
	SecondIngredientName string `json:"second_ingredient_name"`
	SellPrice            int    `json:"sell_price"`
	BuyPrice             int    `json:"buy_price"`
	MixTime              int    `json:"mix_time"`
}

type baseMaterial struct {
	ID                   int    `json:"id,omitempty"`
	UserID               int    `json:"user_id,omitempty"`
	Name                 string `json:"name"`
	Username             string `json:"username"`
	SellPrice            int    `json:"sell_price"`
	BuyPrice             int    `json:"buy_price"`
}

func (r *repositoryManager) FindBaseMaterialByID(id int) (*baseMaterial, error) {

	var m baseMaterial

	err := r.db.QueryRow(`
		SELECT id, user_id, username, name,	sell_price, buy_price 
		FROM material WHERE id = ? `, id,
	).Scan(&m.ID, &m.UserID, &m.Username, &m.Name, &m.SellPrice, &m.BuyPrice)
	if err != nil {
		return nil, fmt.Errorf("Repository material, find material by id: %w ", err)
	}

	return &m, nil
}

func (r *repositoryManager) EmptyMaterialStruct() *material {
	return &material{}
}

func (r *repositoryManager) EmptyMaterialSlice() []material {
	return []material{}
}

func (r *repositoryManager) GetMaterialsByUsername(username string) ([]material, error) {

	rows, err := r.db.Query(`
		SELECT id, user_id, username, name, sell_price, buy_price
		FROM material WHERE username = ?`, username)
	if err != nil {
		return nil, fmt.Errorf("Repository material, get base materials by username select query: %w ", err)
	}
	defer rows.Close()

	var list []material

	for rows.Next() {
		var s material

		err := rows.Scan(
			&s.ID,
			&s.UserID,
			&s.Username,
			&s.Name,
			&s.SellPrice,
			&s.BuyPrice,
		)

		if err != nil {
			return nil, fmt.Errorf("Repository material, get base materials by username rows scan: %w ", err)
		}

		list = append(list, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Repository material, get base materials by username rows error: %w ", err)
	}

	return list, nil
}

func (r *repositoryManager) SaveMaterial(m material) error {

	_, err := r.db.Exec(`
		INSERT INTO material(user_id, first_ingredient_id, second_ingredient_id,
		username, name, first_ingredient_name, second_ingredient_name,
		sell_price, buy_price, mix_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		m.UserID, m.FirstIngredientID, m.SecondIngredientID,
		m.Username, m.Name, m.FirstIngredientName, m.SecondIngredientName,
		m.SellPrice, m.BuyPrice, m.MixTime)
	if err != nil {
		return fmt.Errorf("Repository material, save material: %w ", err)
	}

	return nil
}

func (r *repositoryManager) SaveBaseMaterial(m mat.BaseMaterial) error {

	_, err := r.db.Exec(`
		INSERT INTO material(user_id, username, name, sell_price, buy_price) VALUES (?, ?, ?, ?, ?)`,
		m.UserID, m.Username, m.Name, m.SellPrice, m.BuyPrice)
	if err != nil {
		return fmt.Errorf("Repository material, save base material: %w ", err)
	}

	return nil
}

func (r *repositoryManager) FindIDbyMaterialName(name string) (int, error) {

	var id int
	err := r.db.QueryRow("SELECT id FROM material WHERE name = ?", name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Repository material, find id by material name: %w ", err)
	}
	return id, nil
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
		return nil, fmt.Errorf("Repository material, find material by id: %w ", err)
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
		return nil, fmt.Errorf("Repository material, find material by ingredient id: %w ", err)
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
		return 0, fmt.Errorf("Repository material, find material id by ingredient id: %w ", err)
	}

	return id, nil
}

func (r *repositoryManager) FindMatMixTimeByID(id int) (int, error) {

	var time int
	if err := r.db.QueryRow(`SELECT mix_time FROM material WHERE id = ?`, id).Scan(&time); err != nil {
		return 0, fmt.Errorf("Repository material, find material mix time by id: %w ", err)
	}
	return time, nil
}
