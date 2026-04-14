package repository

func createInventoryTable(r *repositoryManager) {
	query := `
	CREATE TABLE IF NOT EXISTS inventory (
		id INTEGER NOT NULL UNIQUE PRIMARY KEY AUTOINCREMENT,
		user_id	INTEGER NOT NULL,
		material_id	INTEGER NOT NULL,
		number INTEGER NOT NULL DEFAULT 0 CHECK("number" >= 1),
		UNIQUE("user_id","material_id"),
		FOREIGN KEY("material_id") REFERENCES "material"("id"),
		FOREIGN KEY("user_id") REFERENCES "user"("id")
	)`
	if _, err := r.db.Exec(query); err != nil {
		panic("Could not create inventory table.")
	}
}
