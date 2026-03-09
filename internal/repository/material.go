package repository

import (
	"chem-factory/internal/model"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func (r *repositoryManager) SaveMaterial(m *model.Material) error {

	r.db.Exec("INSERT OR IGNORE INTO material(name,sell_price,buy_price,mix_time) VALUES (?, ?, ?, ?)",
		m.Name, m.SellPrice, m.BuyPrice, m.MixTime)

	r.db.QueryRow("SELECT id FROM material WHERE name = ?", m.Name).Scan(&m.ID)
	r.db.QueryRow("SELECT id FROM material WHERE name = ?", m.FirstIngredientName).Scan(&m.FirstIngredientID)
	r.db.QueryRow("SELECT id FROM material WHERE name = ?", m.SecondIngredientName).Scan(&m.SecondIngredientID)

	r.db.Exec("UPDATE material SET first_ingredient_id = ? WHERE id = ?", m.FirstIngredientID, m.ID)
	r.db.Exec("UPDATE material SET second_ingredient_id = ? WHERE id = ?", m.SecondIngredientID, m.ID)

	return nil
}

func (r *repositoryManager) createMaterials() {

	var materials []model.Material

	data, err := os.ReadFile(filepath.Join(".", "configs", "materials.json"))
	if err != nil {
		panic("Could not open materials file.")
	}

	json.Unmarshal(data, &materials)

	for _, m := range materials {
		if err = r.SaveMaterial(&m); err != nil {
			panic(fmt.Sprintf("Could not add the %s material to database.", m.Name))
		}
	}
}
