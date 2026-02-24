package repository

type user struct {
	id       int
	username string
	password string
	balance  int
	xp       int
	level    int
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
