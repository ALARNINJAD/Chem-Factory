package repository

type repoMaterial struct {
	id                 int    `db:"id"`
	name               string `db:"name"`
	sellPrice          int    `db:"sell_price"`
	buyPrice           int    `db:"buy_price"`
	firstIngredientID  int    `db:"first_ingredient_id"`
	secondIngredientID int    `db:"second_ingredient_id"`
	mixTime            int    `db:"mix_time"`
}

func (rm *repoMaterial) new(firstIngredientName string, secondIngredientName string) error {

	query := `
		INSERT INTO material(name,sell_price,buy_price,mix_time)
		VALUES (?, ?, ?, ?)`
	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	stmt.Exec(rm.name, rm.sellPrice, rm.buyPrice, rm.mixTime)

	db.QueryRow("SELECT id FROM material WHERE name = ?", rm.name).Scan(&rm.id)
	db.QueryRow("SELECT id FROM material WHERE name = ?", firstIngredientName).Scan(&rm.firstIngredientID)
	db.QueryRow("SELECT id FROM material WHERE name = ?", secondIngredientName).Scan(&rm.secondIngredientID)

	db.Exec("UPDATE material SET first_ingredient_id = ? WHERE id = ?", rm.firstIngredientID, rm.id)
	db.Exec("UPDATE material SET second_ingredient_id = ? WHERE id = ?", rm.secondIngredientID, rm.id)

	return nil
}
