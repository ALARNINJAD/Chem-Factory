package repository

import "chem-factory/internal/model"

type material struct {
	model.Material
}

func (m *material) new() error {

	query := `INSERT OR IGNORE INTO material(name,sell_price,buy_price,mix_time) VALUES (?, ?, ?, ?)`
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	stmt.Exec(m.Name, m.SellPrice, m.BuyPrice, m.MixTime)

	db.QueryRow("SELECT id FROM material WHERE name = ?", m.Name).Scan(&m.ID)
	db.QueryRow("SELECT id FROM material WHERE name = ?", m.FirstIngredientName).Scan(&m.FirstIngredientID)
	db.QueryRow("SELECT id FROM material WHERE name = ?", m.SecondIngredientName).Scan(&m.SecondIngredientID)

	db.Exec("UPDATE material SET first_ingredient_id = ? WHERE id = ?", m.FirstIngredientID, m.ID)
	db.Exec("UPDATE material SET second_ingredient_id = ? WHERE id = ?", m.SecondIngredientID, m.ID)

	return nil
}
