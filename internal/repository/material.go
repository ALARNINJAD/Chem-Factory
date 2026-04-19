package repository

import (
	"chem-factory/internal/model"
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

type baseMaterial struct {
	Username  string `json:"username"`
	Name      string `json:"name"`
	SellPrice int    `json:"sell_price"`
	BuyPrice  int    `json:"buy_price"`
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
		panic("Could not create material table.")
	}
}

func (r *repositoryManager) SaveMaterial(mtrl model.Material) error {

	m := material{
		Username:             mtrl.Username,
		Name:                 mtrl.Name,
		FirstIngredientName:  mtrl.FirstIngredientName,
		SecondIngredientName: mtrl.SecondIngredientName,
		SellPrice:            mtrl.SellPrice,
		BuyPrice:             mtrl.BuyPrice,
		MixTime:              mtrl.MixTime,
	}

	var err error
	m.UserID, err = r.ExportIDbyUsername(m.Username)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(`INSERT OR IGNORE INTO
		material(name,user_id,username,first_ingredient_name,second_ingredient_name,sell_price,buy_price,mix_time)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		m.Name, m.UserID, m.Username, m.FirstIngredientName, m.SecondIngredientName, m.SellPrice, m.BuyPrice, m.MixTime)
	if err != nil {
		return err
	}

	r.db.QueryRow("SELECT id FROM material WHERE name = ?", m.FirstIngredientName).Scan(&m.FirstIngredientID)
	r.db.QueryRow("SELECT id FROM material WHERE name = ?", m.SecondIngredientName).Scan(&m.SecondIngredientID)

	r.db.QueryRow("SELECT id FROM material WHERE name = ?", m.Name).Scan(&m.ID)
	_, err = r.db.Exec("UPDATE material SET first_ingredient_id = ? WHERE id = ?", m.FirstIngredientID, m.ID)
	if err != nil {
		return err
	}
	_, err = r.db.Exec("UPDATE material SET second_ingredient_id = ? WHERE id = ?", m.SecondIngredientID, m.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *repositoryManager) ExportIDbyMaterialName(name string) (int, error) {
	var id int
	err := r.db.QueryRow("SELECT id FROM material WHERE name = ?", name).Scan(&id)
	return id, err
}

func (r *repositoryManager) ExportBaseMaterials() ([]baseMaterial, error) {

	rows, err := r.db.Query(`
        SELECT username, name, sell_price, buy_price
        FROM material 
        WHERE second_ingredient_id = first_ingredient_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []baseMaterial

	for rows.Next() {
		var m baseMaterial

		err := rows.Scan(
			&m.Username,
			&m.Name,
			&m.SellPrice,
			&m.BuyPrice,
		)

		if err != nil {
			return nil, err
		}

		list = append(list, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *repositoryManager) ExportMaterialByID(id int) (material, error) {

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
		return material{}, err
	}

	return m, nil
}

func (r *repositoryManager) ExportMaterialByIngrID(firstID int, secondID int) (material, error) {

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
		return material{}, err
	}

	return m, nil
}

func (r *repositoryManager) ExportMatMixTimeByID(id int) (int, error) {

	var mixTime int

	if err := r.db.QueryRow(`SELECT mix_time FROM material WHERE id = ?`, id).Scan(&mixTime); err != nil {
		return 0, err
	}

	return mixTime, nil
}
