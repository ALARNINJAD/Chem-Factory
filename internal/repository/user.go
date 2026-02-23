package repository

type user struct {
	id       int    `db:"id" json:"id"`
	username string `db:"username" json:"username"`
	password string `db:"password" json:"password"`
	balance  int    `db:"balance" json:"balance"`
	xp       int    `db:"xp" json:"xp"`
	level    int    `db:"level" json:"level"`
}

func (usr *user) new() error {
	query := `
		INSERT INTO user(username, password)
		VALUES (?, ?)`

	stmt, err := db.Prepare(query)
	if err != nil {
		return err
	}
	stmt.Exec(usr.username, usr.password)
	return nil
}
