package repository

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type material struct {
	ID                   int    `json:"id,omitempty"`
	FirstIngredientID    int    `json:"first_ingredient_id,omitempty"`
	SecondIngredientID   int    `json:"second_ingredient_id,omitempty"`
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
		first_ingredient_id INTEGER,
		second_ingredient_id INTEGER,
		name TEXT NOT NULL UNIQUE,
		first_ingredient_name TEXT NOT NULL,
		second_ingredient_name TEXT NOT NULL,
		sell_price INTEGER NOT NULL CHECK("sell_price" >= 0),
		buy_price INTEGER NOT NULL CHECK("buy_price" >= 0),
		mix_time INTEGER NOT NULL CHECK("mix_time" >= 0),
		UNIQUE("first_ingredient_id","second_ingredient_id"),
		UNIQUE("first_ingredient_name","second_ingredient_name"),
		FOREIGN KEY("first_ingredient_id") REFERENCES "material"("id"),
		FOREIGN KEY("second_ingredient_id") REFERENCES "material"("id")
	)`
	if _, err := r.db.Exec(query); err != nil {
		panic("Could not create material table.")
	}
}

func (r *repositoryManager) createMaterials() {

	var materials []material

	data, err := os.ReadFile(filepath.Join(".", "configs", "materials.json"))
	if err != nil {
		panic("Could not open materials file.")
	}

	json.Unmarshal(data, &materials)

	for _, m := range materials {
		if err = r.SaveMaterial(m); err != nil {
			panic(fmt.Sprintf("Could not add the %s material to database.", m.Name))
		}
	}
}

func (r *repositoryManager) SaveMaterial(m material) error {

	_, err := r.db.Exec("INSERT OR IGNORE INTO material(name,first_ingredient_name,second_ingredient_name,sell_price,buy_price,mix_time) VALUES (?, ?, ?, ?, ?, ?)",
		m.Name, m.FirstIngredientName, m.SecondIngredientName, m.SellPrice, m.BuyPrice, m.MixTime)
	if err != nil {
		return err
	}

	r.db.QueryRow("SELECT id FROM material WHERE name = ?", m.FirstIngredientName).Scan(&m.FirstIngredientID)
	r.db.QueryRow("SELECT id FROM material WHERE name = ?", m.SecondIngredientName).Scan(&m.SecondIngredientID)

	r.db.QueryRow("SELECT id FROM material WHERE name = ?", m.Name).Scan(&m.ID)
	r.db.Exec("UPDATE material SET first_ingredient_id = ? WHERE id = ?", m.FirstIngredientID, m.ID)
	r.db.Exec("UPDATE material SET second_ingredient_id = ? WHERE id = ?", m.SecondIngredientID, m.ID)

	return nil
}
