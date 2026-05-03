package repository

import (
	mat "chem-factory/internal/dto/material"
	"database/sql"
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
	ID        int    `json:"id,omitempty"`
	UserID    int    `json:"user_id,omitempty"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	SellPrice int    `json:"sell_price"`
	BuyPrice  int    `json:"buy_price"`
}

type materialManager struct{ db *sql.DB }

func NewMaterialManager(db *sql.DB) *materialManager { return &materialManager{db: db} }

func (m *materialManager) FindBaseByID(id int) (*baseMaterial, error) {
	var mat baseMaterial
	err := m.db.QueryRow(`
		SELECT id, user_id, username, name,	sell_price, buy_price 
		FROM material WHERE id = ? `, id,
	).Scan(&mat.ID, &mat.UserID, &mat.Username, &mat.Name, &mat.SellPrice, &mat.BuyPrice)
	if err != nil {
		return nil, fmt.Errorf("Repository material, find material by id: %w", err)
	}
	return &mat, nil
}

func (m *materialManager) EmptyStruct() *material { return &material{} }

func (m *materialManager) EmptySlice() []material {	return []material{} }

func (m *materialManager) FindByUsername(username string) ([]baseMaterial, error) {

	rows, err := m.db.Query(`
		SELECT id, user_id, username, name, sell_price, buy_price
		FROM material WHERE username = ?`, username)
	if err != nil {
		return nil, fmt.Errorf("Repository material, get base materials by username select query: %w", err)
	}
	defer rows.Close()

	var list []baseMaterial

	for rows.Next() {
		var mat baseMaterial
		err := rows.Scan(
			&mat.ID,
			&mat.UserID,
			&mat.Username,
			&mat.Name,
			&mat.SellPrice,
			&mat.BuyPrice,
		)
		if err != nil {
			return nil, fmt.Errorf("Repository material, get base materials by username rows scan: %w", err)
		}
		list = append(list, mat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("Repository material, get base materials by username rows error: %w", err)
	}
	return list, nil
}

func (m *materialManager) Add(mat material) error {
	_, err := m.db.Exec(`
		INSERT INTO material(user_id, first_ingredient_id, second_ingredient_id,
		username, name, first_ingredient_name, second_ingredient_name,
		sell_price, buy_price, mix_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		mat.UserID, mat.FirstIngredientID, mat.SecondIngredientID,
		mat.Username, mat.Name, mat.FirstIngredientName, mat.SecondIngredientName,
		mat.SellPrice, mat.BuyPrice, mat.MixTime)
	if err != nil {
		return fmt.Errorf("Repository material, save material: %w", err)
	}
	return nil
}

func (m *materialManager) AddBase(tx *sql.Tx, mat mat.BaseMaterial) error {
	_, err := tx.Exec(`
		INSERT INTO material(user_id, username, name, sell_price, buy_price) VALUES (?, ?, ?, ?, ?)`,
		mat.UserID, mat.Username, mat.Name, mat.SellPrice, mat.BuyPrice)
	if err != nil {
		return fmt.Errorf("Repository material, save base material: %w", err)
	}
	return nil
}

func (m *materialManager) FindIDByName(name string) (int, error) {
	var id int
	err := m.db.QueryRow("SELECT id FROM material WHERE name = ?", name).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Repository material, find id by material name: %w", err)
	}
	return id, nil
}

func (m *materialManager) FindByID(id int) (*material, error) {
	var mat material
	err := m.db.QueryRow(`
		SELECT id, user_id, first_ingredient_id, second_ingredient_id,
		username, name, first_ingredient_name, second_ingredient_name,
		sell_price, buy_price, mix_time
		FROM material WHERE id = ? `, id,
	).Scan(&mat.ID, &mat.UserID, &mat.FirstIngredientID, &mat.SecondIngredientID,
		&mat.Username, &mat.Name, &mat.FirstIngredientName, &mat.SecondIngredientName,
		&mat.SellPrice, &mat.BuyPrice, &mat.MixTime)
	if err != nil {
		return nil, fmt.Errorf("Repository material, find material by id: %w", err)
	}
	return &mat, nil
}

func (m *materialManager) FindByIngrID(firstID int, secondID int) (*material, error) {
	var mat material
	err := m.db.QueryRow(`
		SELECT id, user_id, first_ingredient_id, second_ingredient_id,
		username, name, first_ingredient_name, second_ingredient_name,
		sell_price, buy_price, mix_time
		FROM material  WHERE first_ingredient_id = ? AND second_ingredient_id = ?`, firstID, secondID,
	).Scan(&mat.ID, &mat.UserID, &mat.FirstIngredientID, &mat.SecondIngredientID,
		&mat.Username, &mat.Name, &mat.FirstIngredientName, &mat.SecondIngredientName,
		&mat.SellPrice, &mat.BuyPrice, &mat.MixTime)
	if err != nil {
		return nil, fmt.Errorf("Repository material, find material by ingredient id: %w", err)
	}
	return &mat, nil
}

func (m *materialManager) FindIDByIngrID(firstID int, secondID int) (int, error) {
	var id int
	err := m.db.QueryRow(`
		SELECT id FROM material WHERE first_ingredient_id = ? AND second_ingredient_id = ?`,
		firstID, secondID,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("Repository material, find material id by ingredient id: %w", err)
	}
	return id, nil
}

func (m *materialManager) FindMixTimeByID(id int) (int, error) {
	var time int
	if err := m.db.QueryRow(`SELECT mix_time FROM material WHERE id = ?`, id).Scan(&time); err != nil {
		return 0, fmt.Errorf("Repository material, find material mix time by id: %w", err)
	}
	return time, nil
}
