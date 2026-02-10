package repository

type Database struct {
	Path string
}

func New(path string) *Database {
	return &Database{Path: path}
}

// func (database *Database) Open() (*sql.DB, error) {
// 	db, err := sql.Open("sqlite3", "../../migration/db.sqlite3")
// 	return db, err
// }

// func (database *Database) Close() (*sql.DB, error) {
// 	db, err := sql.Open("sqlite3", "../../migration/db.sqlite3")
// 	return db, err
// }
