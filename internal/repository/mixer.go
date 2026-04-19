package repository

import "chem-factory/internal/model"

type mixer struct {
	ID                   int    `json:"id,omitempty"`
	UserID               int    `json:"user_id,omitempty"`
	FirstIngredientID    int    `json:"first_ingredient_id,omitempty"`
	SecondIngredientID   int    `json:"second_ingredient_id,omitempty"`
	Username             string `json:"username"`
	FirstIngredientName  string `json:"first_ingredient_name"`
	SecondIngredientName string `json:"second_ingredient_name"`
	Number               int    `json:"number"`
	DateTime             int    `json:"date_time"`
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
		panic("Could not create mixer table.")
	}
}

func (r *repositoryManager) AddToMixer(m model.Mixer) (int, error) {

	userID, err := r.ExportIDbyUsername(m.Username)
	if err != nil {
		return 0, err
	}

	firstIngredientID, err := r.ExportIDbyMaterialName(m.FirstIngredientName)
	if err != nil {
		return 0, err
	}

	secondIngredientID, err := r.ExportIDbyMaterialName(m.SecondIngredientName)
	if err != nil {
		return 0, err
	}

	mxr := mixer{
		UserID:               userID,
		FirstIngredientID:    firstIngredientID,
		SecondIngredientID:   secondIngredientID,
		Username:             m.Username,
		FirstIngredientName:  m.FirstIngredientName,
		SecondIngredientName: m.SecondIngredientName,
		Number:               m.Number,
	}

	err = r.db.QueryRow(`
		SELECT id FROM mixer
		WHERE user_id = ? AND first_ingredient_id = ? AND second_ingredient_id = ?`,
		mxr.UserID, mxr.FirstIngredientID, mxr.SecondIngredientID).Scan(&mxr.ID)
	if err != nil {
		err = r.db.QueryRow(`
			SELECT id FROM mixer
			WHERE user_id = ? AND first_ingredient_id = ? AND second_ingredient_id = ?`,
			mxr.UserID, mxr.SecondIngredientID, mxr.FirstIngredientID).Scan(&mxr.ID)
		if err != nil {
			_, err = r.db.Exec(`
				INSERT INTO mixer(
					user_id,
					first_ingredient_id,
					second_ingredient_id,
					username,
					first_ingredient_name,
					second_ingredient_name,
					number)
				VALUES(?,?,?,?,?,?,?)`,
				mxr.UserID,
				mxr.FirstIngredientID,
				mxr.SecondIngredientID,
				mxr.Username,
				mxr.FirstIngredientName,
				mxr.SecondIngredientName,
				mxr.Number)
			if err != nil {
				return 0, err
			}
		}
	}

	return mxr.ID, nil
}

func (r *repositoryManager) ExportMixRowByID(id int) (mixer, error) {

	var m mixer

	err := r.db.QueryRow(`
		SELECT id, user_id, first_ingredient_id, second_ingredient_id,
		username, first_ingredient_name, second_ingredient_name,
		number, date_time
		FROM mixer WHERE id = ? `, id,
	).Scan(&m.ID, &m.UserID, &m.FirstIngredientID, &m.SecondIngredientID,
		&m.Username, &m.FirstIngredientName, &m.SecondIngredientName,
		&m.Number, &m.DateTime)
	if err != nil {
		return mixer{}, err
	}

	return m, nil
}
